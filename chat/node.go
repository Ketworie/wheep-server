package chat

import (
	"wheep-server/hub"
)

type Node struct {
	clients    map[*Client]bool
	receive    chan hub.Message
	register   chan *Client
	unregister chan *Client
	subId      uint16
}

func (n *Node) run() {
	for {
		select {
		case c := <-n.unregister:
			if _, ok := n.clients[c]; ok {
				close(c.send)
				delete(n.clients, c)
			}
		case c := <-n.register:
			n.clients[c] = true
		case m := <-n.receive:
			for c := range n.clients {
				n.subId++
				m.SubId = n.subId
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(n.clients, c)
				}
			}
		}
	}
}
