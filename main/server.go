package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/MaksymYesipov/hasher/hasher"
	"github.com/google/uuid"
	"net/http"
)

type User struct {
	UserName string
	Password string
}

type UserBean struct {
	Id       string
	UserName string
}

type LoginResponse struct {
	Url string
}

const url = "ws://fancy-chat.io/ws&token=%s"
const tokenLength = 64

var users = make(map[string]User)

func createUser(w http.ResponseWriter, req *http.Request) {
	var u User

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, v := range users {
		if v.UserName == u.UserName {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}
	}

	u.Password, err = hasher.HashPassword(u.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uuidValue := uuid.New().String()
	users[uuidValue] = u

	userBean := UserBean{Id: uuidValue, UserName: u.UserName}
	response, _ := json.Marshal(userBean)
	fmt.Fprint(w, string(response))
}

func login(w http.ResponseWriter, req *http.Request) {
	var u User

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range users {
		if v.UserName == u.UserName {
			if !hasher.CheckPasswordHash(u.Password, v.Password) {
				http.Error(w, "Invalid login or password", http.StatusUnauthorized)
				return
			}
			responseData, _ := json.Marshal(LoginResponse{Url: fmt.Sprintf(url, generateAccessToken())})
			fmt.Fprint(w, string(responseData))
			return
		}
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
	}
}

func withJsonMimeType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

func generateAccessToken() string {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func main() {
	http.Handle("/user", withJsonMimeType(http.HandlerFunc(createUser)))
	http.Handle("/user/login", withJsonMimeType(http.HandlerFunc(login)))

	http.ListenAndServe(":8090", nil)
}
