package ru.bmstr.java2go.listener;

import java.math.BigDecimal;
import java.util.ArrayList;

public class CurrencyRateMessage extends ArrayList<CurrencyRateMessage.Rate> {

    public record Rate(String currencyPair, BigDecimal rate) {
    }
}
