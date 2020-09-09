package chat

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"wheep-server/hub"
	"wheep-server/message"
)

func Send(mv message.View) (message.Model, error) {
	err := hub.GetService().AssertMember(mv.HubId, mv.UserId)
	if err != nil {
		return message.Model{}, err
	}
	model, err := message.GetService().Add(message.Model{
		ID:     primitive.ObjectID{},
		UserId: mv.UserId,
		HubId:  mv.HubId,
		Text:   mv.Text,
		Date:   time.Now(),
		PrevId: primitive.ObjectID{},
	})
	return model, err
}
