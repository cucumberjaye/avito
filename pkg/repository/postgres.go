package repository

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = "5432"
	userName = "iskander"
	dbName   = "postgres"
	password = ""
	sslMode  = ""
)

func NewPostgresDB() (*sql.DB, error) {
	dbSource := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, userName, dbName, password, sslMode)
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}
