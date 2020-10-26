package event

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	wheepTime "wheep-server/time"
)

type Typer interface {
	Type() string
}

type View struct {
	Date wheepTime.JSONTime `json:"date"`
	Type string             `json:"type"`
	Body string             `json:"body"`
}

func (v View) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{\"date\": \"%v\", \"type\": \"%v\", \"body\": %v}", v.Date.UTC().Format(wheepTime.Zoned), v.Type, v.Body)
	return []byte(s), nil
}

type Model struct {
	Date   time.Time          `bson:"date"`
	Type   string             `bson:"type"`
	Body   string             `bson:"body"`
	UserId primitive.ObjectID `bson:"userId"`
}

func (m Model) View() View {
	return View{
		Date: wheepTime.JSONTime{m.Date},
		Type: m.Type,
		Body: m.Body,
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
