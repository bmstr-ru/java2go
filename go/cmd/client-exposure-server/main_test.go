package main

import (
	"context"
	"github.com/rs/zerolog/log"
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

	cfg := ConfigStruct{
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
	log.Print(cfg)
	log.Info().Msg(pgPort.Port())
}
