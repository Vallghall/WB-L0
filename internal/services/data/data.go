package data

import (
	"encoding/json"
	nats "github.com/Vallghall/wb-test-l0/internal/model/message"
	"github.com/Vallghall/wb-test-l0/internal/storage"
	"log"
)

type Service struct {
	storage *storage.Storage
}

func New(s *storage.Storage) *Service {
	return &Service{s}
}

func (s *Service) CashMessage(data []byte) {
	var msg nats.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("message parsing failed: %v", err)
		return
	}

	s.storage.InMemoryStorage.Store(msg.OrderUID, &msg)
}

func (s *Service) PersistMessage(data []byte) {
	var msg nats.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("message parsing failed: %v", err)
		return
	}

	s.storage.PersistentStorage.Store(msg.OrderUID, &msg)
}

func (s *Service) LoadCachedMsgById(id string) (*nats.Message, error) {
	return s.storage.InMemoryStorage.LoadMessageById(id)
}

func (s *Service) LoadPersistedMsgById(id string) (*nats.Message, error) {
	return s.storage.PersistentStorage.LoadMessageById(id)
}

func (s *Service) SynchronizeCash() {
	data, err := s.storage.PersistentStorage.LoadAllMessages()
	if err != nil {
		log.Fatalf("synchronization failed: %v\n", err)
	}
	s.storage.InMemoryStorage.Synchronize(data)
}
