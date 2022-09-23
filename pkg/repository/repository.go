package repository

import (
	"database/sql"
	"github.com/cucumberjaye/balanceAPI"
)

type Balance interface {
	Add(userData balanceAPI.UserData) error
	Decrease(userData balanceAPI.UserData) error
	Transfer(usersData balanceAPI.TwoUsers) error
	GetBalance(userId int) (int, error)
}

type Repository struct {
	Balance
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{Balance: NewBalancePostgres(db)}
}
