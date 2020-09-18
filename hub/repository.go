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

func (r *Repository) Add(ctx context.Context, hub Model) (Model, error) {
	hub.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, hub)
	return hub, err
}

func (r *Repository) Get(ctx context.Context, id primitive.ObjectID) (Model, error) {
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	return m, err
}

func (r *Repository) GetUserIds(ctx context.Context, id primitive.ObjectID) ([]primitive.ObjectID, error) {
	var m Model
	err := r.collection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"users": 1})).Decode(&m)
	return m.Users, err
}

func (r *Repository) GetList(ctx context.Context, id []primitive.ObjectID) ([]Model, error) {
	var m []Model
	find, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": id}})
	if err != nil {
		return nil, err
	}
	err = find.All(ctx, &m)
	return m, err
}

func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repository) Rename(ctx context.Context, id primitive.ObjectID, name string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": name}})
	return err
}

func (r *Repository) UpdateAvatar(ctx context.Context, hubId primitive.ObjectID, image string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": hubId}, bson.M{"$set": bson.M{"image": image}})
	return err
}

func (r *Repository) FindByUser(ctx context.Context, userId primitive.ObjectID) ([]Model, error) {
	find, err := r.collection.Find(ctx, bson.M{"users": bson.M{"$in": []primitive.ObjectID{userId}}})
	if err != nil {
		return nil, err
	}
	var hubs []Model
	err = find.All(ctx, &hubs)
	return hubs, err
}

func (r *Repository) AddUsers(ctx context.Context, id primitive.ObjectID, users []primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"users": bson.M{"$each": users}}})
	return err
}

func (r *Repository) RemoveUser(ctx context.Context, id primitive.ObjectID, user primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"users": user}})
	return err
}

func (r Repository) AssertMember(ctx context.Context, hubId primitive.ObjectID, userId primitive.ObjectID) error {
	isMember, err := r.IsMember(ctx, hubId, userId)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member of this hub")
	}
	return nil
}

func (r *Repository) IsMember(ctx context.Context, hubId primitive.ObjectID, userId primitive.ObjectID) (bool, error) {
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
