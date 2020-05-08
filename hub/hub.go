package hub

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Model struct {
	ID    primitive.ObjectID   `bson:"_id"`
	Name  string               `bson:"name"`
	Users []primitive.ObjectID `bson:"users"`
}

func (h Model) View() View {
	return View{
		ID:        h.ID,
		Name:      h.Name,
		UserCount: len(h.Users),
	}
}

type View struct {
	ID        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	UserCount int                `json:"userCount"`
}

type AddView struct {
	Name  string               `json:"name"`
	Users []primitive.ObjectID `json:"users"`
}

type Message struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserId primitive.ObjectID `bson:"userId"`
	HubId  primitive.ObjectID `bson:"hubId"`
	Text   string             `bson:"text"`
	Date   time.Time          `bson:"date"`
}
