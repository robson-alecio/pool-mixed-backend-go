package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"

	. "github/RobsonAlecio/pool-mixed-backend-go/app"
)

// AuthenticatedFunction ...
type AuthenticatedFunction func(session *Session, data interface{}) (interface{}, error)

/////// Repositories
var db *sql.DB
var userHandler *UserHandlerImpl
var sessionHandler *SessionHandlerImpl
var pollStore *PollStore
var pollOptionStore *PollOptionStore
var pollVoteStore *PollVoteStore

//SavePoll ...
func SavePoll(poll *Poll) (bool, error) {
	log.Println("Saving Poll", poll)

	return pollStore.Save(poll)
}

// SavePollOption ...
func SavePollOption(pollOption *PollOption) (bool, error) {
	log.Println("Adding Poll Option", pollOption)

	return pollOptionStore.Save(pollOption)
}

// DeletePollOption ...
func DeletePollOption(id kallax.ULID) error {
	log.Println("Removing Poll Option", id)

	query := NewPollOptionQuery().FindByID(id)

	opt, err := pollOptionStore.FindOne(query)

	if err != nil {
		return err
	}

	return pollOptionStore.Delete(opt)
}

//FindPollByID ...
func FindPollByID(ID kallax.ULID) (*Poll, error) {
	query := NewPollQuery().FindByID(ID)
	poll, err := pollStore.FindOne(query)

	if err != nil {
		return poll, err
	}

	options, errOption := FindPollOptions(ID)
	if errOption != nil {
		return poll, errOption
	}

	poll.Options = options

	return poll, nil
}

// FindPollOptions ...
func FindPollOptions(id kallax.ULID) ([]*PollOption, error) {
	query := NewPollOptionQuery().FindByOwner(id)
	return pollOptionStore.FindAll(query)
}

//SaveVote ...
func SaveVote(vote PollVote) PollVote {
	log.Println("Registering vote", vote)

	pollVoteStore.Save(&vote)

	return vote
}

//PollAlreadyVotedByUser ...
func PollAlreadyVotedByUser(pollID, userID kallax.ULID) (bool, error) {
	query := NewPollVoteQuery().
		FindByPollID(pollID).
		FindByUserID(userID)

	count, err := pollVoteStore.Count(query)

	return count > 0, err
}

//ExistsOption ...
func ExistsOption(pollID kallax.ULID, candidate string) (bool, error) {
	query := NewPollOptionQuery().
		FindByOwner(pollID).
		FindByContent(candidate)

	count, err := pollOptionStore.Count(query)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

//CreateUserEndpointEntry ...
func CreateUserEndpointEntry(w http.ResponseWriter, r *http.Request) {
	CreateUser(NewHTTPHelper(w, r), userHandler)
}

//VisitEndpointEntry ...
func VisitEndpointEntry(w http.ResponseWriter, r *http.Request) {
	Visit(NewHTTPHelper(w, r), userHandler, sessionHandler)
}

//LoginEndpointEntry ...
func LoginEndpointEntry(w http.ResponseWriter, r *http.Request) {
	Login(NewHTTPHelper(w, r), userHandler, sessionHandler)
}

//StartCreatePollEndpointEntry ...
func StartCreatePollEndpointEntry(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session *Session, protoData interface{}) (interface{}, error) {
		var data CreatePollData
		mapstructure.Decode(protoData, &data)

		return SavePoll(&Poll{
			ID:      kallax.NewULID(),
			Name:    data.Name,
			Options: make([]*PollOption, 0),
			Owner:   session.UserID,
		})
	})
}

//AddOption ...
func AddOption(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session *Session, protoData interface{}) (interface{}, error) {
		return ChangePollOrCry(w, r, session, func(poll *Poll) {
			var data AddOptionData
			mapstructure.Decode(protoData, &data)

			SavePollOption(&PollOption{
				ID:      kallax.NewULID(),
				Owner:   poll,
				Content: data.Value,
			})
		})
	})
}

//RemoveOption ...
func RemoveOption(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session *Session, protoData interface{}) (interface{}, error) {
		return ChangePollOrCry(w, r, session, func(poll *Poll) {
			var data RemoveOptionData
			mapstructure.Decode(protoData, &data)

			id, err := kallax.NewULIDFromText(data.Value)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			DeletePollOption(id)
		})
	})
}

//Publish ...
func Publish(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session *Session, protoData interface{}) (interface{}, error) {
		return ChangePollOrCry(w, r, session, func(poll *Poll) {
			poll.Published = true
		})
	})
}

//ChangePollOrCry ...
func ChangePollOrCry(w http.ResponseWriter, r *http.Request, session *Session, fChange func(aPoll *Poll)) (bool, error) {
	vars := mux.Vars(r)
	id, err := kallax.NewULIDFromText(vars["id"])
	if err != nil {
		return false, ErrNotChangePoll(err.Error())
	}

	poll, errFind := FindPollByID(id)

	if errFind != nil {
		http.Error(w, errFind.Error(), http.StatusBadRequest)
	}

	if poll.Published {
		return false, ErrNotChangePoll("Can't change a published poll.")
	}

	if poll.Owner != session.UserID {
		return false, ErrNotChangePoll("Can't change a poll from other user.")
	}

	fChange(poll)

	return SavePoll(poll)
}

