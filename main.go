package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

/////// Types
type User struct {
	id       string
	login    string
	password string
}

type UserCreationData struct {
	login           string `json:"login,omitempty"`
	password        string `json:"password,omitempty"`
	passwordConfirm string `json:"passwordConfirm,omitempty"`
}

/////// Repositories
var users []User

/////// Functions

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var data UserCreationData
	_ = json.NewDecoder(r.Body).Decode(&data)
	log.Println(data)
	var user = CreateUserFromData(data)
	users = append(users, user)
	log.Println(user)
	json.NewEncoder(w).Encode(user)
}

func CreateUserFromData(d UserCreationData) User {
	return User{
		id:       uuid.New(),
		login:    d.login,
		password: d.password,
	}
}

/////// Main
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
