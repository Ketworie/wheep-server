package hub

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wheep-server/db"
)

type Repository struct {
	collection *mongo.Collection
}

func (r *Repository) Add(hub Model) (Model, error) {
	hub.ID = uuid.New()
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, hub)
	return hub, err
}

func (r *Repository) Get(id uuid.UUID) (Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	return m, err
}

func (r *Repository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repository) Update(hub Model) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": hub.ID}, bson.M{"$set": bson.M{"name": hub.Name, "users": hub.Users}})
	return err
}

func (r *Repository) FindByUser(userId uuid.UUID) ([]Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	find, err := r.collection.Find(ctx, bson.M{"users": bson.M{"$in": userId}})
	if err != nil {
		return nil, err
	}
	var hubs []Model
	ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	err = find.All(ctx, &hubs)
	return hubs, err
}
