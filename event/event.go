package event

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Typer interface {
	Type() string
}

type View struct {
	Date time.Time `json:"date"`
	Type string    `json:"type"`
	Body string    `json:"body"`
}

type Model struct {
	Date   time.Time          `bson:"date"`
	Type   string             `bson:"type"`
	Body   string             `bson:"body"`
	UserId primitive.ObjectID `bson:"userId"`
}

func (e Model) View() View {
	return View{
		Date: e.Date,
		Type: e.Type,
		Body: e.Body,
	}
}

func NewEventModel(body Typer, recipient primitive.ObjectID) (Model, error) {
	marshal, err := json.Marshal(body)
	return Model{
		UserId: recipient,
		Date:   time.Now(),
		Type:   body.Type(),
		Body:   string(marshal),
	}, err
}
