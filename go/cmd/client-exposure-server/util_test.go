package main

import (
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/artemis"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type TestLogConsumer struct {
	Msgs []string
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	log.Info().Msg(string(l.Content))
}

const (
	pgDbName   = "postgres"
	pgUser     = "postgres"
	pgPassword = "password"
)

func startPostgres() *postgres.PostgresContainer {
	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:latest",
		postgres.WithDatabase(pgDbName),
		postgres.WithUsername(pgUser),
		postgres.WithPassword(pgPassword),
		testcontainers.WithWaitStrategy(
			wait.ForExec([]string{"pg_isready"}).
				WithStartupTimeout(30*time.Second),
		),
		testcontainers.WithLogConsumers(&TestLogConsumer{}),
	)
	if err != nil {
		panic(err)
	}
	return postgresContainer
}

func stopPostgres(pgContainer *postgres.PostgresContainer) {
	if err := testcontainers.TerminateContainer(pgContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}

func startActiveMq() *artemis.Container {
	mqContainer, err := artemis.Run(ctx,
		"docker.io/apache/activemq-classic:6.1.2",
		testcontainers.WithLogConsumers(&TestLogConsumer{}),
		artemis.WithCredentials("admin", "admin"),
	)
	if err != nil {
		panic(err)
	}
	return mqContainer
}

func stopActiveMq(mqContainer *artemis.Container) {
	if err := testcontainers.TerminateContainer(mqContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
