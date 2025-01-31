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
import ru.bmstr.java2go.service.DealService;

import java.nio.charset.StandardCharsets;

@Slf4j
@Component
@RequiredArgsConstructor
public class DealListener {

    private final DealService dealService;
    private final ObjectMapper objectMapper;

    @JmsListener(destination = "${jms.deal.queue}")
    public void onMessage(@Payload Message message) throws JMSException {
        String strMessage = switch (message) {
            case TextMessage textMessage -> textMessage.getText();
            case BytesMessage bytesMessage -> readBytesMessage(bytesMessage);
            default -> "";
        };
        log.info("Received deal message: {}", strMessage);
        try {
            DealMessage dealMessage = objectMapper.readValue(strMessage, DealMessage.class);
            dealService.receiveDeal(dealMessage);
        } catch (JsonProcessingException e) {
            log.error("Failed to process deal: {}", strMessage, e);
        }
    }

    private String readBytesMessage(BytesMessage bytesMessage) throws JMSException {
        byte[] byteData = new byte[(int) bytesMessage.getBodyLength()];
        bytesMessage.readBytes(byteData);
        return new String(byteData, StandardCharsets.UTF_8);
    }

}
