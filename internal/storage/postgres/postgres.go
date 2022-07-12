package postgres

import (
	"fmt"
	nats "github.com/Vallghall/wb-test-l0/internal/model/message"
	"github.com/jmoiron/sqlx"
)

const (
	tableName = "messages"
	idField   = "id"
	msg       = "message"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db}
}

func (s *Storage) Store(id string, m *nats.Message) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (%s, %s) 
			VALUES ($1, $2::jsonb);
	`, tableName, idField, msg)

	tx, err := s.db.Begin()
	if _, err = tx.Exec(query, id, m); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Storage) LoadMessageById(id string) (*nats.Message, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
			WHERE %s = $1;
	`, msg, tableName, idField)

	row := s.db.QueryRow(query, id)
	var message nats.Message
	if err := row.Scan(&message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *Storage) LoadAllMessages() ([]*nats.Message, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s;`, msg, tableName)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	res := make([]*nats.Message, 0)
	for rows.Next() {
		var row nats.Message
		err = rows.Scan(&row)
		if err != nil {
			return nil, err
		}

		res = append(res, &row)
	}

	return res, nil
}
