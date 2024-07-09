package ru.bmstr.java2go.listener;

import com.fasterxml.jackson.core.JacksonException;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.jms.annotation.JmsListener;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Component;
import ru.bmstr.java2go.model.DealDto;

@Component
public class DealListener {

    private static final Logger log = LoggerFactory.getLogger(DealListener.class);

    public DealListener(ObjectMapper objectMapper) {
        this.objectMapper = objectMapper;
    }

    private final ObjectMapper objectMapper;

    @JmsListener(destination = "${jms.deal.queue}")
    public void onMessage(@Payload String message) {
        try {
            DealDto dealDto = objectMapper.readValue(message, DealDto.class);
            log.info("{}", dealDto);
        } catch (JacksonException e) {
            log.error("Invalid message received", e);
        }
    }
}
