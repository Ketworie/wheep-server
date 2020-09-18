package notebook

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"wheep-server/db"
)

type Repository struct {
	collection *mongo.Collection
}

var r *Repository
var rOnce sync.Once

func GetRepository() *Repository {
	rOnce.Do(initRepository)
	return r
}

func initRepository() {
	r = &Repository{db.GetDB().Collection("notebook")}
}

func (r *Repository) GetContacts(ctx context.Context, userId primitive.ObjectID) ([]primitive.ObjectID, error) {
	var n Model
	cs := []primitive.ObjectID{}
	err := r.collection.FindOne(ctx, bson.M{"_id": userId}, options.FindOne().SetProjection(bson.M{"contacts": 1})).Decode(&n)
	return append(cs, n.Contacts...), err
}

func (r *Repository) AddContact(ctx context.Context, userId primitive.ObjectID, contact primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$addToSet": bson.M{"contacts": contact}}, options.Update().SetUpsert(true))
	return err
}

func (r *Repository) RemoveContact(ctx context.Context, userId primitive.ObjectID, contact primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$pullAll": bson.M{"contacts": contact}})
	return err
}
