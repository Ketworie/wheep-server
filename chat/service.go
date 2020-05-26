package chat

import (
	"errors"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"wheep-server/hub"
)

var nodes = sync.Map{}

func connect(userId primitive.ObjectID, hubId primitive.ObjectID, conn *websocket.Conn) error {
	isMember, err := hub.GetService().IsMember(hubId, userId)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member of this hub")
	}

	n, ok := nodes.LoadOrStore(
		hubId,
		&Node{
			clients:    make(map[*Client]bool),
			receive:    make(chan hub.Message),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		})

	node := n.(*Node)

	if !ok {
		node.subId = 1 // TODO
		node.run()
	}

	node.register <- &Client{
		node:       node,
		connection: conn,
		send:       make(chan hub.Message),
		queue:      make([]hub.Message, 0),
		sender:     make(chan bool),
	}

	return nil
}
