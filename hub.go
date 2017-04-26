/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import(
    "encoding/json"
    "log"
)

var _ = json.Marshal
var _ = log.Println

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

    history [][]byte
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
    h.history = append(h.history, msg)
    for client := range h.clients {
        select {
        case client.send <- msg:
            
        default:
            close(client.send)
            delete(h.clients, client)
        }
    }
}

func (h *Hub) run() {
	for {
		select {

		case client := <-h.register:
			h.clients[client] = true
            log.Println("REGISTER", client.Username)

        case client := <-h.unregister:
            log.Println("UNREGISTER", client.Username)
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
            temp, _ := json.Marshal(NewChannelLeave(client.Username,
                client.Channel))
            h.gobroadcast(temp)
            //h.broadcast <- temp
            
		case message := <-h.broadcast:
            h.gobroadcast(message)

		}
	}
}
