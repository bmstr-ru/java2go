package ru.bmstr.java2go.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import ru.bmstr.java2go.model.MonetaryAmount;
import ru.bmstr.java2go.model.TotalExposure;
import ru.bmstr.java2go.persistence.ClientExposureDetailRecord;
import ru.bmstr.java2go.persistence.ClientExposureDetailRepository;
import ru.bmstr.java2go.persistence.ClientExposureRecord;
import ru.bmstr.java2go.persistence.ClientExposureRepository;
import ru.bmstr.java2go.persistence.CurrencyRateRecord;
import ru.bmstr.java2go.persistence.CurrencyRateRepository;

import java.math.BigDecimal;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class TotalExposureServiceImpl implements TotalExposureService{

    private static final String BASE_CURRENCY = "EUR";
    private static final MonetaryAmount ZERO = new MonetaryAmount(BigDecimal.ZERO, BASE_CURRENCY);

    private final CurrencyRateRepository currencyRateRepository;
    private final ClientExposureRepository clientExposureRepository;
    private final ClientExposureDetailRepository clientExposureDetailRepository;

    @Override
    public void recalculateAllTotalExposure() {
        Map<Long, List<ClientExposureDetailRecord>> dealsByClientId = clientExposureDetailRepository.findAll().stream()
                .collect(Collectors.groupingBy(ClientExposureDetailRecord::getClientId, Collectors.toList()));
        dealsByClientId.forEach(this::recalculateTotalExposure);
    }

    private void recalculateTotalExposure(Long clientId, List<ClientExposureDetailRecord> detailRecords) {
        MonetaryAmount totalExposure = detailRecords.stream()
                .map(this::toMonetaryAmounts)
                .map(this::toBaseCurrency)
                .reduce(MonetaryAmount::add)
                .orElse(ZERO);
        clientExposureRepository.saveClientExposure(clientId, totalExposure.amount(), totalExposure.currency());
        log.info("Recalculated client exposure: clientId={}", clientId);
    }

    private MonetaryAmount toMonetaryAmounts(ClientExposureDetailRecord record) {
        return new MonetaryAmount(record.getExposureAmount(), record.getExposureCurrency());
    }

    private MonetaryAmount toBaseCurrency(MonetaryAmount monetaryAmount) {
        if (BASE_CURRENCY.equals(monetaryAmount.currency())) {
            return monetaryAmount;
        }
        CurrencyRateRecord rate = currencyRateRepository.findByBaseCurrencyAndQuotedCurrency(BASE_CURRENCY, monetaryAmount.currency());
        return monetaryAmount.convert(BASE_CURRENCY, rate.getRate());
    }

    @Override
    public TotalExposure getClientsTotalExposure(Long clientId) {
        MonetaryAmount total = clientExposureRepository.findByClientId(clientId)
                .map(this::toMonetaryAmount)
                .orElse(ZERO);
        List<MonetaryAmount> details = clientExposureDetailRepository.findAllByClientId(clientId).stream()
                .map(this::toMonetaryAmount)
                .sorted(Comparator.comparing(MonetaryAmount::currency))
                .toList();

        return new TotalExposure(details, total);
    }

    private MonetaryAmount toMonetaryAmount(ClientExposureRecord record) {
        return new MonetaryAmount(record.getTotalExposureAmount(), record.getTotalExposureCurrency());
    }

    private MonetaryAmount toMonetaryAmount(ClientExposureDetailRecord record) {
        return new MonetaryAmount(record.getExposureAmount(), record.getExposureCurrency());
    }

    @Override
    public void considerNewAmounts(Long clientId, MonetaryAmount... monetaryAmounts) {
        Arrays.stream(monetaryAmounts).forEach(monetaryAmount -> {
            ClientExposureDetailRecord record = clientExposureDetailRepository
                    .findByClientIdAndExposureCurrency(clientId, monetaryAmount.currency())
                    .map(r -> add(r, monetaryAmount))
                    .orElse(toClientExposureDetailRecord(clientId, monetaryAmount));
            clientExposureDetailRepository.save(record);
        });
        recalculateTotalExposure(clientId, clientExposureDetailRepository.findAllByClientId(clientId));
    }

    private ClientExposureDetailRecord toClientExposureDetailRecord(Long clientId, MonetaryAmount monetaryAmount) {
        return ClientExposureDetailRecord.builder()
                .clientId(clientId)
                .exposureAmount(monetaryAmount.amount())
                .exposureCurrency(monetaryAmount.currency())
                .build();
    }

    private ClientExposureDetailRecord add(ClientExposureDetailRecord record, MonetaryAmount monetaryAmount) {
        return record.toBuilder()
                .exposureAmount(record.getExposureAmount().add(monetaryAmount.amount()))
                .build();
    }

}