package event

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
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
	r = &Repository{db.GetDB().Collection("event")}
}

func (r *Repository) Add(ctx context.Context, m Model) (Model, error) {
	_, err := r.collection.InsertOne(ctx, m)
	return m, err
}

func (r *Repository) Last(ctx context.Context, userId primitive.ObjectID, date time.Time) (Model, error) {
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"userId": userId, "$gt": bson.M{"date": date}}, options.FindOne().SetSort(bson.M{"date": -1})).Decode(&m)
	return m, err
}
