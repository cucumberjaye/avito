package repository

import (
	"database/sql"
	"github.com/cucumberjaye/avito"
)

type Balance interface {
	Add(userData avito.UserData) error
	Decrease(userData avito.UserData) error
	Transfer(firstUserData avito.UserData, secondUserData avito.UserData) error
	GetBalance(userId int) int
}

type Repository struct {
	Balance
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{Balance: NewBalancePostgres(db)}
}
