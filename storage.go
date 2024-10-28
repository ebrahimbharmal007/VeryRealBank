package main

import "fmt"

type Storage interface {
	SaveAccount(*Account) error
	DeleteAccount(accountNumber int64) error
	RetrieveAccounts() ([]*Account, error)
	RetrieveAccount(accountNumber int64) (*Account, error)
}

type InMemoryRepository struct {
	accounts map[int64]*Account
}

func NewInMemoryRepository() (*InMemoryRepository, error) {
	return &InMemoryRepository{
		accounts: make(map[int64]*Account),
	}, nil
}

func (r *InMemoryRepository) Init() error {
	return nil
}

func (r *InMemoryRepository) SaveAccount(account *Account) error {
	r.accounts[account.AccountNumber] = account
	return nil
}

func (r *InMemoryRepository) RetrieveAccounts() ([]*Account, error) {
	accounts := []*Account{}
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *InMemoryRepository) RetrieveAccount(accountNumber int64) (*Account, error) {
	if account, ok := r.accounts[accountNumber]; ok {
		return account, nil
	}

	return nil, fmt.Errorf("key does not exist")
}

func (r *InMemoryRepository) DeleteAccount(accountNumber int64) error {
	delete(r.accounts, accountNumber)
	return nil
}
