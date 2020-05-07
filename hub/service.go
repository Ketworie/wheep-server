package hub

import (
	"github.com/google/uuid"
	"sync"
	"wheep-server/db"
)

type Service struct {
	*Repository
}

var s *Service
var once sync.Once

func initService() {
	s = &Service{&Repository{
		collection: db.GetDB().Collection("hub"),
	}}
}

func GetService() *Service {
	once.Do(initService)
	return s
}

func (s *Service) Rename(id uuid.UUID, name string) error {
	model, err := s.Get(id)
	if err != nil {
		return err
	}
	model.Name = name
	return s.Update(model)
}

func (s *Service) AddUsers(id uuid.UUID, users []uuid.UUID) error {
	model, err := s.Get(id)
	if err != nil {
		return err
	}
	for _, u := range users {
		model.Users[u] = true
	}
	return s.Update(model)
}

func (s *Service) RemoveUsers(id uuid.UUID, users []uuid.UUID) error {
	model, err := s.Get(id)
	if err != nil {
		return err
	}
	for _, u := range users {
		delete(model.Users, u)
	}
	return s.Update(model)
}

func (s *Service) FindHubs(userId uuid.UUID) ([]View, error) {
	hubs, err := s.FindByUser(userId)
	if err != nil {
		return nil, err
	}
	views := make([]View, len(hubs))
	for i, v := range hubs {
		views[i] = View{
			ID:        v.ID,
			Name:      v.Name,
			UserCount: len(v.Users),
		}
	}
	return views, nil
}
