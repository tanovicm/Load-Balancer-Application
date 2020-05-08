package main

import (
	"net/http"
	"time"
	"encoding/json"

	"github.com/go-redis/redis/v7"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	gcontext "github.com/gorilla/context"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

type Claims struct {
	UserID int32
	jwt.StandardClaims
}

func auth(endpoint func(http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			http.Error(w, "No token header", http.StatusBadRequest)
			return
		}

		tokenStr := r.Header["Token"][0]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !token.Valid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		gcontext.Set(r, "userId", claims.UserID)
		endpoint(w, r)
	}
}

func GenerateJWT(userID int32) (string, error) {

	claims := &Claims{
		UserID: userID, 
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(300 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

type Credentials struct {
	Password string 
	Username string
}

type User struct {
	Password string
	Username string
	UserID int32
}

func register(client *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var request Credentials
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		exsists, err := client.HExists("user", request.Username).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if exsists {
			http.Error(w, "Username taken", http.StatusBadRequest)
			return
		}

		userID, err := client.Get("userID").Int()
		if err != nil {
			userID = 1
		}
		client.Set("userID", userID+1, 0) // TODO: Check for errors

		user := User{
			Password: request.Password,
			Username: request.Username,
			UserID: int32(userID),
		}
		userJson, err := json.Marshal(&user) // TODO: Check for errors
		client.HSet("user", request.Username, userJson).Result() // TODO: Check error

		json.NewEncoder(w).Encode(user)
	}
}

func login(client *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		
		var request Credentials
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userJson, err := client.HGet("user", request.Username).Result()
		if err != nil {
			http.Error(w, "Wrong username or password", http.StatusBadRequest)
			return
		}

		var user User
		err = json.Unmarshal([]byte(userJson), &user) // TODO:Check error

		if user.Password != request.Password {
			http.Error(w, "Wrong username or password", http.StatusBadRequest)
			return
		}

		validToken, err := GenerateJWT(user.UserID)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusBadRequest)
			return
		}

		w.Write([]byte(validToken))
	}
}

func registerAuthRoutes(r *mux.Router) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r.HandleFunc("/login", login(client))
	r.HandleFunc("/register", register(client))
}