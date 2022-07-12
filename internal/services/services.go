package services

import (
	nats "github.com/Vallghall/wb-test-l0/internal/model/message"
	"github.com/Vallghall/wb-test-l0/internal/services/data"
	"github.com/Vallghall/wb-test-l0/internal/storage"
)

type DataService interface {
	CashMessage(data []byte)
	PersistMessage(data []byte)
	LoadCachedMsgById(id string) (*nats.Message, error)
	LoadPersistedMsgById(id string) (*nats.Message, error)
	SynchronizeCash()
}

type Service struct {
	DataService
}

func New(s *storage.Storage) *Service {
	return &Service{
		DataService: data.New(s),
	}
}
