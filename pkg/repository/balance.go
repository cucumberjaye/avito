package repository

import (
	"database/sql"
	"github.com/cucumberjaye/avito"
)

type BalancePostgres struct {
	db *sql.DB
}

func NewBalancePostgres(db *sql.DB) *BalancePostgres {
	return &BalancePostgres{db: db}
}

func checkIdInDB(db *sql.DB, userId int) bool {
	return false
}

func (b *BalancePostgres) Add(userData avito.UserData) error {
	return nil
}

func (b *BalancePostgres) Decrease(userData avito.UserData) error {
	return nil
}

func (b *BalancePostgres) Transfer(firstUserData avito.UserData, secondUserData avito.UserData) error {
	return nil
}

func (b *BalancePostgres) GetBalance(userId int) int {
	return 0
}
