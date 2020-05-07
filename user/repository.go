package user

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

func (r *Repository) Add(user Model) (Model, error) {
	user.ID = uuid.New()
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, user)
	return user, err
}

func (r *Repository) Get(id uuid.UUID) (Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	return m, err
}

func (r *Repository) GetByLogin(login string) (Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"login": login}).Decode(&m)
	return m, err
}

func (r *Repository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repository) Update(user Model) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{
		"alias":    user.Alias,
		"login":    user.Name,
		"password": user.Password,
		"name":     user.Name,
	}})
	return err
}
