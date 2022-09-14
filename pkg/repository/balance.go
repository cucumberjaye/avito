package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/cucumberjaye/avito"
)

type BalancePostgres struct {
	db *sql.DB
}

func NewBalancePostgres(db *sql.DB) *BalancePostgres {
	return &BalancePostgres{db: db}
}

func checkIdInDB(db *sql.DB, userId int) (bool, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", usersTable)
	var sum int
	row, err := db.Query(query, userId)
	if err != nil {
		return false, err
	}
	if err := row.Scan(&sum); err != nil {
		return false, nil
	}
	return true, nil
}

func createUser(db *sql.DB, userData avito.UserData) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	createUserQuery := fmt.Sprintf("INSERT INTO %s (id, name, surname) VALUES ($1, $2, $3)", usersTable)
	_, err = tx.Exec(createUserQuery, userData.Id, userData.Name, userData.Surname)
	if err != nil {
		tx.Rollback()
		return err
	}
	createBalanceQuery := fmt.Sprintf("INSERT INTO %s (user_id, sum) VALUES ($1, $2)", balanceTable)
	_, err = tx.Exec(createBalanceQuery, userData.Id, userData.Sum)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (b *BalancePostgres) Add(userData avito.UserData) error {
	exists, err := checkIdInDB(b.db, userData.Id)
	if err != nil {
		return err
	}
	if !exists {
		return createUser(b.db, userData)
	}
	query := fmt.Sprintf("UPDATE %s SET balance=(SELECT balance FROM %s WHERE user_id=$1)+$2 WHERE user_id=$1", balanceTable, balanceTable)
	_, err = b.db.Exec(query, userData.Id, userData.Sum)
	if err != nil {
		return err
	}

	return nil
}

func (b *BalancePostgres) Decrease(userData avito.UserData) error {
	var balance int
	checkQuery := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=$1", balanceTable)
	row, err := b.db.Query(checkQuery, userData.Id)
	if err != nil {
		return err
	}
	if err := row.Scan(&balance); err != nil {
		return err
	}
	if balance < userData.Sum {
		return errors.New("not ")
	}
	decreaseQuery := fmt.Sprintf("UPDATE %s SET balance=(SELECT balance FROM %s WHERE user_id=$1)-$2 WHERE user_id=$1", balanceTable, balanceTable)
	_, err = b.db.Exec(decreaseQuery, userData.Id, userData.Sum)
	if err != nil {
		return err
	}

	return nil
}

func (b *BalancePostgres) Transfer(addUserData avito.UserData, decreaseUserData avito.UserData) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	err = b.Decrease(decreaseUserData)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = b.Add(addUserData)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (b *BalancePostgres) GetBalance(userId int) (int, error) {
	query := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=$1", balanceTable)
	row, err := b.db.Query(query, userId)
	if err != nil {
		return 0, err
	}
	var balance int
	if err = row.Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}
