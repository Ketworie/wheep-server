package message

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	wheepTime "wheep-server/time"
)

type Model struct {
	ID     primitive.ObjectID   `bson:"_id"`
	UserId primitive.ObjectID   `bson:"userId"`
	HubId  primitive.ObjectID   `bson:"hubId"`
	Text   string               `bson:"text"`
	Date   time.Time            `bson:"date"`
	PrevId primitive.ObjectID   `bson:"prevId"`
	ReadBy []primitive.ObjectID `bson:"readBy"`
}

type View struct {
	ID     primitive.ObjectID `json:"id"`
	UserId primitive.ObjectID `json:"userId"`
	HubId  primitive.ObjectID `json:"hubId"`
	Text   string             `json:"text"`
	Date   wheepTime.JSONTime `json:"date"`
	PrevId primitive.ObjectID `json:"prevId"`
	//ReadBy []primitive.ObjectID `json:"readBy"`
}

func (m Model) View() View {
	return View{
		ID:     m.ID,
		UserId: m.UserId,
		HubId:  m.HubId,
		Text:   m.Text,
		Date:   wheepTime.JSONTime{m.Date},
		PrevId: m.PrevId,
		//ReadBy: m.ReadBy,
	}
}

func (v View) Type() string {
	return "message"
}

type ModelList []Model

func (ml ModelList) View() []View {
	vl := []View{}
	for _, model := range ml {
		vl = append(vl, model.View())
	}
	return vl
}
