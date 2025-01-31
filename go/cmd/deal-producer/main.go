package main

import (
	"encoding/json"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/go-stomp/stomp/v3"
	"github.com/rs/zerolog/log"
	"math"
	"math/rand"
	"net"
)

import (
	"time"
)

func main() {
	conn, err := establish_amq_connection("localhost:61613")
	if err != nil {
		panic(err)
	}
	go start_rates_producer(conn)
	go start_deal_producer(conn)
	select {}
}

func start_deal_producer(conn *stomp.Conn) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		deal := java2go.Deal{
			Id:       time.Now().Unix(),
			ClientId: random.Int63n(10) + 1,
			AmountBought: java2go.MonetaryAmount{
				Currency: "USD",
				Amount:   math.Round(rand.Float64() * 1000000),
			},
			AmountSold: java2go.MonetaryAmount{
				Currency: "EUR",
				Amount:   math.Round(rand.Float64() * 1000000),
			},
		}
		message, err := json.Marshal(&deal)
		err = conn.Send("DEAL.QUEUE", "text/plain", message)
		if err != nil {
			log.Err(err).Msg("Could not send deal message")
		}
		log.Info().Any("message", string(message)).Msg("Sent deal message")
		time.Sleep(300 * time.Millisecond)
	}
}

func start_rates_producer(conn *stomp.Conn) {
	for {
		rateMessage := java2go.CurrencyRateMessage{
			{
				BaseCurrency:   "EUR",
				QuotedCurrency: "USD",
				Rate:           0.9,
			},
		}
		message, err := json.Marshal(&rateMessage)
		err = conn.Send("RATE.QUEUE", "text/plain", message)
		if err != nil {
			log.Err(err).Msg("Could not send rates message")
		}
		log.Info().Any("message", string(message)).Msg("Sent rate message")
		time.Sleep(10 * time.Second)
	}
}

func establish_amq_connection(url string) (*stomp.Conn, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}

	stompConn, err := stomp.Connect(conn, stomp.ConnOpt.HeartBeat(time.Second*3, time.Second*3))
	if err != nil {
		conn.Close()
		return nil, err
	}

	return stompConn, nil
}
