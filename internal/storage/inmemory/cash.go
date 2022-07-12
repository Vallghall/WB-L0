package inmemory

import (
	"errors"
	nats "github.com/Vallghall/wb-test-l0/internal/model/message"
)

func New() *Cash {
	return &Cash{
		storage: make(map[string]*nats.Message),
	}
}

type Cash struct {
	storage map[string]*nats.Message
}

func (c *Cash) Store(id string, m *nats.Message) {
	c.storage[id] = m
}

func (c *Cash) LoadMessageById(id string) (*nats.Message, error) {
	msg, ok := c.storage[id]
	if !ok {
		return nil, errors.New("no such value")
	}
	return msg, nil
}

func (c *Cash) Synchronize(data []*nats.Message) {
	for _, datum := range data {
		c.Store(datum.OrderUID, datum)
	}
}
