package ru.bmstr.java2go.service;

import ru.bmstr.java2go.listener.CurrencyRateMessage;

public interface CurrencyRateService {

    void receiveRate(CurrencyRateMessage rateMessage);
}
