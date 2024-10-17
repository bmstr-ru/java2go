package ru.bmstr.java2go.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.bmstr.java2go.listener.CurrencyRateMessage;
import ru.bmstr.java2go.persistence.CurrencyRateRepository;

@Slf4j
@Service
@RequiredArgsConstructor
public class CurrencyRateServiceImpl implements CurrencyRateService {

    private final CurrencyRateRepository currencyRateRepository;
    private final TotalExposureService totalExposureService;

    @Override
    @Transactional
    public void receiveRate(CurrencyRateMessage rateMessage) {
        rateMessage.forEach(rate -> {
            currencyRateRepository.saveRate(
                    rate.baseCurrency(),
                    rate.quotedCurrency(),
                    rate.rate()
            );

            log.info("New currency rate record saved: currencyPair={}/{}, rate={}", rate.baseCurrency(), rate.quotedCurrency(), rate.rate());
        });

        totalExposureService.recalculateAllTotalExposure();
    }
}
