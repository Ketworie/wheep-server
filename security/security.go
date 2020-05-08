package security

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
	"wheep-server/db"
	"wheep-server/user"
)

const sessionLimit int = 5

type Session struct {
	ID   primitive.ObjectID `bson:"_id"`
	U    user.Model         `bson:"u"`
	Last time.Time          `bson:"last"`
}

type Gate struct {
	us *user.Service
	sc *mongo.Collection
}

var g *Gate
var once sync.Once

func initGate() {
	g = &Gate{
		us: user.GetService(),
		sc: db.GetDB().Collection("session"),
	}
}

func GetGate() *Gate {
	once.Do(initGate)
	return g
}

func (g *Gate) Login(login string, password string) (primitive.ObjectID, error) {
	u, err := g.us.GetByLogin(login)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if u.Password != password {
		return primitive.NilObjectID, errors.New("password incorrect")
	}
	session := Session{
		ID:   primitive.NewObjectID(),
		U:    u,
		Last: time.Now(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err = g.sc.InsertOne(ctx, session)
	g.checkSessionLimit(login)
	return session.ID, err
}

func (g *Gate) checkSessionLimit(login string) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	find, _ := g.sc.Find(ctx, bson.M{"u.login": login}, options.Find().SetProjection(bson.M{"_id": 1}), options.Find().SetSort(bson.M{"last": -1}), options.Find().SetSkip(int64(sessionLimit)))
	var res []struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	find.All(ctx, &res)
	if len(res) != 0 {
		ids := make([]primitive.ObjectID, len(res))
		for i, re := range res {
			ids[i] = re.ID
		}
		ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
		defer cancel()
		g.sc.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	}
}

func (g *Gate) Authorize(sid primitive.ObjectID) (user.Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	one := g.sc.FindOne(ctx, bson.M{"_id": sid})
	var s Session
	if one.Err() != nil {
		return s.U, one.Err()
	}
	err := one.Decode(&s)
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
		defer cancel()
		g.sc.UpdateOne(ctx, bson.M{"_id": s.ID}, bson.M{"$set": bson.M{"last": time.Now()}})
	}
	return s.U, err
}