//CreateVote ...
func CreateVote(w http.ResponseWriter, r *http.Request) {
	session, err := CheckSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	vars := mux.Vars(r)
	pollID, err := kallax.NewULIDFromText(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data PollVoteData
	_ = json.NewDecoder(r.Body).Decode(&data)

	exists, err := ExistsOption(pollID, data.Value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if !exists {
		message := fmt.Sprintf("There is no option %s for vote on this poll.", data.Value)
		http.Error(w, message, http.StatusConflict)
		return
	}

	voted, errVoted := PollAlreadyVotedByUser(pollID, session.UserID)

	if errVoted != nil {
		http.Error(w, errVoted.Error(), http.StatusConflict)
		return
	}

	if voted {
		message := fmt.Sprintf("You already voted in this poll.")
		http.Error(w, message, http.StatusConflict)
		return
	}

	vote := PollVote{
		ID:           kallax.NewULID(),
		PollID:       pollID,
		UserID:       session.UserID,
		ChosenOption: data.Value,
	}

	SaveVote(vote)

	voteResult := PollVoteResult{
		VoteID:       vote.ID.String(),
		VoteCounting: CountVotes(pollID),
	}

	json.NewEncoder(w).Encode(voteResult)
}

//CountVotes ...
func CountVotes(pollID kallax.ULID) map[string]float64 {
	options, err := FindPollOptions(pollID)

	if err != nil {
		badResult := make(map[string]float64)
		badResult[err.Error()] = -1.0
		return badResult
	}

	count := make(map[string]int64)

	for _, opt := range options {
		query := NewPollVoteQuery().FindByPollID(pollID).FindByChosenOption(opt.Content)
		votesOption, err := pollVoteStore.Count(query)
		if err != nil {
			count[opt.Content] = 0
		} else {
			count[opt.Content] = votesOption
		}
	}

	total := int64(0)

	for _, votes := range count {
		total += votes
	}

	result := make(map[string]float64)
	result["total"] = float64(total)

	for _, opt := range options {
		countVote, ok := count[opt.Content]

		if ok {
			perct := float64(countVote*100) / float64(total)
			result[opt.Content] = math.Round(perct*100) / 100
		} else {
			result[opt.Content] = 0
		}
	}

	return result
}

//GetPoll ...
func GetPoll(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session *Session) interface{} {
		vars := mux.Vars(r)
		id, err := kallax.NewULIDFromText(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}

		poll, errFind := FindPollByID(id)
		if errFind != nil {
			http.Error(w, errFind.Error(), http.StatusBadRequest)
			return nil
		}

		return poll
	})
}

//CountingPollVotes ...
func CountingPollVotes(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session *Session) interface{} {
		vars := mux.Vars(r)
		id, err := kallax.NewULIDFromText(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return make(map[string]float64)
		}
		return CountVotes(id)
	})
}

//GetPolls ...
func GetPolls(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session *Session) interface{} {
		log.Println("Getting all")
		query := NewPollQuery().
			Order(kallax.Asc(Schema.Poll.CreatedAt))
		return FindPolls(w, query)
	})
}

//GetPollsMine ...
func GetPollsMine(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session *Session) interface{} {
		query := NewPollQuery().
			FindByOwner(session.UserID).
			Order(kallax.Asc(Schema.Poll.CreatedAt))

		return FindPolls(w, query)
	})
}

// FindPolls ...
func FindPolls(w http.ResponseWriter, query *PollQuery) []*Poll {
	polls, err := pollStore.FindAll(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return polls
}

//ExecuteSessioned ...
func ExecuteSessioned(w http.ResponseWriter, r *http.Request, f func(session *Session) interface{}) {
	session, err := CheckSession(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	result := f(session)

	json.NewEncoder(w).Encode(result)
}

//ExecuteAuthenticated ...
func ExecuteAuthenticated(w http.ResponseWriter, r *http.Request, f AuthenticatedFunction) {
	session, err := CheckAuthentication(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var data interface{}
	_ = json.NewDecoder(r.Body).Decode(&data)

	result, err := f(session, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(result)
}

//CheckAuthentication ...
func CheckAuthentication(w http.ResponseWriter, r *http.Request) (*Session, error) {
	session, errCheck := CheckSession(w, r)

	if errCheck != nil {
		return nil, errCheck
	}

	if session.RegisteredUser {
		return nil, ErrUserNotLogged("Must be logged to perform this action. Not authenticated.")
	}

	return session, nil
}

//CheckSession ...
func CheckSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	sessionID := r.Header.Get("sessionId")

	if sessionID == "" {
		return nil, ErrUserNotLogged("Must be logged to perform this action. Missing value.")
	}

	ID, err := kallax.NewULIDFromText(sessionID)
	if err != nil {
		return nil, err
	}

	return sessionHandler.FindSessionByID(ID)
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
	pollStore = NewPollStore(db)
	pollOptionStore = NewPollOptionStore(db)
	pollVoteStore = NewPollVoteStore(db)

	log.Println("Successfuly connected!")
}

// ConfigStartServer ...
func ConfigStartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUserEndpointEntry).Methods("POST")

	router.HandleFunc("/visit", VisitEndpointEntry).Methods("POST")
	router.HandleFunc("/login", LoginEndpointEntry).Methods("POST")

	router.HandleFunc("/polls", StartCreatePollEndpointEntry).Methods("POST")
	router.HandleFunc("/polls/{id}", AddOption).Methods("PUT")
	router.HandleFunc("/polls/{id}", RemoveOption).Methods("DELETE")
	router.HandleFunc("/polls/{id}/publish", Publish).Methods("PUT")
	router.HandleFunc("/polls/{id}/vote", CreateVote).Methods("POST")
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
