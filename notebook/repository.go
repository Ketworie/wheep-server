package notebook

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"wheep-server/db"
)

type Repository struct {
	collection *mongo.Collection
}

func (r *Repository) GetContacts(userId primitive.ObjectID) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var n Model
	cs := []primitive.ObjectID{}
	err := r.collection.FindOne(ctx, bson.M{"_id": userId}, options.FindOne().SetProjection(bson.M{"contacts": 1})).Decode(&n)
	return append(cs, n.Contacts...), err
}

func (r *Repository) AddContact(userId primitive.ObjectID, contact primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$addToSet": bson.M{"contacts": contact}}, options.Update().SetUpsert(true))
	return err
}

func (r *Repository) RemoveContact(userId primitive.ObjectID, contact primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$pullAll": bson.M{"contacts": contact}})
	return err
}
