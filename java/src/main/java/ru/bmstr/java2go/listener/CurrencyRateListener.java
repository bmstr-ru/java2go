package ru.bmstr.java2go.listener;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.jms.BytesMessage;
import jakarta.jms.JMSException;
import jakarta.jms.Message;
import jakarta.jms.TextMessage;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.jms.annotation.JmsListener;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Component;
import ru.bmstr.java2go.service.CurrencyRateService;

import java.io.IOError;
import java.io.IOException;
import java.nio.charset.StandardCharsets;

@Slf4j
@Component
@RequiredArgsConstructor
public class CurrencyRateListener {

    private final CurrencyRateService currencyRateService;
    private final ObjectMapper objectMapper;

    @JmsListener(destination = "${jms.rate.queue}")
    public void onMessage(@Payload Message message) throws JMSException {
        String strMessage = switch (message) {
            case TextMessage textMessage -> textMessage.getText();
            case BytesMessage bytesMessage -> readBytesMessage(bytesMessage);
            default -> "";
        };
        log.info("Received rate message: {}", strMessage);
        try {
            CurrencyRateMessage rateMessage = objectMapper.readValue(strMessage, CurrencyRateMessage.class);
            currencyRateService.receiveRate(rateMessage);
        } catch (JsonProcessingException e) {
            log.error("Failed to process rate message: {}", strMessage, e);
        }
    }

    private String readBytesMessage(BytesMessage bytesMessage) throws JMSException {
        byte[] byteData = new byte[(int) bytesMessage.getBodyLength()];
        bytesMessage.readBytes(byteData);
        return new String(byteData, StandardCharsets.UTF_8);
    }

}
