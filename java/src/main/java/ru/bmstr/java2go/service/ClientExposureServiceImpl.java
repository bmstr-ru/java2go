package ru.bmstr.java2go.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import ru.bmstr.java2go.controller.ClientExposureSummary;
import ru.bmstr.java2go.model.TotalExposure;

@Slf4j
@Service
@RequiredArgsConstructor
public class ClientExposureServiceImpl implements ClientExposureService {

    private final TotalExposureService totalExposureService;

    @Override
    public ClientExposureSummary getClientExposureSummary(Long clientId) {
        TotalExposure totalExposure = totalExposureService.getClientsTotalExposure(clientId);
        return ClientExposureSummary.builder()
                .clientId(clientId)
                .total(totalExposure.total())
                .amounts(totalExposure.amounts())
                .build();
    }
}
