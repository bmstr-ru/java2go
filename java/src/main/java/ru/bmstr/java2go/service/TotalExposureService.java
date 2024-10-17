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
import ru.bmstr.java2go.persistence.DealRecord;
import ru.bmstr.java2go.persistence.DealRepository;

import java.math.BigDecimal;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import java.util.stream.StreamSupport;

@Slf4j
@Service
@RequiredArgsConstructor
public class TotalExposureService {

    private static final String BASE_CURRENCY = "EUR";
    private static final MonetaryAmount ZERO = new MonetaryAmount(BigDecimal.ZERO, BASE_CURRENCY);

    private final DealRepository dealRepository;
    private final CurrencyRateRepository currencyRateRepository;
    private final ClientExposureRepository clientExposureRepository;
    private final ClientExposureDetailRepository clientExposureDetailRepository;

    public void recalculateAllTotalExposure() {
        Map<Long, List<DealRecord>> dealsByClientId = StreamSupport.stream(dealRepository.findAll().spliterator(), false)
                .collect(Collectors.groupingBy(DealRecord::getClientId, Collectors.toList()));
        dealsByClientId.forEach(this::recalculateTotalExposure);
    }

    public void recalculateTotalExposure(Long clientId) {
        recalculateTotalExposure(clientId, dealRepository.findAllByClientId(clientId));
    }

    private void recalculateTotalExposure(Long clientId, List<DealRecord> deals) {
        MonetaryAmount totalExposure = deals.stream()
                .flatMap(this::toMonetaryAmounts)
                .map(this::toBaseCurrency)
                .reduce(MonetaryAmount::add)
                .orElse(ZERO);
        clientExposureRepository.saveClientExposure(clientId, totalExposure.amount(), totalExposure.currency());
        log.info("Recalculated client exposure: clientId={}", clientId);
    }

    private Stream<MonetaryAmount> toMonetaryAmounts(DealRecord record) {
        return Stream.of(
                new MonetaryAmount(record.getBoughtAmount(), record.getBoughtCurrency()),
                new MonetaryAmount(record.getSoldAmount().negate(), record.getSoldCurrency())
        );
    }

    private MonetaryAmount toBaseCurrency(MonetaryAmount monetaryAmount) {
        if (BASE_CURRENCY.equals(monetaryAmount.currency())) {
            return monetaryAmount;
        }
        CurrencyRateRecord rate = currencyRateRepository.findByBaseCurrencyAndQuotedCurrency(BASE_CURRENCY, monetaryAmount.currency());
        return monetaryAmount.convert(BASE_CURRENCY, rate.getRate());
    }

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

    public void considerNewAmounts(Long clientId, MonetaryAmount... monetaryAmounts) {
        Arrays.stream(monetaryAmounts).forEach(monetaryAmount -> {
            ClientExposureDetailRecord record = clientExposureDetailRepository
                    .findByClientIdAndExposureCurrency(clientId, monetaryAmount.currency())
                    .map(r -> add(r, monetaryAmount))
                    .orElse(toClientExposureDetailRecord(clientId, monetaryAmount));
            clientExposureDetailRepository.save(record);
        });
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