package security

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
	"wheep-server/db"
	"wheep-server/user"
)

const sessionLimit int = 5

type Session struct {
	ID     uuid.UUID          `bson:"_id"`
	UserId primitive.ObjectID `bson:"userId"`
	Last   time.Time          `bson:"last"`
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

func (g *Gate) Login(login string, password string) (string, error) {
	u, err := g.us.GetByLogin(login)
	if err != nil {
		return "", err
	}
	if u.Password != password {
		return "", errors.New("password incorrect")
	}
	session := Session{
		ID:     uuid.New(),
		UserId: u.ID,
		Last:   time.Now(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	_, err = g.sc.InsertOne(ctx, session)
	go func() {
		e := g.checkSessionLimit(u.ID)
		if e != nil {
			log.Fatalf("Cannot check session limit. %v", e)
		}
	}()
	return session.ID.String(), err
}

func (g *Gate) checkSessionLimit(userId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	find, _ := g.sc.Find(ctx, bson.M{"userId": userId}, options.Find().SetProjection(bson.M{"_id": 1}), options.Find().SetSort(bson.M{"last": -1}), options.Find().SetSkip(int64(sessionLimit)))
	var res []struct {
		ID uuid.UUID `bson:"_id"`
	}
	ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	err := find.All(ctx, &res)
	if err != nil {
		return err
	}
	if len(res) != 0 {
		ids := make([]uuid.UUID, len(res))
		for i, re := range res {
			ids[i] = re.ID
		}
		ctx, cancel = context.WithTimeout(context.Background(), db.DBTimeout)
		defer cancel()
		_, err = g.sc.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Gate) Authorize(sid uuid.UUID) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
	defer cancel()
	one := g.sc.FindOne(ctx, bson.M{"_id": sid})
	var s Session
	if one.Err() != nil {
		return primitive.NilObjectID, one.Err()
	}
	err := one.Decode(&s)
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), db.DBTimeout)
		defer cancel()
		_, err := g.sc.UpdateOne(ctx, bson.M{"_id": s.ID}, bson.M{"$set": bson.M{"last": time.Now()}})
		if err != nil {
			return primitive.NilObjectID, err
		}
	}
	return s.UserId, err
}
