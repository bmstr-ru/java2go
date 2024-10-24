package exposure

import (
	"context"
	java2go "github.com/bmstr-ru/java2go/go"
)

const baseCurrency = "EUR"

type TotalExposureServiceImpl struct {
	DealStorage           java2go.DealStorage
	ExposureDetailStorage java2go.ClientExposureDetailStorage
	TotalExposureStorage  java2go.ClientExposureStorage
	RateStorage           java2go.CurrencyRateStorage
}

func (s *TotalExposureServiceImpl) RecalculateAllTotalExposure() error {
	allDetails, err := s.ExposureDetailStorage.FindAll(context.Background())
	if err != nil {
		return err
	}
	mapByClient := make(map[int64][]*java2go.ClientExposure)
	for _, detail := range allDetails {
		exposureDetails, ok := mapByClient[detail.ClientId]
		if !ok {
			exposureDetails = []*java2go.ClientExposure{}
			mapByClient[detail.ClientId] = exposureDetails
		}
		exposureDetails = append(exposureDetails, detail)
	}

	for clientId, exposureDetails := range mapByClient {
		if err = s.recalculateClientTotalExposure(clientId, exposureDetails); err != nil {
			return err
		}
	}
	return nil
}

func (s *TotalExposureServiceImpl) GetClientsTotalExposure(clientId int64) (*java2go.TotalExposure, error) {
	ctx := context.WithValue(context.Background(), "clientId", clientId)
	clientExposure, err := s.TotalExposureStorage.FindByClientId(ctx, clientId)
	if err != nil {
		return nil, err
	}
	if clientExposure == nil {
		clientExposure = &java2go.ClientExposure{
			ClientId: clientId,
			Exposure: &java2go.MonetaryAmount{
				Currency: baseCurrency,
				Amount:   0,
			},
		}
	}

	exposureDetails, err := s.ExposureDetailStorage.FindAllByClientId(ctx, clientId)
	if err != nil {
		return nil, err
	}
	totalExposure := &java2go.TotalExposure{
		ClientId: clientId,
		Total:    *clientExposure.Exposure,
		Amounts:  []java2go.MonetaryAmount{},
	}
	for _, details := range exposureDetails {
		totalExposure.Amounts = append(totalExposure.Amounts, *details.Exposure)
	}

	return totalExposure, nil
}

func (s *TotalExposureServiceImpl) ConsiderNewAmounts(clientId int64, monetaryAmounts ...java2go.MonetaryAmount) error {
	ctx := context.WithValue(context.Background(), "clientId", clientId)
	for _, amount := range monetaryAmounts {
		exposure, err := s.ExposureDetailStorage.FindByClientIdAndExposureCurrency(ctx, clientId, amount.Currency)
		if err != nil {
			return err
		}

		if exposure != nil {
			exposure.Exposure.Amount += amount.Amount
		} else {
			exposure = &java2go.ClientExposure{ClientId: clientId, Exposure: &amount}
		}

		if err = s.ExposureDetailStorage.Save(ctx, exposure); err != nil {
			return err
		}
	}
	return s.recalculateClientExposure(clientId)
}

func (s *TotalExposureServiceImpl) recalculateClientExposure(clientId int64) error {
	ctx := context.WithValue(context.Background(), "clientId", clientId)
	details, err := s.ExposureDetailStorage.FindAllByClientId(ctx, clientId)
	if err != nil {
		return err
	}
	return s.recalculateClientTotalExposure(clientId, details)
}

func (s *TotalExposureServiceImpl) recalculateClientTotalExposure(clientId int64, details []*java2go.ClientExposure) error {
	ctx := context.WithValue(context.Background(), "clientId", clientId)
	totalExposure := java2go.MonetaryAmount{
		Currency: baseCurrency,
		Amount:   0,
	}
	for _, detail := range details {
		amount := detail.Exposure
		if amount.Currency == baseCurrency {
			totalExposure.Amount += amount.Amount
			continue
		}
		rate, err := s.RateStorage.FindByBaseCurrencyAndQuotedCurrency(ctx, baseCurrency, amount.Currency)
		if err != nil {
			return err
		}
		totalExposure.Amount += amount.Amount * rate.Rate
	}
	clientExposure := java2go.ClientExposure{
		ClientId: clientId,
		Exposure: &totalExposure,
	}
	return s.TotalExposureStorage.Save(ctx, &clientExposure)
}
