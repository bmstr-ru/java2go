package postgres

import (
	"context"
	java2go "github.com/bmstr-ru/java2go/go"
)

type CurrencyRateStorageImpl struct {
	Postgres *PgPool
}

func (s *CurrencyRateStorageImpl) SaveRate(ctx context.Context, rate java2go.CurrencyRate) error {

}

func (s *CurrencyRateStorageImpl) FindByBaseCurrencyAndQuotedCurrency(ctx context.Context, baseCurrency, quotedCurrency string) (*java2go.CurrencyRate, error) {

}

func (s *CurrencyRateStorageImpl) FindAll(ctx context.Context) ([]*java2go.CurrencyRate, error) {

}
