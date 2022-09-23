package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/cucumberjaye/balanceAPI"
	"log"
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
	row := db.QueryRow(query, userId)
	if err := row.Scan(&sum); err != nil {
		log.Printf("first %s", err.Error())
		return false, nil
	}
	return true, nil
}

func createUser(db *sql.DB, userData balanceAPI.UserData) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}
	createUserQuery := fmt.Sprintf("INSERT INTO %s (id, name, surname) VALUES ($1, $2, $3)", usersTable)
	_, err = tx.Exec(createUserQuery, userData.Id, userData.Name, userData.Surname)
	if err != nil {
		tx.Rollback()
		log.Fatalf(err.Error())
		return err
	}
	createBalanceQuery := fmt.Sprintf("INSERT INTO %s (user_id, balance) VALUES ($1, $2)", balanceTable)
	_, err = tx.Exec(createBalanceQuery, userData.Id, userData.Sum)
	if err != nil {
		tx.Rollback()
		log.Fatalf(err.Error())
		return err
	}
	return tx.Commit()
}

func divisionUsersData(user balanceAPI.User, sum int) balanceAPI.UserData {
	return balanceAPI.UserData{Id: user.Id, Name: user.Name, Surname: user.Surname, Sum: sum}
}

func (b *BalancePostgres) Add(userData balanceAPI.UserData) error {
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
		log.Fatalf(err.Error())
		return err
	}

	return nil
}

func (b *BalancePostgres) Decrease(userData balanceAPI.UserData) error {
	var balance int
	checkQuery := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=$1", balanceTable)
	row := b.db.QueryRow(checkQuery, userData.Id)
	if err := row.Scan(&balance); err != nil {
		log.Fatalf(err.Error())
		return err
	}
	if balance < userData.Sum {
		err := errors.New("balance < 0 ")
		log.Printf(err.Error())
		return err
	}
	decreaseQuery := fmt.Sprintf("UPDATE %s SET balance=(SELECT balance FROM %s WHERE user_id=$1)-$2 WHERE user_id=$1", balanceTable, balanceTable)
	_, err := b.db.Exec(decreaseQuery, userData.Id, userData.Sum)
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}

	return nil
}

func (b *BalancePostgres) Transfer(usersData balanceAPI.TwoUsers) error {
	tx, err := b.db.Begin()
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}
	decreaseUserData := divisionUsersData(usersData.DecreaseMoneyUser, usersData.Sum)
	err = b.Decrease(decreaseUserData)
	if err != nil {
		tx.Rollback()
		log.Fatalf(err.Error())
		return err
	}
	addUserData := divisionUsersData(usersData.AddMoneyUser, usersData.Sum)
	err = b.Add(addUserData)
	if err != nil {
		tx.Rollback()
		log.Fatalf(err.Error())
		return err
	}

	return tx.Commit()
}

func (b *BalancePostgres) GetBalance(userId int) (int, error) {
	query := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=$1", balanceTable)
	row := b.db.QueryRow(query, userId)
	var balance int
	if err := row.Scan(&balance); err != nil {
		log.Printf(err.Error())
		return 0, err
	}

	return balance, nil
}
