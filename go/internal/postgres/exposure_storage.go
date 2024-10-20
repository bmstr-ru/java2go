package postgres

import (
	"context"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/jackc/pgx/v5"
)

type ClientExposureStorageImpl struct {
	Postgres *PgPool
}

func (s *ClientExposureStorageImpl) Save(ctx context.Context, exposure *java2go.ClientExposure) error {
	tx, err := s.Postgres.DbPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := `insert into client_exposure
    		(client_id, total_exposure_amount, total_exposure_currency)
			values
    		($1, $2, $3)
    		on conflict (client_id)
			do update set total_exposure_amount = $2, total_exposure_currency = $3`
	if _, err = tx.Exec(ctx, query, exposure.ClientId, exposure.Exposure.Amount, exposure.Exposure.Currency); err != nil {
		return err
	}
	return nil
}

func (s *ClientExposureStorageImpl) FindByClientId(ctx context.Context, clientId int64) (*java2go.ClientExposure, error) {
	query := `select total_exposure_amount, total_exposure_currency
				from client_exposure
				where client_id = $1`

	row := s.Postgres.DbPool.QueryRow(ctx, query, clientId)
	exposure := java2go.ClientExposure{
		ClientId: clientId,
		Exposure: &java2go.MonetaryAmount{},
	}

	if err := row.Scan(&exposure.Exposure.Amount, &exposure.Exposure.Currency); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &exposure, nil
}
