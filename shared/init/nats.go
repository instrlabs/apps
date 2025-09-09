package initx

import (
	"time"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	Conn *nats.Conn
}

func NewNats(url string) (*Nats, error) {
	c, err := nats.Connect(url, nats.Timeout(10*time.Second))
	if err != nil {
		return nil, err
	}
	return &Nats{Conn: c}, nil
}

func (n *Nats) Close() {
	if n.Conn != nil {
		n.Conn.Close()
	}
}
