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

type createExpenseRequest struct {
	AccountID int32
	Name string
	Amount int32
}

func createExpense(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("createExpense")

		userId := gcontext.Get(r, "userId").(int32)
		var request createExpenseRequest
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

		expense, err := worker.CreateExpense(ctx, &communication.CreateExpenseRequest{Amount: request.Amount,AccountID: request.AccountID, Name: request.Name, UserID: userId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(expense)
	}
}

type fetchExpenseRequest struct {
	ExpenseID int32
}

func fetchExpense(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("fetchExpense")

		userId := gcontext.Get(r, "userId").(int32)
		var request fetchExpenseRequest
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

		expense, err := worker.FetchExpense(ctx, &communication.FetchExpenseRequest{UserID: userId, ExpenseID: request.ExpenseID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(expense)
	}
}

type deleteExpenseRequest struct {
	ExpenseID int32
}

func deleteExpense(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("deleteExpense")

		userId := gcontext.Get(r, "userId").(int32)
		var request deleteExpenseRequest
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

		_, err = worker.DeleteExpense(ctx, &communication.DeleteExpenseRequest{UserID: userId, ExpenseID: request.ExpenseID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func registerExpenseRoutes(r *mux.Router, server *server) {

	r.HandleFunc("/expense", auth(createExpense(server))).Methods(http.MethodPost)
	r.HandleFunc("/expense", auth(fetchExpense(server))).Methods(http.MethodGet)
	r.HandleFunc("/expense", auth(deleteExpense(server))).Methods(http.MethodDelete)
}
