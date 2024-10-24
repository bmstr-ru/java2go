package main

import (
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type TestLogConsumer struct {
	Service string
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	log.Info().Str("service", g.Service).Msg(string(l.Content))
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
		testcontainers.WithLogConsumers(&TestLogConsumer{Service: "postgres"}),
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

//
//func startActiveMq() *testcontainers.Container {
//	mqContainer, err := artemis.Run(ctx,
//		,
//		,
//		testcontainers.WithWaitStrategy(,
//		,
//	),
//)
//	if err != nil {
//		panic(err)
//	}
//	return mqContainer
//}

func startActiveMq() testcontainers.Container {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "docker.io/apache/activemq-classic:6.1.2",
			ExposedPorts: []string{"61613/tcp", "8161/tcp"},
			WaitingFor:   wait.ForListeningPort("8161/tcp"),
			LogConsumerCfg: &testcontainers.LogConsumerConfig{
				Consumers: []testcontainers.LogConsumer{
					&TestLogConsumer{Service: "activemq"},
				},
			},
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		panic(err)
	}

	return container
}

func stopActiveMq(mqContainer testcontainers.Container) {
	if err := testcontainers.TerminateContainer(mqContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
