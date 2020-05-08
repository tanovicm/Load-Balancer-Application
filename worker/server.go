package main

import (
	"errors"
	"context"
	"log"

	"github.com/go-redis/redis/v7"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/tanovicm/tenderly/communication"
)

type server struct {
	Client *redis.Client
	communication.UnimplementedWorkerServer
}

func (s *server) CreateBankAccount(ctx context.Context, in *communication.CreateBankAccountRequest) (*communication.Account, error) {

	log.Printf("CreateBankAccount %v", *in)

	account := Account{
		AccountID: getAutoIncrementID(s.Client, "accountID"),
		UserID: in.GetUserID(),
		Name: in.GetName(),
	}
	err := setTable(s.Client, &account)
	if err != nil {
		return nil, err
	}

	return &communication.Account{
		UserID: account.UserID,
		AccountID: account.AccountID,
		Name: account.Name,
		TotalAmount: account.TotalAmount,
	}, nil
}

func (s *server) FetchBankAccount(ctx context.Context, in *communication.FetchBankAccountRequest) (*communication.Account, error) {

	log.Printf("FetchBankAccount %v", *in)

	var account Account
	err := getTable(s.Client, in.GetAccountID(), &account)
	if err != nil {
		return nil, err
	}

	if account.UserID != in.GetUserID() {
		return nil, errors.New("Trying to fetc other users account")
	}

	return &communication.Account{
		UserID: account.UserID,
		AccountID: account.AccountID,
		Name: account.Name,
		TotalAmount: account.TotalAmount,
	}, nil
}

func (s *server) DeleteBankAccount(ctx context.Context, in *communication.DeleteBankAccountRequest) (*empty.Empty, error) {

	log.Printf("DeleteBankAccount %v", *in)

	var account Account
	err := getTable(s.Client, in.GetAccountID(), &account)
	if err != nil {
		return nil, err
	}

	if account.UserID != in.GetUserID() {
		return nil, errors.New("Trying to delete account that doesn't belong to given user")
	}

	delTable(s.Client, &account)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}


func (s *server) updateAccount(accountId int32, amount int32) error {

	log.Printf("Updating account %v %v", accountId, amount)

	var account Account
	err := getTable(s.Client, accountId, &account)
	if err != nil {
		return err
	}

	account.TotalAmount = account.TotalAmount + amount
	log.Printf("Updated account %v", account)

	return setTable(s.Client, &account)
}

func (s *server) CreateExpense(ctx context.Context, in *communication.CreateExpenseRequest) (*communication.Expense, error) {

	log.Printf("CreateExpense %v", *in)

	expense := Expense{
		UserID: in.GetUserID(),
		AccountID: in.GetAccountID(),
		Name: in.GetName(),
		ExpenseID: getAutoIncrementID(s.Client, "expenseID"),
		Amount: in.GetAmount(),
	}
	err := setTable(s.Client, &expense)
    if err != nil {
       return nil, err
    }

	s.updateAccount(expense.AccountID, expense.Amount)

	return &communication.Expense{
		Amount: expense.Amount,
		AccountID: expense.AccountID,
		Name: expense.Name,
		ExpenseID: expense.ExpenseID,
	}, nil
}

func (s *server) FetchExpense(ctx context.Context, in *communication.FetchExpenseRequest) (*communication.Expense, error) {
	
	log.Printf("FetchExpense %v", *in)

	var expense Expense
	err := getTable(s.Client, in.GetExpenseID(), &expense)
	if err != nil {
		return nil, err
	}

	if expense.UserID != in.GetUserID() {
		return nil, errors.New("Trying to fetch other users expense")
	}

	return &communication.Expense{
		Amount: expense.Amount,
		AccountID: expense.AccountID,
		Name: expense.Name,
		ExpenseID: expense.ExpenseID,
	}, nil
}

func (s *server) DeleteExpense(ctx context.Context, in *communication.DeleteExpenseRequest) (*empty.Empty, error) {

	log.Printf("DeleteExpense %v", *in)

	var expense Expense
	err := getTable(s.Client, in.GetExpenseID(), &expense)
	if err != nil {
		return nil, err
	}

	if expense.UserID != in.GetUserID() {
		return nil, errors.New("Trying to delete other users expense")
	}

	s.updateAccount(expense.AccountID, -expense.Amount)

	delTable(s.Client, &expense)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func createWorkerServer() server {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	log.Println(pong, err)

	return server{Client: client}
}
