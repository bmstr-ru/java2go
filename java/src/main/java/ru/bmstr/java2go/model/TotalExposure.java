package ru.bmstr.java2go.model;

import java.util.List;

public record TotalExposure(
        List<MonetaryAmount> amounts,
        MonetaryAmount total
) {
}
