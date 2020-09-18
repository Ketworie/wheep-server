package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
	"wheep-server/config"
)

var client *mongo.Client
var db *mongo.Database
var cOnce sync.Once
var dbOnce sync.Once

// TODO: SWITCH TO REQUEST CONTEXT
var DBTimeout = time.Minute

func GetClient() *mongo.Client {
	cOnce.Do(initClient)
	return client
}

func GetDB() *mongo.Database {
	dbOnce.Do(initDB)
	return db
}

func initDB() {
	db = GetClient().Database("wheep")
}

func initClient() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(config.Get().MongoAddress))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
}
