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

func (b *BalancePostgres) Add(userData balanceAPI.UserData) error {
	tx, err := b.db.Begin()
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	tx, err = addWithTx(userData, tx, b.db)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			log.Fatalf(rErr.Error())
			return rErr
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}

func (b *BalancePostgres) Decrease(userData balanceAPI.UserData) error {
	tx, err := b.db.Begin()
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	tx, err = decreaseWithTx(userData, tx)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			log.Fatalf(rErr.Error())
			return rErr
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	return nil
}

func (b *BalancePostgres) Transfer(usersData balanceAPI.TwoUsers) error {
	tx, err := b.db.Begin()
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	decreaseUserData := divisionUsersData(usersData.DecreaseMoneyUser, usersData.Sum)
	tx, err = decreaseWithTx(decreaseUserData, tx)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			log.Fatalf(rErr.Error())
			return rErr
		}
		log.Printf(err.Error())
		return err
	}
	addUserData := divisionUsersData(usersData.AddMoneyUser, usersData.Sum)
	tx, err = addWithTx(addUserData, tx, b.db)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			log.Fatalf(rErr.Error())
			return rErr
		}
		log.Printf(err.Error())
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	return nil
}

func (b *BalancePostgres) GetBalance(userId int) (int, error) {
	query := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=$1", balanceTable)
	row := b.db.QueryRow(query, userId)
	var balance int
	if err := row.Scan(&balance); err != nil {
		if err.Error() == "sql: no rows in result set" {
			errString := fmt.Sprintf("id = %d user is not in the database", userId)
			err := errors.New(errString)
			log.Printf(err.Error())
			return 0, err
		}
		log.Printf(err.Error())
		return 0, err
	}

	return balance, nil
}

func checkIdInDB(db *sql.DB, userId int) (bool, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", usersTable)
	var sum int
	row := db.QueryRow(query, userId)
	if err := row.Scan(&sum); err != nil {
		log.Printf(err.Error())
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
		rErr := tx.Rollback()
		if rErr != nil {
			log.Fatalf(rErr.Error())
			return rErr
		}
		log.Fatalf(err.Error())
		return err
	}
	createBalanceQuery := fmt.Sprintf("INSERT INTO %s (user_id, balance) VALUES ($1, $2)", balanceTable)
	_, err = tx.Exec(createBalanceQuery, userData.Id, userData.Sum)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			log.Fatalf(rErr.Error())
			return rErr
		}
		log.Fatalf(err.Error())
		return err
	}
	return tx.Commit()
}

func divisionUsersData(user balanceAPI.User, sum int) balanceAPI.UserData {
	return balanceAPI.UserData{User: user, Sum: sum}
}

func decreaseWithTx(userData balanceAPI.UserData, tx *sql.Tx) (*sql.Tx, error) {
	var balance int
	checkQuery := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=$1", balanceTable)
	row := tx.QueryRow(checkQuery, userData.Id)
	if err := row.Scan(&balance); err != nil {
		if err.Error() == "sql: no rows in result set" {
			errString := fmt.Sprintf("id = %d user is not in the database", userData.Id)
			err := errors.New(errString)
			log.Printf(err.Error())
			return tx, err
		}
		log.Printf(err.Error())
		return tx, err
	}
	if balance < userData.Sum {
		err := errors.New("balance < 0 ")
		log.Printf(err.Error())
		return tx, err
	}
	decreaseQuery := fmt.Sprintf("UPDATE %s SET balance=(SELECT balance FROM %s WHERE user_id=$1)-$2 WHERE user_id=$1", balanceTable, balanceTable)
	_, err := tx.Exec(decreaseQuery, userData.Id, userData.Sum)
	if err != nil {
		log.Fatalf(err.Error())
		return tx, err
	}

	return tx, nil
}

func addWithTx(userData balanceAPI.UserData, tx *sql.Tx, db *sql.DB) (*sql.Tx, error) {
	exists, err := checkIdInDB(db, userData.Id)
	if err != nil {
		return tx, err
	}
	if !exists {
		err = createUser(db, userData)
		if err != nil {
			return tx, err
		}
		return tx, nil
	}
	query := fmt.Sprintf("UPDATE %s SET balance=(SELECT balance FROM %s WHERE user_id=$1)+$2 WHERE user_id=$1", balanceTable, balanceTable)
	_, err = tx.Exec(query, userData.Id, userData.Sum)
	if err != nil {
		log.Fatalf(err.Error())
		return tx, err
	}

	return tx, nil
}
