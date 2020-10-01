package mq

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
	"wheep-server/config"
)

var once sync.Once
var conn *amqp.Connection
var ch *amqp.Channel

func GetConnection() *amqp.Connection {
	once.Do(initClient)
	return conn
}

func ChannelSupplier() func() *amqp.Channel {
	once.Do(initClient)
	return func() *amqp.Channel {
		return ch
	}
}

func initClient() {
	connect()
}

func connect() {
	var err error
	conn, err = amqp.Dial(config.Get().MQAddress)
	if err != nil {
		onError(err)
	}
	ch, err = conn.Channel()
	if err != nil {
		onError(err)
	}
	go sustainConnection()
}

func onError(err error) {
	log.Print(err)
	time.Sleep(time.Minute)
	connect()
}

func sustainConnection() {
	errorCh := make(chan *amqp.Error)
	errorCh = conn.NotifyClose(errorCh)
	err := <-errorCh
	log.Print(err)
	connect()
}
