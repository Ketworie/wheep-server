package hub

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"wheep-server/db"
)

type Repository struct {
	collection *mongo.Collection
}

func (r *Repository) Add(hub Model) (Model, error) {
	hub.ID = primitive.NewObjectID()
	hub.LastModified = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, hub)
	return hub, err
}

func (r *Repository) Get(id primitive.ObjectID) (Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, db.M{"_id": id}).Decode(&m)
	return m, err
}

func (r *Repository) GetList(id []primitive.ObjectID) ([]Model, error) {
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

func (r *Repository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, db.M{"_id": id})
	return err
}

func (r *Repository) Rename(id primitive.ObjectID, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, db.M{"_id": id}, db.M{"$set": db.M{"name": name}}.LastModified())
	return err
}

func (r *Repository) UpdateAvatar(hubId primitive.ObjectID, image string) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, db.M{"_id": hubId}, db.M{"$set": db.M{"image": image}}.LastModified())
	return err
}

func (r *Repository) FindByUser(userId primitive.ObjectID) ([]Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	find, err := r.collection.Find(ctx, db.M{"users": db.M{"$in": []primitive.ObjectID{userId}}})
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
	_, err := r.collection.UpdateOne(ctx, db.M{"_id": id}, db.M{"$addToSet": db.M{"users": db.M{"$each": users}}}.LastModified())
	return err
}

func (r *Repository) RemoveUser(id primitive.ObjectID, user primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, db.M{"_id": id}, db.M{"$pull": db.M{"users": user}}.LastModified())
	return err
}

func (r Repository) AssertMember(hubId primitive.ObjectID, userId primitive.ObjectID) error {
	isMember, err := r.IsMember(hubId, userId)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member of this hub")
	}
	return nil
}

func (r *Repository) IsMember(hubId primitive.ObjectID, userId primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	projection := options.FindOne().SetProjection(db.M{"_id": 1})
	err := r.collection.FindOne(ctx, db.M{"_id": hubId, "users": db.M{"$in": []primitive.ObjectID{userId}}}, projection).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
