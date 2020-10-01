package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sync"
	"wheep-server/hub"
	"wheep-server/message"
	"wheep-server/mq"
	"wheep-server/security"
)

var s *Service
var sOnce sync.Once

func GetService() *Service {
	sOnce.Do(initService)
	return s
}

func initService() {
	supplier := mq.ChannelSupplier()
	s = &Service{
		chanSupplier: supplier,
		hubSync: &hubSync{
			RWMutex:  sync.RWMutex{},
			hubUsers: make(map[primitive.ObjectID]*idSync),
		},
		exchangeSync: &exchangeSync{
			RWMutex:  sync.RWMutex{},
			channels: make(map[string]bool),
		},
		repo: struct {
			*hub.Repository
			*security.Gate
		}{
			Repository: hub.GetRepository(),
			Gate:       security.GetGate(),
		},
	}
}

type Service struct {
	chanSupplier func() *amqp.Channel
	hubSync      *hubSync
	exchangeSync *exchangeSync
	repo         userRepository
}

type userRepository interface {
	GetUserIds(ctx context.Context, hubId primitive.ObjectID) ([]primitive.ObjectID, error)
	GetOffline(ctx context.Context) ([]primitive.ObjectID, error)
}

type exchangeSync struct {
	sync.RWMutex
	channels map[string]bool
}

type idSync struct {
	sync.RWMutex
	ids []primitive.ObjectID
}

type hubSync struct {
	sync.RWMutex
	hubUsers map[primitive.ObjectID]*idSync
}

func (s *Service) Fanout(ctx context.Context, m message.View) {
	hs := s.hubSync
	hubId := m.HubId
	hs.RLock()
	uSync, ok := hs.hubUsers[hubId]
	hs.RUnlock()
	if !ok {
		userIds, err := s.repo.GetUserIds(ctx, hubId)
		if err != nil {
			log.Printf("Error during message fanout. Can't get hub's users': %v", err)
			return
		}
		uSync = &idSync{ids: userIds}
		hs.Lock()
		hs.hubUsers[hubId] = uSync
		hs.Unlock()
	}
	body, err := json.Marshal(m)
	if err != nil {
		log.Printf("Error during message fanout. Can't marshall message': %v", err)
		return
	}
	for _, userId := range uSync.ids {
		s.publishJSON(userId.Hex(), body)
	}
}

func (s *Service) publishJSON(exchange string, body []byte) {
	e := s.exchangeSync
	e.RLock()
	defer e.RUnlock()
	if _, ok := e.channels[exchange]; !ok {
		return
	}
	err := s.chanSupplier().Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Can't publish message': %v", err)
	}
}

func (s *Service) SetupExchange(userId primitive.ObjectID, token string) error {
	e := s.exchangeSync
	exchangeName := userId.Hex()
	err := s.chanSupplier().ExchangeDeclare(
		exchangeName,
		"fanout",
		false,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return err
	}
	e.Lock()
	e.channels[exchangeName] = true
	e.Unlock()
	qName := fmt.Sprintf("q-%v", token)
	queue, err := s.chanSupplier().QueueDeclare(
		qName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = s.chanSupplier().QueueBind(
		queue.Name,
		"",
		exchangeName,
		false,
		nil,
	)
	return err
}
