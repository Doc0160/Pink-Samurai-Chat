/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

    "encoding/json"

	"github.com/gorilla/websocket"
)

var _ = bytes.TrimSpace
var _ = log.Println

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is an middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

    Username string

    Channel string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
        c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

    temp, _ := json.Marshal(NewHello())
    c.send <- temp

    for _, t := range c.hub.history {
        c.send <- t
    }
    
    temp, _ = json.Marshal(NewChannelJoin(c.Username, c.Channel))
    c.hub.broadcast <- temp
    
    for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		//message = bytes.TrimSpace(message)

        stub := Stub{}
        json.Unmarshal(message, &stub)
        if stub.Type == _ChannelJoin {
            uc := ChannelJoin{}
            json.Unmarshal(message, &uc)
            c.Channel = uc.Channel
            uc.Username = c.Username
            c.hub.broadcast <- message
            for _, t := range c.hub.history {
                select {
                    case _, ok := <- c.send:
                    if ok != true {
                        break
                    }
                    
                    case c.send <- t:
                }
            }
            continue
            
        } else if stub.Type == _Command {
            uc := Command{}
            json.Unmarshal(message, &uc)
            switch uc.Command {
            case "/tentacules":
                m := NewMessage("Tentacule-Sama", uc.Channel, "I bet you'd love my tentacules! <3")
                b, _ := json.Marshal(m)
                c.hub.broadcast <- b
                
            default:
                log.Println(uc)
                m := NewMessage("Tentacule-Sama", uc.Channel, "I don't understand you.")
                b, _ := json.Marshal(m)
                c.send <- b
            }
            continue
            
        } else {
            var uc map[string]interface{}
            json.Unmarshal(message, &uc)
            uc["username"] = c.Username
            uc["time"] = time.Now().UnixNano()
            b, _ := json.Marshal(uc)
            c.hub.broadcast <- b
        }
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
                temp, _ := json.Marshal(NewDisconnect(c.Username))
                c.hub.broadcast <- temp
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
    cookie, _ := r.Cookie("session")
	client := &Client{
        hub: hub,
        conn: conn,
        send: make(chan []byte, 256),
        Username: members.hashs[cookie.Value].Username,
    }    
    
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
