package rate

import (
	"context"
	java2go "github.com/bmstr-ru/java2go/go"
)

type CurrencyRateServiceImpl struct {
	Storage         java2go.CurrencyRateStorage
	ExposureService java2go.TotalExposureService
}

func (s *CurrencyRateServiceImpl) ReceiveRate(rate *java2go.CurrencyRate) error {
	ctx := context.WithValue(context.Background(), "rate", rate)
	if err := s.Storage.SaveRate(ctx, rate); err != nil {
		return err
	}

	return s.ExposureService.RecalculateAllTotalExposure()
}
