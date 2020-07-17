package chat

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"wheep-server/hub"
	"wheep-server/message"
)

func Send(mv message.View) (message.Model, error) {
	isMember, err := hub.GetService().IsMember(mv.HubId, mv.UserId)
	if err != nil {
		return message.Model{}, err
	}
	if !isMember {
		return message.Model{}, errors.New("you are not a member of this hub")
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
