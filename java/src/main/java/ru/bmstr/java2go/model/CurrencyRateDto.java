package ru.bmstr.java2go.model;

import java.math.BigDecimal;

public record CurrencyRateDto(
        String currencyPair,
        BigDecimal rate
) {
}
