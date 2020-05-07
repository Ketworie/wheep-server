package hub

import (
	"github.com/google/uuid"
	"time"
)

type Model struct {
	ID    uuid.UUID          `bson:"_id"`
	Name  string             `bson:"name"`
	Users map[uuid.UUID]bool `bson:"users"`
}

type View struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UserCount int       `json:"userCount"`
}

type Message struct {
	ID     uuid.UUID `bson:"_id"`
	UserId uuid.UUID `bson:"userId"`
	HubId  uuid.UUID `bson:"hubId"`
	Text   string    `bson:"text"`
	Date   time.Time `bson:"date"`
}
