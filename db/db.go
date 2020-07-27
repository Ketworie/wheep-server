package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var client *mongo.Client
var db *mongo.Database
var cOnce sync.Once
var dbOnce sync.Once
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
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://185.162.10.137:8333"))
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

type M bson.M

func (m M) LastModified() M {
	m["$currentDate"] = M{"lastModified": true}
	return m
}
