/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import(
    "encoding/json"
    //"fmt"
)

var _ = json.Marshal

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) gobroadcast(msg []byte){
    temp = append(temp, msg)
    for client := range h.clients {
        select {
        case client.send <- msg:
            
        default:
            close(client.send)
            delete(h.clients, client)
        }
    }
}

var temp [][]byte

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
            for _, t := range temp {
                client.send <- t
            }
            /*temp, _ := json.Marshal(Message{Type: _Message, Username:"Someone", Text:"a challenger has entered"})
            h.gobroadcast(temp)*/

        case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
            temp, _ := json.Marshal(Disconnect{Type: _Disconnect, Username:client.Username})
            h.gobroadcast(temp)
            
		case message := <-h.broadcast:
            h.gobroadcast(message)
		}
	}
}
