package main

import (
	"github.com/bmstr-ru/java2go/go/internal/activemq"
	"github.com/bmstr-ru/java2go/go/internal/deal"
	"github.com/bmstr-ru/java2go/go/internal/httphandler"
	"github.com/bmstr-ru/java2go/go/internal/postgres"
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

	dealStorage := &postgres.DealStorageImpl{
		Postgres: pgPool,
	}
	dealService := &deal.DealServiceImpl{
		Storage: dealStorage,
	}

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
