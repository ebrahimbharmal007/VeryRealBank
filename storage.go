package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `
	create table if not exists account(
	id serial primary key,
	first_name varchar(50),
	last_name varchar(50),
	number serial,
	balance serial,
	created_at timestamp
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) (*Account, error) {
	createAccountQuery := `
	insert into account (first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5)
	returning id, first_name, last_name, number, balance, created_at`

	// Create a new Account to hold the returned data
	newAccount := &Account{}

	// QueryRow is used for a single row result
	err := s.db.QueryRow(createAccountQuery, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt).Scan(
		&newAccount.ID,
		&newAccount.FirstName,
		&newAccount.LastName,
		&newAccount.Number,
		&newAccount.Balance,
		&newAccount.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Return the newly created account
	return newAccount, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {

	rows, err := s.db.Query(`SELECT * from account`)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query(`delete from account where id = $1`, id)
	return err
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {

	rows, err := s.db.Query(`SELECT * from account where id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account  %d not found", id)

	// return account, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)

	return account, err
}
