package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
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

type CreatePollData struct {
	Name string `json:"name,omitempty"`
}

type Poll struct {
	ID      string
	Name    string
	Options []string
	Owner   string
}

type AuthenticatedFunction func(b Session, c interface{}) interface{}

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

type ErrUserNotLogged struct {
	message string
}

func (e *ErrUserNotLogged) Error() string {
	return e.message
}

/////// Repositories
var users map[string]User = make(map[string]User)
var usersByLogin map[string]string = make(map[string]string)
var sessions map[string]Session = make(map[string]Session)
var polls map[string]Poll = make(map[string]Poll)
var pollsByUser map[string][]Poll = make(map[string][]Poll)

func SaveUser(user User) {
	log.Println("Saving User", user)
	users[user.ID] = user
	usersByLogin[user.Login] = user.ID
}

func FindUserByLogin(login string) (string, bool) {
	id, ok := usersByLogin[login]
	return id, ok
}

func FindUserById(id string) User {
	return users[id]
}

func SaveSession(session Session) {
	log.Println("Saving Session", session)
	sessions[session.ID] = session
}

func FindSessionById(id string) (Session, bool) {
	session, ok := sessions[id]
	return session, ok
}

func SavePoll(poll Poll) {
	log.Println("Saving Poll", poll)
	polls[poll.ID] = poll

	_, has := pollsByUser[poll.Owner]

	if !has {
		pollsByUser[poll.Owner] = make([]Poll, 1)
	}
	pollsByUser[poll.Owner] = append(pollsByUser[poll.Owner], poll)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(session)
}

func Authenticate(data LoginData) (Session, error) {
	log.Println("Trying authenticate", data)
	userID, ok := FindUserByLogin(data.Login)

	log.Println("Result ID", userID, ok)
	if !ok {
		return Session{}, &ErrUserNotFound{fmt.Sprintf("User %s not found.", data.Login)}
	}

	user := FindUserById(userID)
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

func StartCreatePoll(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, DoCreatePoll)
}

func DoCreatePoll(session Session, protoData interface{}) interface{} {
	var data CreatePollData
	mapstructure.Decode(protoData, &data)

	poll := Poll{
		ID:      uuid.New(),
		Name:    data.Name,
		Options: make([]string, 0),
		Owner:   session.UserID,
	}

	SavePoll(poll)

	return poll
}

func ExecuteAuthenticated(w http.ResponseWriter, r *http.Request, f AuthenticatedFunction) {
	session, err := CheckAuthentication(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var data interface{}
	_ = json.NewDecoder(r.Body).Decode(&data)

	result := f(session, data)

	json.NewEncoder(w).Encode(result)
}

func CheckAuthentication(w http.ResponseWriter, r *http.Request) (Session, error) {
	sessionID := r.Header.Get("sessionId")

	log.Println(sessionID)
	if sessionID == "" {
		return Session{}, &ErrUserNotLogged{"Must be logged to perform this action. Missing value."}
	}

	session, ok := FindSessionById(sessionID)

	log.Println(session, ok)
	if !ok {
		return Session{}, &ErrUserNotLogged{"Must be logged to perform this action. Session invalid."}
	}

	return session, nil
}

/////// Main
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")

	router.HandleFunc("/visit", Visit).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")

	router.HandleFunc("/polls", StartCreatePoll).Methods("POST")
	// router.HandleFunc("/polls/{id}", AddOption).Methods("PUT")
	// router.HandleFunc("/polls/{id}", RemoveOption).Methods("PUT")
	// router.HandleFunc("/polls/{id}", Finish).Methods("PUT")
	// router.HandleFunc("/polls/{id}", CreateVote).Methods("POST")
	// router.HandleFunc("/polls/{id}", GetPolls).Methods("GET")
	// router.HandleFunc("/polls", GetPolls).Methods("GET")
	// router.HandleFunc("/polls/mine", GetPolls).Methods("GET")

	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":8000", router))
}
