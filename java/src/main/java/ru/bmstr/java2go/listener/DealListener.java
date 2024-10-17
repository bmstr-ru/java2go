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
import ru.bmstr.java2go.service.DealService;

@Slf4j
@Component
@RequiredArgsConstructor
public class DealListener {

    private final DealService dealService;
    private final ObjectMapper objectMapper;

    @JmsListener(destination = "${jms.deal.queue}")
    public void onMessage(@Payload TextMessage message) throws JMSException {
        log.info("Received deal message: {}", message.getText());
        try {
            DealMessage dealMessage = objectMapper.readValue(message.getText(), DealMessage.class);
            dealService.receiveDeal(dealMessage);
        } catch (JsonProcessingException e) {
            log.error("Failed to process deal: {}", message.getText(), e);
        }
    }

}
