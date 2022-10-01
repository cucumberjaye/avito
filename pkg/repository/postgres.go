package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

const (
	usersTable        = "users"
	balanceTable      = "balance"
	transactionsTable = "transactions"
)

var (
	host              = os.Getenv("HOST")
	port              = os.Getenv("PORT")
	userName          = os.Getenv("USER")
	dbName            = os.Getenv("DBNAME")
	password          = os.Getenv("PASSWORD")
	sslMode           = os.Getenv("SSLMODE")
)

func NewPostgresDB() (*sqlx.DB, error) {
	dbSource := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, userName, dbName, password, sslMode)
	db, err := sqlx.Open("postgres", dbSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}
