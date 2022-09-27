package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host              = "localhost"
	port              = "5436"
	userName          = "postgres"
	dbName            = "postgres"
	password          = "qwerty"
	sslMode           = "disable"
	usersTable        = "users"
	balanceTable      = "balance"
	transactionsTable = "transactions"
)

func NewPostgresDB() (*sqlx.DB, error) {
	dbSource := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, userName, dbName, password, sslMode)
	db, err := sqlx.Open("postgres", dbSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}
