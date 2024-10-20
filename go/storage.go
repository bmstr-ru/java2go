package java2go

import "context"

type DealStorage interface {
	SaveDeal(ctx context.Context, deal *Deal) error
	FindAll(ctx context.Context) ([]*Deal, error)
	FindAllByClientId(ctx context.Context, clientId int64) ([]*Deal, error)
}

type CurrencyRateStorage interface {
	SaveRate(ctx context.Context, rate *CurrencyRate) error
	FindByBaseCurrencyAndQuotedCurrency(ctx context.Context, baseCurrency, quotedCurrency string) (*CurrencyRate, error)
	FindAll(ctx context.Context) ([]*CurrencyRate, error)
}

type ClientExposure struct {
	ClientId int64
	Exposure *MonetaryAmount
}

type ClientExposureDetailStorage interface {
	FindByClientIdAndExposureCurrency(ctx context.Context, clientId int64, exposureCurrency string) (*ClientExposure, error)
	FindAllByClientId(ctx context.Context, clientId int64) ([]*ClientExposure, error)
	FindAll(ctx context.Context) ([]*ClientExposure, error)
	Save(ctx context.Context, detail *ClientExposure) error
}

type ClientExposureStorage interface {
	Save(ctx context.Context, totalExposure *ClientExposure) error
	FindByClientId(ctx context.Context, clientId int64) (*ClientExposure, error)
}
