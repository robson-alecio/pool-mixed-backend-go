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
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserCreationData struct {
	Login           string `json:"login,omitempty"`
	Name            string `json:"name,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
}

type ErrPasswordDoNotMatch struct {
	message string
}

func (e *ErrPasswordDoNotMatch) Error() string {
	return e.message
}

type Session struct {
	Id     string
	UserId string
}

/////// Repositories
var users map[string]User = make(map[string]User)
var sessions map[string]Session = make(map[string]Session)

func SaveUser(user User) {
	users[user.Id] = user
}

func SaveSession(session Session) {
	sessions[session.Id] = session
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

	users[user.Id] = user
	json.NewEncoder(w).Encode(user)
}

func CreateUserFromData(d UserCreationData) (User, error) {
	if d.Password != d.PasswordConfirm {
		return User{}, &ErrPasswordDoNotMatch{"Passwords don't match"}
	}

	user := User{
		Id:       uuid.New(),
		Login:    d.Login,
		Name:     d.Name,
		Password: d.Password,
	}
	return user, nil
}

func IsAnon(u User) bool {
	return u.Login == ""
}

func Visit(w http.ResponseWriter, r *http.Request) {
	user := CreateAnonUser()
	session := CreateSession(user)
	json.NewEncoder(w).Encode(session)
}

func CreateAnonUser() User {
	user := User{
		Id: uuid.New(),
	}
	user.Name = "Anon" + user.Id

	SaveUser(user)

	return user
}

func CreateSession(u User) Session {
	session := Session{
		Id:     uuid.New(),
		UserId: u.Id,
	}

	SaveSession(session)

	return session
}

/////// Main
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")

	router.HandleFunc("/visit", Visit).Methods("POST")
	// router.HandleFunc("/login", Login).Methods("POST")

	// router.HandleFunc("/polls", StartCreatePoll).Methods("POST")
	// router.HandleFunc("/polls", AddOption).Methods("PUT")
	// router.HandleFunc("/polls", RemoveOption).Methods("PUT")
	// router.HandleFunc("/polls", Finish).Methods("PUT")
	// router.HandleFunc("/polls", CreateVote).Methods("POST")
	// router.HandleFunc("/polls", GetPolls).Methods("GET")

	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":8000", router))
}
