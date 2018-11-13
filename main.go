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
	Id       string `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserCreationData struct {
	Login           string `json:"login,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
}

type ErrPasswordDoNotMatch struct {
	message string
}

func (e *ErrPasswordDoNotMatch) Error() string {
	return e.message
}

/////// Repositories
var users []User

/////// Functions

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var data UserCreationData
	_ = json.NewDecoder(r.Body).Decode(&data)
	user, err := CreateUserFromData(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

func CreateUserFromData(d UserCreationData) (User, error) {
	if d.Password != d.PasswordConfirm {
		return User{}, &ErrPasswordDoNotMatch{"Passwords don't match"}
	}

	user := User{
		Id:       uuid.New(),
		Login:    d.Login,
		Password: d.Password,
	}
	return user, nil
}

/////// Main
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")
	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":8000", router))
}
