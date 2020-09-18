package mq

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
	"wheep-server/config"
)

var once sync.Once
var conn *amqp.Connection

func GetConnection() *amqp.Connection {
	once.Do(initClient)
	return conn
}

func initClient() {
	var err error
	conn, err = amqp.Dial(config.Get().MQAddress)
	if err != nil {
		log.Fatal(err)
	}
}
