package ru.bmstr.java2go.service;

import ru.bmstr.java2go.model.MonetaryAmount;
import ru.bmstr.java2go.model.TotalExposure;

public interface TotalExposureService {

    void recalculateAllTotalExposure();

    TotalExposure getClientsTotalExposure(Long clientId);

    void considerNewAmounts(Long clientId, MonetaryAmount... monetaryAmounts);

}