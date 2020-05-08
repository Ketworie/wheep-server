package user

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

func (r *Repository) Add(user Model) (Model, error) {
	user.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, user)
	return user, err
}

func (r *Repository) Get(id primitive.ObjectID) (Model, error) {
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

func (r *Repository) Delete(id primitive.ObjectID) error {
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

func (r *Repository) CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	login := mongo.IndexModel{
		Keys: bson.M{"login": 1},
		Options: &options.IndexOptions{
			Unique: &[]bool{true}[0],
		},
	}
	alias := mongo.IndexModel{
		Keys: bson.M{"alias": 1},
		Options: &options.IndexOptions{
			Unique: &[]bool{true}[0],
		},
	}
	indexes := []mongo.IndexModel{login, alias}
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
