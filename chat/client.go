package chat

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
	"wheep-server/hub"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Minute

	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	node       *Node
	connection *websocket.Conn
	send       chan hub.Message
	queue      []hub.Message
	sender     chan bool
}

func (c *Client) read() {
	defer func() {
		c.node.unregister <- c
		c.connection.Close()
	}()
	c.connection.SetReadLimit(maxMessageSize)
	c.connection.SetReadDeadline(time.Now().Add(pongWait))
	c.connection.SetPongHandler(func(string) error { c.connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var m hub.Message
		err := c.connection.ReadJSON(&m)
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				break
			}
			log.Printf("Error during reading JSON from client %v", err)
		}
		c.node.receive <- m
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.connection.Close()
		close(c.sender)
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.connection.SetWriteDeadline(time.Now().Add(writeWait))
				// The hub closed the channel.
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.queue = append(c.queue, message)
			select {
			case c.sender <- true:
				c.sendMessages(c.queue)
				c.queue = make([]hub.Message, 0)
			default:
				//Just proceed to next iteration
			}

		case <-ticker.C:
			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) sendMessages(messages []hub.Message) {
	for _, message := range messages {
		c.connection.SetWriteDeadline(time.Now().Add(writeWait))
		if err := c.connection.WriteJSON(message); err != nil {
			c.connection.Close()
		}
	}
}
