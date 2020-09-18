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
	ur *user.Repository
	sc *mongo.Collection
}

var g *Gate
var once sync.Once

func initGate() {
	g = &Gate{
		ur: user.GetRepository(),
		sc: db.GetDB().Collection("session"),
	}
}

func GetGate() *Gate {
	once.Do(initGate)
	return g
}

func (g *Gate) Login(ctx context.Context, login string, password string) (string, error) {
	u, err := g.ur.GetByLogin(ctx, login)
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
	_, err = g.sc.InsertOne(ctx, session)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()
		e := g.checkSessionLimit(ctx, u.ID)
		if e != nil {
			log.Fatalf("Cannot check session limit. %v", e)
		}
	}()
	return session.ID.String(), err
}

func (g *Gate) checkSessionLimit(ctx context.Context, userId primitive.ObjectID) error {
	find, _ := g.sc.Find(ctx, bson.M{"userId": userId}, options.Find().SetProjection(bson.M{"_id": 1}), options.Find().SetSort(bson.M{"last": -1}), options.Find().SetSkip(int64(sessionLimit)))
	var res []struct {
		ID uuid.UUID `bson:"_id"`
	}
	err := find.All(ctx, &res)
	if err != nil {
		return err
	}
	if len(res) != 0 {
		ids := make([]uuid.UUID, len(res))
		for i, re := range res {
			ids[i] = re.ID
		}
		_, err = g.sc.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Gate) Authorize(ctx context.Context, sid uuid.UUID) (primitive.ObjectID, error) {
	one := g.sc.FindOne(ctx, bson.M{"_id": sid})
	var s Session
	if one.Err() != nil {
		return primitive.NilObjectID, one.Err()
	}
	err := one.Decode(&s)
	if err == nil {
		_, err := g.sc.UpdateOne(ctx, bson.M{"_id": s.ID}, bson.M{"$set": bson.M{"last": time.Now()}})
		if err != nil {
			return primitive.NilObjectID, err
		}
	}
	return s.UserId, err
}

func (g *Gate) GetOffline(ctx context.Context) ([]primitive.ObjectID, error) {
	find, err := g.sc.Find(ctx, bson.M{"$lt": bson.M{"last": time.Now().Add(-time.Hour * 24 * 30)}})
	if err != nil {
		return nil, err
	}
	var s []Session
	err = find.All(ctx, &s)
	if err != nil {
		return nil, err
	}
	var ids []primitive.ObjectID
	for _, session := range s {
		ids = append(ids, session.UserId)
	}
	return ids, nil
}
