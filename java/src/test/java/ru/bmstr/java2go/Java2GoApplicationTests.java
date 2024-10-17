package ru.bmstr.java2go;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.SneakyThrows;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
import org.springframework.jms.core.JmsTemplate;
import org.springframework.test.web.servlet.MockMvc;
import ru.bmstr.java2go.listener.CurrencyRateMessage;
import ru.bmstr.java2go.listener.DealMessage;
import ru.bmstr.java2go.model.MonetaryAmount;
import ru.bmstr.java2go.persistence.ClientExposureRepository;
import ru.bmstr.java2go.persistence.CurrencyRateRepository;

import java.math.BigDecimal;
import java.time.Duration;

import static org.awaitility.Awaitility.await;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@Import(TestcontainersConfiguration.class)
@SpringBootTest
@AutoConfigureMockMvc
class Java2GoApplicationTests {

    private static final Duration FIVE_SECONDS = Duration.ofSeconds(5);
    private static final Long CLIENT_ID = 53L;
    private static final MonetaryAmount USD_BOUGHT = new MonetaryAmount(new BigDecimal("871982.48"), "USD");
    private static final MonetaryAmount EUR_SOLD = new MonetaryAmount(new BigDecimal("800000.0"), "EUR");
    private static final MonetaryAmount EUR_TOTAL = new MonetaryAmount(new BigDecimal("-6061.6589"), "EUR");

    @Autowired
    JmsTemplate jmsTemplate;

    @Autowired
    ObjectMapper objectMapper;

    @Autowired
    CurrencyRateRepository rateRepository;

    @Autowired
    ClientExposureRepository exposureRepository;

    @Autowired
    MockMvc mockMvc;

    @Value("${jms.deal.queue}")
    private String dealQueue;

    @Value("${jms.rate.queue}")
    private String rateQueue;

    @Test
    void integrationTest() {
        sendRates();
        awaitRatesStored();

        sendDeal();
        awaitDealProcessed();

        verifyClientExposure();
    }

    @SneakyThrows
    private void verifyClientExposure() {
        mockMvc.perform(get("/client/" + CLIENT_ID + "/summary"))
                .andDo(print())
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.clientId").value(CLIENT_ID))
                .andExpect(jsonPath("$.total.amount").value(EUR_TOTAL.amount()))
                .andExpect(jsonPath("$.total.currency").value(EUR_TOTAL.currency()))
                .andExpect(jsonPath("$.amounts[0].amount").value(EUR_SOLD.amount().negate()))
                .andExpect(jsonPath("$.amounts[0].currency").value(EUR_SOLD.currency()))
                .andExpect(jsonPath("$.amounts[1].amount").value(USD_BOUGHT.amount()))
                .andExpect(jsonPath("$.amounts[1].currency").value(USD_BOUGHT.currency()));
    }

    @SneakyThrows
    private void sendRates() {
        CurrencyRateMessage rateMessage = new CurrencyRateMessage();
        rateMessage.add(new CurrencyRateMessage.Rate("EUR", "USD", new BigDecimal("1.0983")));
        rateMessage.add(new CurrencyRateMessage.Rate("EUR", "GBP", new BigDecimal("0.845675")));
        rateMessage.add(new CurrencyRateMessage.Rate("EUR", "CHF", new BigDecimal("0.970912")));
        String message = objectMapper.writeValueAsString(rateMessage);
        jmsTemplate.send(rateQueue, s -> s.createTextMessage(message));
    }

    private void awaitRatesStored() {
        await().atMost(FIVE_SECONDS)
                .until(rateRepository::count, c -> c > 0);
    }

    @SneakyThrows
    private void sendDeal() {
        DealMessage dealMessage = DealMessage.builder()
                .id(7L)
                .clientId(CLIENT_ID)
                .amountBought(USD_BOUGHT)
                .amountSold(EUR_SOLD)
                .build();
        String message = objectMapper.writeValueAsString(dealMessage);
        jmsTemplate.send(dealQueue, s -> s.createTextMessage(message));
    }

    private void awaitDealProcessed() {
        await().atMost(FIVE_SECONDS)
                .until(exposureRepository::count, c -> c > 0);
    }
}
