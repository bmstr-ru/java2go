package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"os"
)

type PgPool struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Schema   string
	DbPool   *pgxpool.Pool
}

func (p *PgPool) Init() {
	connstring := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s search_path=%s target_session_attrs=read-write",
		p.Host, p.Port, p.Database, p.Username, p.Password, p.Schema)

	connConfig, err := pgxpool.ParseConfig(connstring)
	if err != nil {
		log.Printf("Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	if p.DbPool, err = pgxpool.NewWithConfig(context.Background(), connConfig); err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
}
