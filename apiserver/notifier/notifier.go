package notifier

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Notifier struct {
	eventq  chan interface{}
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewNotifier() *Notifier {
	return &Notifier{
		eventq:  make(chan interface{}, 300),
		clients: make(map[*websocket.Conn]bool),
		mu:      sync.Mutex{},
	}
}

func (n *Notifier) AddClient(client *websocket.Conn) {
	n.mu.Lock()
	n.clients[client] = true
	n.mu.Unlock()
	go n.readPump(client)

}

func (n *Notifier) Notify(event interface{}) {
	n.eventq <- event
}

func (n *Notifier) readPump(client *websocket.Conn) {
	for {
		if _, _, err := client.NextReader(); err != nil {
			client.Close()
			break
		}
	}
}

func (n *Notifier) Start() {
	for {
		select {
		case event := <-n.eventq:
			n.broadcast(event)
		}
	}
}

func (n *Notifier) broadcast(event interface{}) {
	for conn, _ := range n.clients {
		err := conn.WriteJSON(event)
		if err != nil {
			conn.Close()
			delete(n.clients, conn)
		}
	}
}
