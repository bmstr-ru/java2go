package activemq

import (
	"encoding/json"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/rs/zerolog/log"
)

func StartRateListener(url, queue string) (<-chan *java2go.CurrencyRate, error) {
	conn, sub, err := subscribe(url, queue)
	if err != nil {
		return nil, err
	}

	rateChannel := make(chan (*java2go.CurrencyRate))

	go func() {
		for {
			msg, err := sub.Read()
			if err != nil {
				log.Err(err)
				conn.Disconnect()
				conn, sub, err = subscribe(url, queue)
				continue
			}

			log.Info().Msg("Received message " + msg.ContentType)
			log.Debug().Msg(string(msg.Body))

			var rates java2go.CurrencyRateMessage
			err = json.Unmarshal(msg.Body, &rates)

			if err != nil {
				log.Warn().Err(err).Msg("Got unparseable rate message: " + string(msg.Body))
			} else {
				for _, rate := range rates {
					rateChannel <- &rate
				}
			}

			err = conn.Ack(msg)
			if err != nil {
				log.Error().Msg("Could not ack deal message")
			}
		}
	}()
	return rateChannel, nil
}
