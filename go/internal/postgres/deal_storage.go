package postgres

import (
	"context"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/jackc/pgx/v5"
)

type DealStorageImpl struct {
	Postgres *PgPool
}

func (p *DealStorageImpl) SaveDeal(ctx context.Context, deal *java2go.Deal) error {
	tx, err := p.Postgres.DbPool.Begin(ctx)
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

	query := `insert into deal
    		(deal_id, client_id, bought_amount, bought_currency, sold_amount, sold_currency)
			values
    		($1, $2, $3, $4, $5, $6)`
	if _, err = tx.Exec(ctx, query,
		deal.Id,
		deal.ClientId,
		deal.AmountBought.Amount,
		deal.AmountBought.Currency,
		deal.AmountSold.Amount,
		deal.AmountSold.Currency); err != nil {
		return err
	}
	return nil
}

func (p *DealStorageImpl) FindAll(ctx context.Context) ([]*java2go.Deal, error) {
	query := `select deal_id, client_id, bought_amount, bought_currency, sold_amount, sold_currency
				from deal d
				order by d.deal_id`

	dealRows, err := p.Postgres.DbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer dealRows.Close()
	return mapRowsToDeal(dealRows)
}

func (p *DealStorageImpl) FindAllByClientId(ctx context.Context, clientId int64) ([]*java2go.Deal, error) {
	query := `select deal_id, client_id, bought_amount, bought_currency, sold_amount, sold_currency
				from deal d
				where d.client_id = $1
				order by d.deal_id`

	dealRows, err := p.Postgres.DbPool.Query(ctx, query, clientId)
	if err != nil {
		return nil, err
	}
	defer dealRows.Close()
	return mapRowsToDeal(dealRows)
}

func mapRowsToDeal(dealRows pgx.Rows) ([]*java2go.Deal, error) {
	deals := []*java2go.Deal{}
	for dealRows.Next() {
		deal := java2go.Deal{
			AmountBought: java2go.MonetaryAmount{},
			AmountSold:   java2go.MonetaryAmount{},
		}

		if err := dealRows.Scan(
			&deal.Id,
			&deal.ClientId,
			&deal.AmountBought.Amount,
			&deal.AmountBought.Currency,
			&deal.AmountSold.Amount,
			&deal.AmountSold.Currency); err != nil {
			return nil, err
		}
		deals = append(deals, &deal)
	}

	return deals, nil
}
