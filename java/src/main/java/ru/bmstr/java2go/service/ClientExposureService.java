package ru.bmstr.java2go.service;

import ru.bmstr.java2go.controller.ClientExposureSummary;

public interface ClientExposureService {
    ClientExposureSummary getClientExposureSummary(Long clientId);
}
