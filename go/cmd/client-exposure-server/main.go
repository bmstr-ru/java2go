package main

import (
	"fmt"
	"github.com/bmstr-ru/java2go/go/internal/activemq"
	"github.com/bmstr-ru/java2go/go/internal/deal"
	"github.com/bmstr-ru/java2go/go/internal/exposure"
	"github.com/bmstr-ru/java2go/go/internal/httphandler"
	"github.com/bmstr-ru/java2go/go/internal/postgres"
	"github.com/bmstr-ru/java2go/go/internal/rate"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := GetDefaultConfig()
	pgPool := &postgres.PgPool{
		Host:     cfg.Db.Host,
		Port:     cfg.Db.Port,
		Username: cfg.Db.Username,
		Password: cfg.Db.Password,
		Database: cfg.Db.Database,
		Schema:   cfg.Db.Schema,
	}
	pgPool.Init()

	migrateDb(pgPool)

	dealStorage := &postgres.DealStorageImpl{
		Postgres: pgPool,
	}
	rateStorage := &postgres.CurrencyRateStorageImpl{
		Postgres: pgPool,
	}
	exposureDetailsStorage := &postgres.ClientExposureDetailStorageImpl{
		Postgres: pgPool,
	}
	totalExposureStorage := &postgres.ClientExposureStorageImpl{
		Postgres: pgPool,
	}

	exposureService := &exposure.TotalExposureServiceImpl{
		DealStorage:           dealStorage,
		ExposureDetailStorage: exposureDetailsStorage,
		TotalExposureStorage:  totalExposureStorage,
		RateStorage:           rateStorage,
	}
	dealService := &deal.DealServiceImpl{
		Storage:         dealStorage,
		ExposureService: exposureService,
	}
	rateService := &rate.CurrencyRateServiceImpl{
		Storage:         rateStorage,
		ExposureService: exposureService,
	}

	startRateListener(cfg, rateService)
	startDealListener(cfg, dealService)
	startMainServer(cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-quit
	log.Print("I am dying...")
}

func startRateListener(cfg *ConfigStruct, rateService *rate.CurrencyRateServiceImpl) {
	ratesChan, err := activemq.StartRateListener(cfg.Mq.Url, cfg.Mq.Queue.Rates)
	if err != nil {
		panic(err)
	}

	go func() {
		for r := range ratesChan {
			err = rateService.ReceiveRate(r)
			if err != nil {
				log.Error().Err(err).Msg("Error while processing new rates")
			}
		}
	}()
}

func startDealListener(cfg *ConfigStruct, dealService *deal.DealServiceImpl) {
	dealsChan, err := activemq.StartDealListener(cfg.Mq.Url, cfg.Mq.Queue.Deals)
	if err != nil {
		panic(err)
	}

	go func() {
		for d := range dealsChan {
			err = dealService.ReceiveDeal(d)
			if err != nil {
				log.Error().Err(err).Msg("Error while processing new deal")
			}
		}
	}()
}

func startMainServer(cfg *ConfigStruct) {
	router := httprouter.New()
	router.GET("/health", httphandler.Health())

	go func() {
		log.Print("I am starting...")
		log.Fatal().Err(http.ListenAndServe(":"+cfg.ServerPort, router)).Msg("")
	}()
}

func migrateDb(pgPool *postgres.PgPool) {
	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?search_path=%s&sslmode=disable",
			pgPool.Username, pgPool.Password, pgPool.Host, pgPool.Port, pgPool.Database, pgPool.Schema))
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
