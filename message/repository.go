package message

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
	r = &Repository{db.GetDB().Collection("message")}
}

func (r *Repository) Add(ctx context.Context, m Model) (Model, error) {
	session, err := r.collection.Database().Client().StartSession()
	defer func() {
		session.EndSession(ctx)
	}()
	if err != nil {
		return m, err
	}
	m.ID = primitive.NewObjectID()
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		var prev Model
		err = r.collection.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"date": -1}).SetProjection(bson.M{"_id": 1})).Decode(&prev)
		if err != nil {
			documents, _ := r.collection.CountDocuments(ctx, nil)
			if documents > 0 {
				return nil, err
			}
		}
		m.PrevId = prev.ID
		return r.collection.InsertOne(ctx, m)
	})
	return m, err
}

func (r *Repository) Prev(ctx context.Context, hubId primitive.ObjectID, time time.Time) (ModelList, error) {
	ms := []Model{}
	find, err := r.collection.Find(ctx, bson.M{"date": bson.M{"$lt": time}, "hubId": bson.M{"$eq": hubId}}, options.Find().SetLimit(30).SetSort(bson.M{"date": -1}))
	if err != nil {
		return nil, err
	}
	err = find.All(ctx, &ms)
	return ms, err
}

func (r *Repository) Next(ctx context.Context, hubId primitive.ObjectID, time time.Time) (ModelList, error) {
	ms := []Model{}
	find, err := r.collection.Find(ctx, bson.M{"date": bson.M{"$gt": time}, "hubId": bson.M{"$eq": hubId}}, options.Find().SetLimit(30).SetSort(bson.M{"date": 1}))
	if err != nil {
		return nil, err
	}
	err = find.All(ctx, &ms)
	return ms, err
}
