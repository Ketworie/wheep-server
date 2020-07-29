package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"wheep-server/db"
)

type Repository struct {
	collection *mongo.Collection
}

func (r *Repository) Add(user Model) (Model, error) {
	user.ID = primitive.NewObjectID()
	user.LastModified = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, user)
	return user, err
}

func (r *Repository) Get(id primitive.ObjectID) (Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, db.M{"_id": id}).Decode(&m)
	return m, err
}

func (r *Repository) GetList(id []primitive.ObjectID) (ModelList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m []Model
	find, err := r.collection.Find(ctx, db.M{"_id": db.M{"$in": id}})
	if err != nil {
		return nil, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	err = find.All(ctx, &m)
	return m, err
}

func (r *Repository) GetByLogin(login string) (Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, db.M{"login": login}).Decode(&m)
	return m, err
}

func (r *Repository) GetByAlias(alias string) (Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, db.M{"alias": alias}).Decode(&m)
	return m, err
}

func (r *Repository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, db.M{"_id": id})
	return err
}

func (r *Repository) Update(user Model) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, db.M{"_id": user.ID}, db.M{"$set": db.M{
		"alias":    user.Alias,
		"login":    user.Name,
		"password": user.Password,
		"name":     user.Name,
	}}.LastModified())
	return err
}

func (r *Repository) CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	login := mongo.IndexModel{
		Keys: db.M{"login": 1},
		Options: &options.IndexOptions{
			Unique: &[]bool{true}[0],
		},
	}
	alias := mongo.IndexModel{
		Keys: db.M{"alias": 1},
		Options: &options.IndexOptions{
			Unique: &[]bool{true}[0],
		},
	}
	indexes := []mongo.IndexModel{login, alias}
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
