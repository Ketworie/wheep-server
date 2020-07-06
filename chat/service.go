package chat

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"wheep-server/hub"
	"wheep-server/message"
)

func Send(mv message.View) error {
	isMember, err := hub.GetService().IsMember(mv.HubId, mv.UserId)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member of this hub")
	}
	_, err = message.GetService().Add(message.Model{
		ID:     primitive.ObjectID{},
		UserId: mv.UserId,
		HubId:  mv.HubId,
		Text:   mv.Text,
		Date:   time.Now(),
		NextId: primitive.ObjectID{},
	})
	return err
}
