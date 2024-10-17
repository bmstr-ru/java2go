package ru.bmstr.java2go.controller;

import lombok.Builder;
import ru.bmstr.java2go.model.MonetaryAmount;

import java.util.List;

@Builder
public record ClientExposureSummary(
        Long clientId,
        List<MonetaryAmount> amounts,
        MonetaryAmount total
) {
}
