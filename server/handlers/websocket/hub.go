package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type TransportMsg struct {
	Data interface{}
}

type Hub struct {
	broadcast chan uint

	clientMap map[uint](map[*websocket.Conn]bool)

	lock sync.Mutex
}

var hub *Hub = &Hub{
	make(chan uint), make(map[uint]map[*websocket.Conn]bool), sync.Mutex{},
}

func GetHub() *Hub {
	return hub
}

func (hub *Hub) RefreshCart(userId uint) {
	hub.broadcast <- userId
}

func (hub *Hub) RegisterNewClient(userId uint, ws *websocket.Conn) {
	hub.lock.Lock()
	clients, ok := hub.clientMap[userId]
	if !ok {
		clients = make(map[*websocket.Conn]bool)
	}

	clients[ws] = true
	hub.clientMap[userId] = clients
	hub.lock.Unlock()
}

func (hub *Hub) Run() {
	for userId := range hub.broadcast {
		hub.lock.Lock()
		clients, ok := hub.clientMap[userId]
		hub.lock.Unlock()
		fmt.Println(clients)
		if !ok {
			continue
		}

		toRevoke := []*websocket.Conn{}
		for client := range clients {
			err := client.WriteJSON("reload")
			if err != nil {
				toRevoke = append(toRevoke, client)
			}
		}

		hub.lock.Lock()
		for _, client := range toRevoke {
			client.Close()
			delete(clients, client)
		}
		hub.lock.Unlock()
	}
}
