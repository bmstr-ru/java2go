package java2go

import "context"

type DealStorage interface {
	SaveDeal(ctx context.Context, deal *Deal) error
	FindAll(ctx context.Context) ([]*Deal, error)
	FindAllByClientId(ctx context.Context, clientId int64) ([]*Deal, error)
}

type CurrencyRateStorage interface {
	SaveRate(ctx context.Context, rate CurrencyRate) error
	FindByBaseCurrencyAndQuotedCurrency(ctx context.Context, baseCurrency, quotedCurrency string) (*CurrencyRate, error)
	FindAll(ctx context.Context) ([]*CurrencyRate, error)
}

type ClientExposureDetail struct {
	ClientId int64
	Exposure *MonetaryAmount
}

type ClientExposureDetailStorage interface {
	FindByClientIdAndExposureCurrency(ctx context.Context, clientId int64, exposureCurrency string) (*ClientExposureDetail, error)
	FindAllByClientId(ctx context.Context, clientId int64) ([]*ClientExposureDetail, error)
	FindAll(ctx context.Context) ([]*ClientExposureDetail, error)
}
