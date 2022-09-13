package repository

import "database/sql"

type BalancePostgres struct {
	db *sql.DB
}

func (b *BalancePostgres) Add() {

}
