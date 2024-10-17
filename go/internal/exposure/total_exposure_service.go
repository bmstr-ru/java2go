package exposure

import (
	"context"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/rs/zerolog/log"
)

type TotalExposureServiceImpl struct {
	Storage java2go.DealStorage
}

func (s *TotalExposureServiceImpl) RecalculateAllTotalExposure() error {
	return nil
}

func (s *TotalExposureServiceImpl) RecalculateTotalExposure(clientId int64) error {
	ctx := context.WithValue(context.Background(), "clientId", clientId)
	deals, err := s.Storage.FindAllByClientId(ctx, clientId)
	if err != nil {
		return err
	}
	for _, d := range deals {
		log.Print(d)
	}
	return nil
}

func (s *TotalExposureServiceImpl) GetClientsTotalExposure(clientId int64) *java2go.TotalExposure {
	return nil
}

func (s *TotalExposureServiceImpl) ConsiderNewAmounts(clientId int64, monetaryAmounts []java2go.MonetaryAmount) error {
	return nil
}
