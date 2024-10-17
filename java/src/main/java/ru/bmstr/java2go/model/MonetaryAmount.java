package ru.bmstr.java2go.model;

import java.math.BigDecimal;
import java.math.RoundingMode;

public record MonetaryAmount(
        BigDecimal amount,
        String currency
) {

    public MonetaryAmount convert(String targetCurrency, BigDecimal rate) {
        BigDecimal newAmount = amount.divide(rate, 4, RoundingMode.HALF_UP);
        return new MonetaryAmount(newAmount, targetCurrency);
    }

    public MonetaryAmount add(MonetaryAmount other) {
        if (!currency.equals(other.currency)) {
            throw new IllegalArgumentException("Cannot sum amounts of different currencies");
        }
        return new MonetaryAmount(amount.add(other.amount), currency);
    }

    public MonetaryAmount negate() {
        return new MonetaryAmount(amount.negate(), currency);
    }
}
