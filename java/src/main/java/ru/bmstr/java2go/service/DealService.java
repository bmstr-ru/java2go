package ru.bmstr.java2go.service;

import ru.bmstr.java2go.listener.DealMessage;

public interface DealService {

    void receiveDeal(DealMessage dealMessage);
}
