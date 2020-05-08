package hub

import (
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
