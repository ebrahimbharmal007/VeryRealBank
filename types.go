package main

import (
	"time"
)

type User struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email"`
	PhoneNumber          int64  `json:"phone_number"`
	IdentificationNumber int64  `json:"identification_number"`
}

type Transaction struct {
	TransactionNumber int64     `json:"transaction_number"`
	Description       string    `json:"description"`
	Amount            int64     `json:"amount"`
	Date              time.Time `json:"date"`
}

type Account struct {
	PrimaryUser   User           `json:"primary_user"`
	AccountNumber int64          `json:"account_number"`
	Balance       int64          `json:"balance"`
	Transactions  []*Transaction `json:"transactions"`
}

// func CreateNewAccount(firstName string, lastName string, email string, phoneNumber int64, identificationNumber int64) *Account {
// 	user := User{
// 		FirstName:            firstName,
// 		LastName:             lastName,
// 		Email:                email,
// 		PhoneNumber:          phoneNumber,
// 		IdentificationNumber: identificationNumber,
// 	}
// 	return &Account{
// 		PrimaryUser:   user,
// 		AccountNumber: 283050400000 + rand.Int63n(99999-10000) + 10000,
// 		Transactions:  []*Transaction{},
// 	}
// }

type CreateAccountRequest struct {
	FirstName            string `json:"first_name" validate:"required"`
	LastName             string `json:"last_name" validate:"required"`
	Email                string `json:"email" validate:"email,required"`
	PhoneNumber          int64  `json:"phone_number" validate:"required,min=1000000000,max=9999999999"`
	IdentificationNumber int64  `json:"identification_number" validate:"required,min=1000000,max=999999999"`
}

type UpdateAccountRequest struct {
	FirstName            string `json:"first_name" validate:"required"`
	LastName             string `json:"last_name" validate:"required"`
	Email                string `json:"email" validate:"email,required"`
	PhoneNumber          int64  `json:"phone_number" validate:"required,min=1000000000,max=9999999999"`
	IdentificationNumber int64  `json:"identification_number" validate:"required,min=1000000,max=999999999"`
}
