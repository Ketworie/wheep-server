package message

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Model struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserId primitive.ObjectID `bson:"userId"`
	HubId  primitive.ObjectID `bson:"hubId"`
	Text   string             `bson:"text"`
	Date   time.Time          `bson:"date"`
	NextId primitive.ObjectID `bson:"nextId"`
}

type View struct {
	ID     primitive.ObjectID `json:"id"`
	UserId primitive.ObjectID `json:"userId"`
	HubId  primitive.ObjectID `json:"hubId"`
	Text   string             `json:"text"`
	Date   time.Time          `json:"date"`
	NextId primitive.ObjectID `json:"nextId"`
}

func (m Model) View() View {
	return View{
		ID:     m.ID,
		UserId: m.UserId,
		HubId:  m.HubId,
		Text:   m.Text,
		Date:   m.Date,
		NextId: m.NextId,
	}
}

type ModelList []Model

func (ml ModelList) View() []View {
	vl := []View{}
	for _, model := range ml {
		vl = append(vl, model.View())
	}
	return vl
}
