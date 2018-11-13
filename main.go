package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

/////// Types
type User struct {
	ID       string `json:"ID,omitempty"`
	Login    string `json:"login,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserCreationData struct {
	Login           string `json:"login,omitempty"`
	Name            string `json:"name,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
}

type Session struct {
	ID     string
	UserID string
}

type LoginData struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

/////// Errors

type ErrPasswordDoNotMatch struct {
	message string
}

func (e *ErrPasswordDoNotMatch) Error() string {
	return e.message
}

type ErrUserNotFound struct {
	message string
}

func (e *ErrUserNotFound) Error() string {
	return e.message
}

/////// Repositories
var users map[string]User = make(map[string]User)
var usersByLogin map[string]string = make(map[string]string)
var sessions map[string]Session = make(map[string]Session)

func SaveUser(user User) {
	log.Println("Saving User", user)
	users[user.ID] = user
	usersByLogin[user.Login] = user.ID
}

func SaveSession(session Session) {
	log.Println("Saving Session", session)
	sessions[session.ID] = session
}

/////// Functions
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var data UserCreationData
	_ = json.NewDecoder(r.Body).Decode(&data)
	user, err := CreateUserFromData(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SaveUser(user)
	json.NewEncoder(w).Encode(user)
}

func CreateUserFromData(d UserCreationData) (User, error) {
	if d.Password != d.PasswordConfirm {
		return User{}, &ErrPasswordDoNotMatch{"Passwords don't match"}
	}

	user := User{
		ID:       uuid.New(),
		Login:    d.Login,
		Name:     d.Name,
		Password: d.Password,
	}
	return user, nil
}

func IsAnon(u User) bool {
	return u.Password == ""
}

func Visit(w http.ResponseWriter, r *http.Request) {
	user := CreateAnonUser()
	session := CreateSession(user)
	json.NewEncoder(w).Encode(session)
}

func CreateAnonUser() User {
	user := User{
		ID: uuid.New(),
	}
	user.Name = "Anon" + user.ID
	user.Login = user.Name

	SaveUser(user)

	return user
}

func Login(w http.ResponseWriter, r *http.Request) {
	var data LoginData
	_ = json.NewDecoder(r.Body).Decode(&data)
	session, err := Authenticate(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(session)
}

func Authenticate(data LoginData) (Session, error) {
	log.Println("Trying authenticate", data)
	userID, ok := usersByLogin[data.Login]

	log.Println("Result ID", userID, ok)
	if !ok {
		return Session{}, &ErrUserNotFound{fmt.Sprintf("User %s not found.", data.Login)}
	}

	user := users[userID]
	log.Println("Result User", user)

	return CreateSession(user), nil
}

func CreateSession(u User) Session {
	session := Session{
		ID:     uuid.New(),
		UserID: u.ID,
	}

	SaveSession(session)

	return session
}

/////// Main
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")

	router.HandleFunc("/visit", Visit).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")

	// router.HandleFunc("/polls", StartCreatePoll).Methods("POST")
	// router.HandleFunc("/polls", AddOption).Methods("PUT")
	// router.HandleFunc("/polls", RemoveOption).Methods("PUT")
	// router.HandleFunc("/polls", Finish).Methods("PUT")
	// router.HandleFunc("/polls", CreateVote).Methods("POST")
	// router.HandleFunc("/polls", GetPolls).Methods("GET")

	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":8000", router))
}
