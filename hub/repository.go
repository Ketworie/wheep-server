package hub

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
	r = &Repository{db.GetDB().Collection("hub")}
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

func (r *Repository) GetUserIds(id primitive.ObjectID) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"users": 1})).Decode(&m)
	return m.Users, err
}

func (r *Repository) GetList(id []primitive.ObjectID) ([]Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	var m []Model
	find, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": id}})
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
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repository) Rename(id primitive.ObjectID, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": name}})
	return err
}

func (r *Repository) UpdateAvatar(hubId primitive.ObjectID, image string) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": hubId}, bson.M{"$set": bson.M{"image": image}})
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

func (r *Repository) RemoveUser(id primitive.ObjectID, user primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"users": user}})
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
	projection := options.FindOne().SetProjection(bson.M{"_id": 1})
	err := r.collection.FindOne(ctx, bson.M{"_id": hubId, "users": bson.M{"$in": []primitive.ObjectID{userId}}}, projection).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
