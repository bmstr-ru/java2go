package activemq

import (
	"encoding/json"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/go-stomp/stomp/v3"
	"github.com/rs/zerolog/log"
	"net"
	"time"
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

func subscribe(url, queue string) (*stomp.Conn, *stomp.Subscription, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, nil, err
	}

	stompConn, err := stomp.Connect(conn, stomp.ConnOpt.HeartBeat(time.Second*3, time.Second*3))
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	sub, err := stompConn.Subscribe(queue, stomp.AckClientIndividual)
	return stompConn, sub, err
}
