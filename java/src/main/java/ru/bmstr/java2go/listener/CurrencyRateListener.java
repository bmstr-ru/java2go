package ru.bmstr.java2go.listener;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.jms.JMSException;
import jakarta.jms.TextMessage;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.jms.annotation.JmsListener;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Component;
import ru.bmstr.java2go.service.CurrencyRateService;

@Slf4j
@Component
@RequiredArgsConstructor
public class CurrencyRateListener {

    private final CurrencyRateService currencyRateService;
    private final ObjectMapper objectMapper;

    @JmsListener(destination = "${jms.rate.queue}")
    public void onMessage(@Payload TextMessage message) throws JMSException {
        log.info("Received rate message: {}", message.getText());
        try {
            CurrencyRateMessage rateMessage = objectMapper.readValue(message.getText(), CurrencyRateMessage.class);
            currencyRateService.receiveRate(rateMessage);
        } catch (JsonProcessingException e) {
            log.error("Failed to process rate message: {}", message.getText(), e);
        }
    }

}
