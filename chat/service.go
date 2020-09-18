package chat

import (
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type Service struct {
	mqChan       *amqp.Channel
	hubs         *hubSync
	userActivity map[string]*timeSync
}

type timeSync struct {
	sync.RWMutex
	time time.Time
}

type idSync struct {
	sync.RWMutex
	ids []string
}

type hubSync struct {
	sync.RWMutex
	hubs map[primitive.ObjectID]idSync
}
