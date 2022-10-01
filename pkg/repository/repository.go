package repository

import (
	"github.com/cucumberjaye/balanceAPI"
	"github.com/jmoiron/sqlx"
)

type Balance interface {
	Add(userData balanceAPI.UserData) error
	Decrease(userData balanceAPI.UserData) error
	Transfer(usersData balanceAPI.TwoUsers) error
	GetBalance(userId int) (int, error)
	GetTransactions(userId int, sortBy string) ([]balanceAPI.Transactions, error)
}

type Repository struct {
	Balance
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Balance: NewBalancePostgres(db)}
}
