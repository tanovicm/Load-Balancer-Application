package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"encoding/json"

	"github.com/gorilla/mux"
	gcontext "github.com/gorilla/context"
	"github.com/tanovicm/tenderly/communication"
)

type createBankAccountRequest struct {
	Name string
}

func createBankAccount(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("CreateBankAccount")

		userId := gcontext.Get(r, "userId").(int32)
		log.Printf("userId %v", userId)
		var request createBankAccountRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		worker, err := s.GetWorker()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		account, err := worker.CreateBankAccount(ctx, &communication.CreateBankAccountRequest{UserID: userId, Name: request.Name})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(account)
	}
}

type fetchBankAccountRequest struct {
	AccountID int32
}

func fetchBankAccount(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("fetchBankAccount")

		userId := gcontext.Get(r, "userId").(int32)
		var request fetchBankAccountRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		worker, err := s.GetWorker()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		account, err := worker.FetchBankAccount(ctx, &communication.FetchBankAccountRequest{UserID: userId, AccountID: request.AccountID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(account)
	}
}

type deleteBankAccountRequest struct {
	AccountID int32
}

func deleteBankAccount(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("deleteBankAccount")

		userId := gcontext.Get(r, "userId").(int32)
		var request deleteBankAccountRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		worker, err := s.GetWorker()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = worker.DeleteBankAccount(ctx, &communication.DeleteBankAccountRequest{UserID: userId, AccountID: request.AccountID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func registerBankRoutes(r *mux.Router, server *server) {

	r.HandleFunc("/bank_account", auth(createBankAccount(server))).Methods(http.MethodPost)
	r.HandleFunc("/bank_account", auth(fetchBankAccount(server))).Methods(http.MethodGet)
	r.HandleFunc("/bank_account", auth(deleteBankAccount(server))).Methods(http.MethodDelete)
}
