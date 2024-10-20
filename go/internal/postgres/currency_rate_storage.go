package postgres

import (
	"context"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/jackc/pgx/v5"
)

type CurrencyRateStorageImpl struct {
	Postgres *PgPool
}

func (s *CurrencyRateStorageImpl) SaveRate(ctx context.Context, rate *java2go.CurrencyRate) error {
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

	query := `insert into currency_rate
    		(base_currency, quoted_currency, rate)
			values
    		($1, $2, $3)
    		on conflict (base_currency, quoted_currency)
			do update set rate = $3`
	if _, err = tx.Exec(ctx, query, rate.BaseCurrency, rate.QuotedCurrency, rate.Rate); err != nil {
		return err
	}
	return nil
}

func (s *CurrencyRateStorageImpl) FindByBaseCurrencyAndQuotedCurrency(ctx context.Context, baseCurrency, quotedCurrency string) (*java2go.CurrencyRate, error) {
	query := `select base_currency, quoted_currency, rate
				from currency_rate
				where base_currency = $1
				and quoted_currency = $2`

	row := s.Postgres.DbPool.QueryRow(ctx, query, baseCurrency, quotedCurrency)
	rate := java2go.CurrencyRate{}

	if err := row.Scan(&rate.BaseCurrency, &rate.QuotedCurrency, &rate.Rate); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rate, nil
}

func (s *CurrencyRateStorageImpl) FindAll(ctx context.Context) ([]*java2go.CurrencyRate, error) {
	query := `select base_currency, quoted_currency, rate
				from currency_rate
				order by base_currency, quoted_currency`

	rows, err := s.Postgres.DbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return mapRowsCurrencyRate(rows)
}

func mapRowsCurrencyRate(rows pgx.Rows) ([]*java2go.CurrencyRate, error) {
	rates := []*java2go.CurrencyRate{}
	for rows.Next() {
		rate := java2go.CurrencyRate{}

		if err := rows.Scan(&rate.BaseCurrency, &rate.QuotedCurrency, &rate.Rate); err != nil {
			return nil, err
		}
		rates = append(rates, &rate)
	}

	return rates, nil
}
