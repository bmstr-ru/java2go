package activemq

import (
	"github.com/go-stomp/stomp/v3"
	"net"
	"time"
)

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
