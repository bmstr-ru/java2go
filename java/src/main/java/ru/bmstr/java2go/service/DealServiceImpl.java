package ru.bmstr.java2go.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.bmstr.java2go.listener.Deal;
import ru.bmstr.java2go.persistence.DealRecord;
import ru.bmstr.java2go.persistence.DealRepository;

@Slf4j
@Service
@RequiredArgsConstructor
public class DealServiceImpl implements DealService {

    private final DealRepository dealRepository;
    private final TotalExposureService totalExposureService;

    @Override
    @Transactional
    public void receiveDeal(Deal deal) {
        DealRecord record = DealRecord.builder()
                .dealId(deal.id())
                .clientId(deal.clientId())
                .boughtAmount(deal.amountBought().amount())
                .boughtCurrency(deal.amountBought().currency())
                .soldAmount(deal.amountSold().amount())
                .soldCurrency(deal.amountSold().currency())
                .build();
        record = dealRepository.save(record);
        log.info("New deal record saved: id={}", record.getId());

        totalExposureService.considerNewAmounts(deal.clientId(), deal.amountBought(), deal.amountSold().negate());
    }
}
