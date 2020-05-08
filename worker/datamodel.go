package main

import (
	"strconv"
	"encoding/json"

	"github.com/go-redis/redis/v7"
)

func getAutoIncrementID(client *redis.Client, name string) int32 {

	id, err := client.Get(name).Int()
	if err != nil {
		id = 1
	}
	client.Set(name, id+1, 0) // TODO: Check for errors
	return int32(id)
}

type Table interface {
	ID() int32
	TableName() string
}

func setTable(client *redis.Client, obj Table) error {

	objJson, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return client.HSet(obj.TableName(), obj.ID(), objJson).Err()
}

func delTable(client *redis.Client, obj Table) error {

	_, err := client.HDel(obj.TableName(), strconv.Itoa(int(obj.ID()))).Result()
	return err
}

func getTable(client *redis.Client, id int32, obj Table) error {

	objJson, err := client.HGet(obj.TableName(), strconv.Itoa(int(id))).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(objJson), obj)
}

type Account struct{
	UserID int32
	AccountID int32
	Name string
	TotalAmount int32
}

func (a *Account) ID() int32{ 

	return a.AccountID
}

func (a *Account) TableName() string {

	return "account"
}

type Expense struct {
	UserID int32
	AccountID int32
	Name string
	Amount int32
	ExpenseID int32
}

func (a *Expense) ID() int32 {

	return a.ExpenseID
}

func (a *Expense) TableName() string {

	return "expense"
}

