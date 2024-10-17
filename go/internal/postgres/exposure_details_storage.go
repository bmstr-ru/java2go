package postgres

import (
	"context"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/jackc/pgx/v5"
)

type ClientExposureDetailStorageImpl struct {
	Postgres *PgPool
}

func (s *ClientExposureDetailStorageImpl) FindByClientIdAndExposureCurrency(ctx context.Context, clientId int64, exposureCurrency string) (*java2go.ClientExposureDetail, error) {
	query := `select client_id, exposure_amount, exposure_currency
				from client_exposure_detail d
				where d.client_id = $1
				and d.exposure_currency = $2`

	row := s.Postgres.DbPool.QueryRow(ctx, query, clientId, exposureCurrency)
	details := java2go.ClientExposureDetail{
		Exposure: &java2go.MonetaryAmount{},
	}

	if err := row.Scan(&details.ClientId, &details.Exposure.Amount, &details.Exposure.Currency); err != nil {
		return nil, err
	}

	return &details, nil
}

func (s *ClientExposureDetailStorageImpl) FindAllByClientId(ctx context.Context, clientId int64) ([]*java2go.ClientExposureDetail, error) {
	query := `select client_id, exposure_amount, exposure_currency
				from client_exposure_detail d
				where d.client_id = $1
				order by d.exposure_currency`

	dealRows, err := s.Postgres.DbPool.Query(ctx, query, clientId)
	if err != nil {
		return nil, err
	}
	defer dealRows.Close()
	return mapRowsToExposureDetails(dealRows)
}

func (s *ClientExposureDetailStorageImpl) FindAll(ctx context.Context) ([]*java2go.ClientExposureDetail, error) {
	query := `select client_id, exposure_amount, exposure_currency
				from client_exposure_detail d
				order by d.client_id, d.exposure_currency`

	dealRows, err := s.Postgres.DbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer dealRows.Close()
	return mapRowsToExposureDetails(dealRows)
}

func mapRowsToExposureDetails(rows pgx.Rows) ([]*java2go.ClientExposureDetail, error) {
	details := []*java2go.ClientExposureDetail{}
	for rows.Next() {
		detail := java2go.ClientExposureDetail{
			Exposure: &java2go.MonetaryAmount{},
		}

		if err := rows.Scan(
			&detail.ClientId,
			&detail.Exposure.Amount,
			&detail.Exposure.Currency); err != nil {
			return nil, err
		}
		details = append(details, &detail)
	}

	return details, nil
}