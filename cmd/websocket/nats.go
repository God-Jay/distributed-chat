package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync/atomic"
	"time"
)

var js = NewJS()

type NC struct {
	room *Room
	js   nats.JetStreamContext
	sub  *nats.Subscription
}

// TODO check
func NewJS() nats.JetStreamContext {
	// Connect to NATS
	// TODO change to config
	nc, _ := nats.Connect("nats://god-jay-chat-nats-1:4222")
	//nc, _ := nats.Connect("nats://localhost:4222")

	// TODO MaxWait
	// Create JetStream Context
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(1024), nats.MaxWait(time.Second*10))

	js.AddStream(&nats.StreamConfig{
		Name:     "chat",
		Subjects: []string{"chat.*"},
	})
	return js
}

func (n *NC) publish(msg []byte) {
	// Simple Stream Publisher
	// TODO publish
	_, err := n.js.Publish("chat."+n.room.info.Name, msg)
	if err != nil {
		// TODO nats:timeout
		log.Print("pub.....", pub)
		panic(err)
	}
	atomic.AddInt64(&pub, 1)
}

// TODO
func (n *NC) subscribe() {
	// Simple Async Ephemeral Consumer
	subscribe, err := n.js.Subscribe(fmt.Sprintf("chat.%s", n.room.info.Name), func(m *nats.Msg) {
		n.room.broadcast <- m.Data
	})
	if err != nil {
		// TODO
		panic(err)
	}
	n.sub = subscribe
}

func (n *NC) unsubscribe() {
	n.sub.Unsubscribe()
}
