package ru.bmstr.java2go;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
import org.springframework.jms.core.JmsTemplate;
import org.springframework.test.context.ActiveProfiles;
import ru.bmstr.java2go.model.DealDto;

import java.math.BigDecimal;

@Import(TestcontainersConfiguration.class)
@SpringBootTest
@ActiveProfiles("test")
class Java2GoApplicationTests {

    @Autowired
    ObjectMapper objectMapper;

    @Autowired
    JmsTemplate jmsTemplate;

    @Value("${jms.deal.queue}")
    String dealQueue;

    @Test
    void contextLoads() throws Exception {
        DealDto dealDto = new DealDto(1L, 2L,
                new DealDto.MonetaryAmount("USD", BigDecimal.ONE),
                new DealDto.MonetaryAmount("EUR", BigDecimal.TEN)
                );
        String message = objectMapper.writeValueAsString(dealDto);

        jmsTemplate.send(dealQueue, session -> session.createTextMessage(message));

        Thread.sleep(5_000);
    }

}
