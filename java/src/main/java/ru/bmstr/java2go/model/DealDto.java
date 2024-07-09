package ru.bmstr.java2go.model;

import java.math.BigDecimal;

public record DealDto(
        Long id,
        Long clientId,
        MonetaryAmount amountBought,
        MonetaryAmount amountSold
) {
    public record MonetaryAmount(
            String currency,
            BigDecimal amount
    ) {}
}
