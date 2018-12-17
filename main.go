package main

import (
	"database/sql"
	"log"
	"net/http"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/gorilla/mux"

	. "github/RobsonAlecio/pool-mixed-backend-go/app"
)

/////// Repositories
var db *sql.DB
var userHandler *UserHandlerImpl
var sessionHandler *SessionHandlerImpl
var pollHandler *PollHandlerImpl
var pollOptionHandler *PollOptionHandlerImpl
var pollVoteHandler *PollVoteHandlerImpl

//CreateUserEndpointEntry ...
func CreateUserEndpointEntry(w http.ResponseWriter, r *http.Request) {
	CreateUser(createHTTPHelper(w, r), userHandler)
}

//VisitEndpointEntry ...
func VisitEndpointEntry(w http.ResponseWriter, r *http.Request) {
	Visit(createHTTPHelper(w, r), userHandler, sessionHandler)
}

//LoginEndpointEntry ...
func LoginEndpointEntry(w http.ResponseWriter, r *http.Request) {
	Login(createHTTPHelper(w, r), userHandler, sessionHandler)
}

//StartCreatePollEndpointEntry ...
func StartCreatePollEndpointEntry(w http.ResponseWriter, r *http.Request) {
	StartCreatePoll(createHTTPHelper(w, r), pollHandler)
}

//AddOptionEndpointEntry ...
func AddOptionEndpointEntry(w http.ResponseWriter, r *http.Request) {
	AddOption(createHTTPHelper(w, r), pollHandler, pollOptionHandler)
}

//RemoveOptionEndpointEntry ...
func RemoveOptionEndpointEntry(w http.ResponseWriter, r *http.Request) {
	RemoveOption(createHTTPHelper(w, r), pollHandler, pollOptionHandler)
}

//PublishEndpointEntry ...
func PublishEndpointEntry(w http.ResponseWriter, r *http.Request) {
	Publish(createHTTPHelper(w, r), pollHandler, pollOptionHandler)
}

//CreateVoteEndpointEntry ...
func CreateVoteEndpointEntry(w http.ResponseWriter, r *http.Request) {
	CreateVote(createHTTPHelper(w, r), pollOptionHandler, pollVoteHandler)
}

//GetPoll ...
func GetPoll(w http.ResponseWriter, r *http.Request) {
	// ExecuteSessioned(w, r, func(session *Session) interface{} {
	// 	vars := mux.Vars(r)
	// 	id, err := kallax.NewULIDFromText(vars["id"])
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return nil
	// 	}

	// 	poll, errFind := FindPollByID(id)
	// 	if errFind != nil {
	// 		http.Error(w, errFind.Error(), http.StatusBadRequest)
	// 		return nil
	// 	}

	// 	return poll
	// })
}

//CountingPollVotes ...
func CountingPollVotes(w http.ResponseWriter, r *http.Request) {
	// ExecuteSessioned(w, r, func(session *Session) interface{} {
	// 	vars := mux.Vars(r)
	// 	id, err := kallax.NewULIDFromText(vars["id"])
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return make(map[string]float64)
	// 	}
	// 	return CountVotes(id)
	// })
}

//GetPolls ...
func GetPolls(w http.ResponseWriter, r *http.Request) {
	// ExecuteSessioned(w, r, func(session *Session) interface{} {
	// 	log.Println("Getting all")
	// 	query := NewPollQuery().
	// 		Order(kallax.Asc(Schema.Poll.CreatedAt))
	// 	return FindPolls(w, query)
	// })
}

//GetPollsMine ...
func GetPollsMine(w http.ResponseWriter, r *http.Request) {
	// ExecuteSessioned(w, r, func(session *Session) interface{} {
	// 	query := NewPollQuery().
	// 		FindByOwner(session.UserID).
	// 		Order(kallax.Asc(Schema.Poll.CreatedAt))

	// 	return FindPolls(w, query)
	// })
}

// FindPolls ...
func FindPolls(w http.ResponseWriter, query *PollQuery) []*Poll {
	// polls, err := pollStore.FindAll(query)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return nil
	// }

	// return polls
	return make([]*Poll, 0)
}

//ConnectToDatabase ...
func ConnectToDatabase() {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=poll password=poll dbname=poll sslmode=disable")

	if err != nil {
		panic(err)
	}

	userHandler = NewUserHandler(db)
	sessionHandler = NewSessionHandler(db)
	pollOptionHandler = NewPollOptionHandler(db)
	pollHandler = NewPollHandler(db, pollOptionHandler)
	pollVoteHandler = NewPollVoteHandler(db)

	log.Println("Successfuly connected!")
}

func createHTTPHelper(w http.ResponseWriter, r *http.Request) *HTTPHelperImpl {
	helper := NewHTTPHelper(w, r)

	helper.CheckSession = func(ID string) error {
		realID, err := kallax.NewULIDFromText(ID)
		if err != nil {
			return err
		}

		session, errFind := sessionHandler.FindSessionByID(realID)
		helper.Session = session
		return errFind
	}

	return helper
}

// ConfigStartServer ...
func ConfigStartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUserEndpointEntry).Methods("POST")

	router.HandleFunc("/visit", VisitEndpointEntry).Methods("POST")
	router.HandleFunc("/login", LoginEndpointEntry).Methods("POST")

	router.HandleFunc("/polls", StartCreatePollEndpointEntry).Methods("POST")
	router.HandleFunc("/polls/{id}", AddOptionEndpointEntry).Methods("PUT")
	router.HandleFunc("/polls/{id}", RemoveOptionEndpointEntry).Methods("DELETE")
	router.HandleFunc("/polls/{id}/publish", PublishEndpointEntry).Methods("PUT")
	router.HandleFunc("/polls/{id}/vote", CreateVoteEndpointEntry).Methods("POST")
	router.HandleFunc("/polls/{id}", GetPoll).Methods("GET")
	router.HandleFunc("/polls/{id}/counting", CountingPollVotes).Methods("GET")
	router.HandleFunc("/polls", GetPolls).Methods("GET")
	router.HandleFunc("/mine/polls", GetPollsMine).Methods("GET")

	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":8000", router))
}

//main ...
func main() {
	ConnectToDatabase()
	ConfigStartServer()
}

// TODOs (Improvements)
// Sessions to expires
// Sessions efemerals
// Endpoint for published polls
// Poll DTO for GETs
// Split files by packages
// Poll Options with defined order
