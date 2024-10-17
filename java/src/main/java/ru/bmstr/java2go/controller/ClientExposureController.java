package ru.bmstr.java2go.controller;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;
import ru.bmstr.java2go.service.ClientExposureService;

@Slf4j
@RestController
@RequiredArgsConstructor
public class ClientExposureController {

    private final ClientExposureService clientExposureService;

    @GetMapping("/client/{clientId}/summary")
    public ClientExposureSummary getClientSummary(@PathVariable("clientId") Long clientId) {
        return clientExposureService.getClientExposureSummary(clientId);
    }
}
