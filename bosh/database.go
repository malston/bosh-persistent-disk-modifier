package bosh

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase(host string, user string) (*sqlx.DB, error) {
	conn := fmt.Sprintf("postgres://%s@%s/bosh?sslmode=disable", user, host)
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return db, err
	}

	return db, err
}
