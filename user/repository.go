package user

import (
	"context"
	"errors"
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
	r = &Repository{db.GetDB().Collection("user")}
}

func (r *Repository) Add(ctx context.Context, user Model) (Model, error) {
	user.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, user)
	return user, err
}

func (r *Repository) Get(ctx context.Context, id primitive.ObjectID) (Model, error) {
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	return m, err
}

func (r *Repository) GetList(ctx context.Context, id []primitive.ObjectID) (ModelList, error) {
	var m []Model
	find, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": id}})
	if err != nil {
		return nil, err
	}
	err = find.All(ctx, &m)
	return m, err
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (Model, error) {
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"login": login}).Decode(&m)
	return m, err
}

func (r *Repository) GetByAlias(ctx context.Context, alias string) (Model, error) {
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"alias": alias}).Decode(&m)
	if err == mongo.ErrNoDocuments {
		return m, errors.New("user not found")
	}
	return m, err
}

func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repository) Update(ctx context.Context, user Model) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{
		"alias":    user.Alias,
		"login":    user.Name,
		"password": user.Password,
		"name":     user.Name,
	}})
	return err
}

func (r *Repository) UpdateAvatar(ctx context.Context, id primitive.ObjectID, uri string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"image": uri}})
	return err
}

func (r *Repository) CreateIndexes(ctx context.Context) error {
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
