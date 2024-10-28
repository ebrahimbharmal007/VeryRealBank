package main

import (
	"encoding/json"
	"fmt"
	"strconv"

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

type Message struct {
	Message string `json:"message"`
}

type APIServer struct {
	port    string
	service Service
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

func NewAPIServer(port string, service Service) *APIServer {
	return &APIServer{
		port:    port,
		service: service,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleGetAllAccounts)).Methods("GET")
	router.HandleFunc("/account/{accnum}", makeHTTPHandleFunc(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account/{accnum}", makeHTTPHandleFunc(s.handleDeleteAccount)).Methods("DELETE")

	log.Printf("Bank JSON API server is running on port %s", s.port)
	if err := http.ListenAndServe(":"+s.port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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

	account, err := s.service.CreateNewAccount(*createAccountReq)

	if err != nil {
		log.Printf("Error creating account account: %v", err)
		return WriteJSON(w, http.StatusInternalServerError, err)
	}
	// account := CreateNewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Email, createAccountReq.PhoneNumber, createAccountReq.IdentificationNumber)

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.service.GetAccounts()

	if err != nil {
		log.Printf("Unable to retrieve all accounts: %v", err)
		return WriteJSON(w, http.StatusInternalServerError, err)
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["accnum"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("account %s not found", idStr)
	}
	account, err := s.service.GetAccount(int64(id))

	if err != nil {
		log.Printf("Unable to retrieve account: %v", err)
		return WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["accnum"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("account %s not found", idStr)
	}
	err = s.service.DeleteAccount(int64(id))

	if err != nil {
		log.Printf("Unable to delete account: %v", err)
		return WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	return WriteJSON(w, http.StatusOK, Message{Message: "Account Deleted Successfully"})
}
