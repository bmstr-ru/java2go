package ru.bmstr.java2go.listener;

import org.springframework.jms.annotation.JmsListener;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Component;

@Component
public class DealListener {

    @JmsListener(destination = "${jms.deal.queue}")
    public void onMessage(@Payload DealMessage dealMessage) {

    }
}
