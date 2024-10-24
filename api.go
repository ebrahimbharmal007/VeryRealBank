package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

var validate = validator.New(validator.WithRequiredStructEnabled())

type ApiError struct {
	Error string `json:"error"`
}

type APIServer struct {
	listenAddr string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8;")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	log.Println("JSON API server running on port: ", s.listenAddr)

	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount)).Methods("POST")
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		log.Printf("\nunable to decode: %+v\n", err)

		return fmt.Errorf("invalid input data")
	}
	err := validate.Struct(createAccountReq)
	if err != nil {
		log.Printf("\nValidation error: %+v\n", err)

		return fmt.Errorf("invalid input data")
	}
	account := CreateNewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Email, createAccountReq.PhoneNumber, createAccountReq.IdentificationNumber)

	return WriteJSON(w, http.StatusOK, account)
}
