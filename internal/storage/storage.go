package storage

import (
	nats "github.com/Vallghall/wb-test-l0/internal/model/message"
	"github.com/Vallghall/wb-test-l0/internal/storage/inmemory"
	"github.com/Vallghall/wb-test-l0/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type InMemoryStorage interface {
	Store(id string, message *nats.Message)
	LoadMessageById(id string) (*nats.Message, error)
	Synchronize(data []*nats.Message)
}

type PersistentStorage interface {
	Store(id string, message *nats.Message) error
	LoadMessageById(id string) (*nats.Message, error)
	LoadAllMessages() ([]*nats.Message, error)
}

type Storage struct {
	InMemoryStorage
	PersistentStorage
}

func New(db *sqlx.DB) *Storage {
	return &Storage{
		InMemoryStorage:   inmemory.New(),
		PersistentStorage: postgres.New(db),
	}
}
