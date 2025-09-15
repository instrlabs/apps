package initx

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	Conn *nats.Conn
}

func NewNats(url string) *Nats {
	c, err := nats.Connect(url, nats.Timeout(10*time.Second))
	if err != nil {
		log.Printf("failed to connect to NATS: %v", err)
		return nil
	}
	return &Nats{Conn: c}
}

func (n *Nats) Close() {
	if n.Conn != nil {
		n.Conn.Close()
	}
}
