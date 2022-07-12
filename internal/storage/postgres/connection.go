package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"

	_ "github.com/lib/pq"
)

type Configs struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
	SSLMode  string
}

func NewConnection(c *Configs) *sqlx.DB {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=%s",
			c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode))

	if err != nil {
		log.Fatalln(fmt.Sprintf("connection failed: %v", err))
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "messages" (
	    id varchar(255) PRIMARY KEY,
	    message jsonb
	)`)

	if err != nil {
		log.Fatalln(fmt.Sprintf("initialization failed: %v", err))
	}

	return db
}
