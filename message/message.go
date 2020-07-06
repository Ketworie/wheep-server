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
	UserId primitive.ObjectID `json:"userId"`
	HubId  primitive.ObjectID `json:"hubId"`
	Text   string             `json:"text"`
	Date   time.Time          `json:"date"`
	PrevID primitive.ObjectID `bson:"prevId"`
}
