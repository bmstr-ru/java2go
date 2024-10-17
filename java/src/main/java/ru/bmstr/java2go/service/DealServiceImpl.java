package ru.bmstr.java2go.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.bmstr.java2go.listener.DealMessage;
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
    public void receiveDeal(DealMessage dealMessage) {
        DealRecord record = DealRecord.builder()
                .dealId(dealMessage.id())
                .clientId(dealMessage.clientId())
                .boughtAmount(dealMessage.amountBought().amount())
                .boughtCurrency(dealMessage.amountBought().currency())
                .soldAmount(dealMessage.amountSold().amount())
                .soldCurrency(dealMessage.amountSold().currency())
                .build();
        record = dealRepository.save(record);
        log.info("New deal record saved: id={}", record.getId());

        totalExposureService.considerNewAmounts(dealMessage.clientId(), dealMessage.amountBought(), dealMessage.amountSold().negate());
    }
}
