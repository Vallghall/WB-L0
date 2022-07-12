package message

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Message struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Items   `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shard_key"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func (m Message) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Message) Scan(obj any) error {
	b, ok := obj.([]byte)
	if !ok {
		return errors.New("type assertion to type []byte failed")
	}

	return json.Unmarshal(b, &m)
}
