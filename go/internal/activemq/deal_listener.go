package activemq

import (
	"encoding/json"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/rs/zerolog/log"
)

func StartDealListener(url, queue string) (<-chan *java2go.Deal, error) {
	conn, sub, err := subscribe(url, queue)
	if err != nil {
		return nil, err
	}

	dealChannel := make(chan (*java2go.Deal))

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

			var deal java2go.Deal
			err = json.Unmarshal(msg.Body, &deal)

			if err != nil {
				log.Warn().Err(err).Msg("Got unparseable deal message: " + string(msg.Body))
			} else {
				dealChannel <- &deal
			}

			err = conn.Ack(msg)
			if err != nil {
				log.Error().Msg("Could not ack deal message")
			}
		}
	}()
	return dealChannel, nil
}
