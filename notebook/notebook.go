package notebook

import "go.mongodb.org/mongo-driver/bson/primitive"

type Model struct {
	ID       primitive.ObjectID   `bson:"_id"`
	Contacts []primitive.ObjectID `bson:"contacts"`
}
