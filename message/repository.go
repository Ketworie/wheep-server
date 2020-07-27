package message

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

func (r *Repository) Add(m Model) (Model, error) {
	session, err := r.collection.Database().Client().StartSession()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
		defer cancel()
		session.EndSession(ctx)
	}()
	if err != nil {
		return m, err
	}
	m.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
		defer cancel()
		var prev Model
		err = r.collection.FindOne(ctx, db.M{}, options.FindOne().SetSort(db.M{"date": -1}).SetProjection(db.M{"_id": 1})).Decode(&prev)
		if err != nil {
			ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
			defer cancel()
			documents, _ := r.collection.CountDocuments(ctx, nil)
			if documents > 0 {
				return nil, err
			}
		}
		m.PrevId = prev.ID
		ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
		defer cancel()
		return r.collection.InsertOne(ctx, m)
	})
	return m, err
}

func (r *Repository) Last(hubId primitive.ObjectID) (Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, db.M{"hubId": hubId}, options.FindOne().SetSort(db.M{"date": -1})).Decode(&m)
	if err != nil {
		return Model{}, err
	}
	return m, nil
}

func (r *Repository) Prev(hubId primitive.ObjectID, time time.Time) (ModelList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	ms := []Model{}
	find, err := r.collection.Find(ctx, db.M{"date": db.M{"$lt": time}, "hubId": db.M{"$eq": hubId}}, options.Find().SetLimit(30).SetSort(db.M{"date": -1}))
	if err != nil {
		return nil, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	err = find.All(ctx, &ms)
	return ms, err
}

func (r *Repository) Next(hubId primitive.ObjectID, time time.Time) (ModelList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	ms := []Model{}
	find, err := r.collection.Find(ctx, db.M{"date": db.M{"$gt": time}, "hubId": db.M{"$eq": hubId}}, options.Find().SetLimit(30).SetSort(db.M{"date": 1}))
	if err != nil {
		return nil, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	err = find.All(ctx, &ms)
	return ms, err
}
