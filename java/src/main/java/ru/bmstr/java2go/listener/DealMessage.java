package ru.bmstr.java2go.listener;

import lombok.Builder;
import ru.bmstr.java2go.model.MonetaryAmount;

@Builder
public record DealMessage(
        Long id,
        Long clientId,
        MonetaryAmount amountBought,
        MonetaryAmount amountSold
) {
}
