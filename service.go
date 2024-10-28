package main

import (
	"log"

	"golang.org/x/exp/rand"
)

type Service struct {
	store Storage
}

func NewService(store Storage) (*Service, error) {
	return &Service{
		store: store,
	}, nil
}

func (s *Service) CreateNewAccount(createAccountReq CreateAccountRequest) (*Account, error) {
	accountNumber := 283050400000 + rand.Int63n(99999-10000) + 10000
	err := s.store.SaveAccount(&Account{
		PrimaryUser:   User(createAccountReq),
		AccountNumber: accountNumber,
	})

	if err != nil {
		log.Printf("Error saving account: %v", err)
		return nil, err
	}

	account, err := s.store.RetrieveAccount(accountNumber)
	if err != nil {
		log.Printf("Error retrieving account: %v", err)
		return nil, err
	}
	return account, nil

}

func (s *Service) GetAccount(accountNumber int64) (*Account, error) {
	account, err := s.store.RetrieveAccount(accountNumber)
	if err != nil {
		log.Printf("Error retrieving account: %v", err)
		return nil, err
	}
	return account, nil
}

func (s *Service) GetAccounts() ([]*Account, error) {
	accounts, err := s.store.RetrieveAccounts()
	if err != nil {
		log.Printf("Unable to retrieve all accounts: %v", err)
		return nil, err
	}
	return accounts, err

}

func (s *Service) DeleteAccount(accountNumber int64) error {
	_, err := s.store.RetrieveAccount(accountNumber)
	if err != nil {
		log.Printf("Error deleting account: %v", err)
		return err
	}
	s.store.DeleteAccount(accountNumber)
	return nil
}

func (s *Service) UpdateAccount(updateAccountReq UpdateAccountRequest, accountNumber int64) (*Account, error) {

	account, err := s.store.RetrieveAccount(accountNumber)
	if err != nil {
		log.Printf("Error retrieving account: %v", err)
		return nil, err
	}

	account.PrimaryUser.FirstName = updateAccountReq.FirstName
	account.PrimaryUser.LastName = updateAccountReq.LastName
	account.PrimaryUser.Email = updateAccountReq.Email
	account.PrimaryUser.PhoneNumber = updateAccountReq.PhoneNumber
	account.PrimaryUser.IdentificationNumber = updateAccountReq.IdentificationNumber

	err = s.store.SaveAccount(account)

	if err != nil {
		log.Printf("Error saving account: %v", err)
		return nil, err
	}

	return account, nil

}
