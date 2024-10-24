package main

import (
	"fmt"
	java2go "github.com/bmstr-ru/java2go/go"
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

var quit = make(chan os.Signal, 1)

func main() {
	cfg := GetDefaultConfig()
	pgPool := createDbPool(cfg)

	migrateDb(pgPool)

	dealService, rateService, exposureService := createServices(pgPool)

	startRateListener(cfg, rateService)
	startDealListener(cfg, dealService)
	startMainServer(cfg, exposureService)

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

func createDbPool(cfg *ConfigStruct) *postgres.PgPool {
	pgPool := &postgres.PgPool{
		Host:     cfg.Db.Host,
		Port:     cfg.Db.Port,
		Username: cfg.Db.Username,
		Password: cfg.Db.Password,
		Database: cfg.Db.Database,
		Schema:   cfg.Db.Schema,
	}
	pgPool.Init()
	return pgPool
}

func createServices(pgPool *postgres.PgPool) (java2go.DealService, java2go.CurrencyRateService, java2go.TotalExposureService) {

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
	return dealService, rateService, exposureService
}

func startRateListener(cfg *ConfigStruct, rateService java2go.CurrencyRateService) {
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

func startDealListener(cfg *ConfigStruct, dealService java2go.DealService) {
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

func startMainServer(cfg *ConfigStruct, exposureService java2go.TotalExposureService) {
	router := createRouter(exposureService)

	go func() {
		log.Print("I am starting...")
		log.Fatal().Err(http.ListenAndServe(":"+cfg.ServerPort, router)).Msg("")
	}()
}

func createRouter(exposureService java2go.TotalExposureService) http.Handler {
	router := httprouter.New()
	router.GET("/health", httphandler.Health)
	router.GET("/client/:clientId/summary", httphandler.GetClientSummary(exposureService))
	return router
}

func migrateDb(pgPool *postgres.PgPool) {
	workingDir, _ := os.Getwd()
	log.Info().Msg(workingDir)
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
