package hub

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"wheep-server/db"
)

type Repository struct {
	collection *mongo.Collection
}

func (r *Repository) Add(hub Model) (Model, error) {
	hub.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, hub)
	return hub, err
}

func (r *Repository) Get(id primitive.ObjectID) (Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	return m, err
}

func (r *Repository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repository) Rename(hub Model) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": hub.ID}, bson.M{"$set": bson.M{"name": hub.Name}})
	return err
}

func (r *Repository) FindByUser(userId primitive.ObjectID) ([]Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	find, err := r.collection.Find(ctx, bson.M{"users": bson.M{"$in": []primitive.ObjectID{userId}}})
	if err != nil {
		return nil, err
	}
	var hubs []Model
	ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	err = find.All(ctx, &hubs)
	return hubs, err
}

func (r *Repository) AddUsers(id primitive.ObjectID, users []primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"users": bson.M{"$each": users}}})
	return err
}

func (r *Repository) RemoveUsers(id primitive.ObjectID, users []primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"users": bson.M{"$in": users}}})
	return err
}
