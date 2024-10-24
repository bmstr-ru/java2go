package main

import (
	"context"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var ctx = context.Background()

func Test_e2e(t *testing.T) {
	pgContainer := startPostgres()
	defer stopPostgres(pgContainer)

	mqContainer := startActiveMq()
	defer stopActiveMq(mqContainer)

	pgPort, _ := pgContainer.MappedPort(context.Background(), "5432/tcp")
	mqPort, _ := mqContainer.MappedPort(context.Background(), "61613/tcp")

	cfg := &ConfigStruct{
		ServerPort: "0",
		Db: Db{
			Host:     "localhost",
			Port:     pgPort.Int(),
			Username: pgUser,
			Password: pgPassword,
			Database: pgDbName,
			Schema:   "public",
		},
		Mq: Mq{
			Url: "localhost:" + mqPort.Port(),
			Queue: Queue{
				Deals: "q.deals",
				Rates: "q.rates",
			},
		},
	}

	pgPool := createDbPool(cfg)

	os.Chdir("../..")
	migrateDb(pgPool)

	dealService, rateService, exposureService := createServices(pgPool)

	startRateListener(cfg, rateService)
	startDealListener(cfg, dealService)
	startMainServer(cfg, exposureService)

	server := httptest.NewServer(createRouter(exposureService))
	defer server.Close()

	e := httpexpect.Default(t, server.URL)
	e.GET("/health").
		Expect().
		Status(http.StatusOK).JSON().Array().IsEmpty()
}
