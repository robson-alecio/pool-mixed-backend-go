package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"time"

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

type AddOptionData struct {
	Value string `json:"value,omitempty"`
}

type RemoveOptionData struct {
	Value int `json:"value,omitempty"`
}

type Poll struct {
	ID        string
	CreatedAt time.Time
	Name      string
	Options   []string
	Owner     string
	Published bool
}

type PollVote struct {
	ID           string
	PoolID       string
	UserID       string
	ChosenOption string
}

type PollVoteData struct {
	Value string `json:"value,omitempty"`
}

type PollVoteResult struct {
	VoteID       string
	VoteCounting map[string]float64
}

type AuthenticatedFunction func(b Session, c interface{}) (interface{}, error)

/////// Errors
type ErrPasswordDoNotMatch string

func (e ErrPasswordDoNotMatch) Error() string {
	return string(e)
}

type ErrUserNotFound struct {
	Login string
}

func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("User %s not found.", e.Login)
}

type ErrUserNotLogged string

func (e ErrUserNotLogged) Error() string {
	return string(e)
}

type ErrNotChangePoll string

func (e ErrNotChangePoll) Error() string {
	return string(e)
}

/////// Repositories
var users map[string]User = make(map[string]User)
var usersByLogin map[string]string = make(map[string]string)
var sessions map[string]Session = make(map[string]Session)
var polls map[string]Poll = make(map[string]Poll)
var pollsByUser map[string][]Poll = make(map[string][]Poll)
var votes map[string]map[string][]PollVote = make(map[string]map[string][]PollVote)
var pollVotedByUser map[string]map[string]bool = make(map[string]map[string]bool)

func SaveUser(user User) User {
	log.Println("Saving User", user)
	users[user.ID] = user
	usersByLogin[user.Login] = user.ID
	return user
}

func FindUserByLogin(login string) (string, bool) {
	id, ok := usersByLogin[login]
	return id, ok
}

func FindUserById(id string) User {
	return users[id]
}

func SaveSession(session Session) Session {
	log.Println("Saving Session", session)
	sessions[session.ID] = session
	return session
}

func FindSessionById(id string) (Session, bool) {
	session, ok := sessions[id]
	return session, ok
}

func SavePoll(poll Poll) Poll {
	log.Println("Saving Poll", poll)
	polls[poll.ID] = poll

	_, has := pollsByUser[poll.Owner]

	if !has {
		pollsByUser[poll.Owner] = make([]Poll, 1)
	}
	pollsByUser[poll.Owner] = append(pollsByUser[poll.Owner], poll)

	return poll
}

func UpdatePoll(poll Poll) Poll {
	log.Println("Updating Poll", poll)
	polls[poll.ID] = poll

	return poll
}

func FindPollById(id string) Poll {
	return polls[id]
}

func SaveVote(vote PollVote) PollVote {
	log.Println("Registering vote", vote)
	pollMap, pollVoted := votes[vote.PoolID]

	if !pollVoted {
		pollMap = make(map[string][]PollVote)
		pollVotedByUser[vote.PoolID] = make(map[string]bool)
		votes[vote.PoolID] = pollMap
	}

	_, optionChosen := pollMap[vote.ChosenOption]

	if !optionChosen {
		pollMap[vote.ChosenOption] = make([]PollVote, 0)
	}

	pollMap[vote.ChosenOption] = append(pollMap[vote.ChosenOption], vote)
	pollVotedByUser[vote.PoolID][vote.UserID] = true

	return vote
}

func PollAlreadyVotedByUser(pollId, userId string) bool {
	_, exists := pollVotedByUser[pollId][userId]

	return exists
}

func ExistsOption(pollId, candidate string) bool {
	poll := FindPollById(pollId)

	for _, opt := range poll.Options {
		if opt == candidate {
			return true
		}
	}

	return false
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
		return User{}, ErrPasswordDoNotMatch("Passwords don't match")
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
		return Session{}, ErrUserNotFound{Login: data.Login}
	}

	user := FindUserById(userID)
	log.Println("Result User", user)

	return CreateSession(user), nil
}

func CreateSession(u User) Session {
	return SaveSession(Session{
		ID:     uuid.New(),
		UserID: u.ID,
	})
}

func StartCreatePoll(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session Session, protoData interface{}) (interface{}, error) {
		var data CreatePollData
		mapstructure.Decode(protoData, &data)

		return SavePoll(Poll{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			Name:      data.Name,
			Options:   make([]string, 0),
			Owner:     session.UserID,
		}), nil
	})
}

func AddOption(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session Session, protoData interface{}) (interface{}, error) {
		return ChangePollOrCry(w, r, session, func(poll *Poll) {
			var data AddOptionData
			mapstructure.Decode(protoData, &data)

			poll.Options = append(poll.Options, data.Value)
		})
	})
}

func RemoveOption(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session Session, protoData interface{}) (interface{}, error) {
		return ChangePollOrCry(w, r, session, func(poll *Poll) {
			var data RemoveOptionData
			mapstructure.Decode(protoData, &data)

			before := poll.Options[:data.Value]
			after := poll.Options[data.Value+1:]
			poll.Options[data.Value] = ""
			poll.Options = before
			for _, v := range after {
				poll.Options = append(poll.Options, v)
			}
		})
	})
}

func Publish(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session Session, protoData interface{}) (interface{}, error) {
		return ChangePollOrCry(w, r, session, func(poll *Poll) {
			poll.Published = true
		})
	})
}

func ChangePollOrCry(w http.ResponseWriter, r *http.Request, session Session, fChange func(aPoll *Poll)) (Poll, error) {
	vars := mux.Vars(r)
	poll := FindPollById(vars["id"])

	if poll.Published {
		return Poll{}, ErrNotChangePoll("Can't change a published poll.")
	}

	if poll.Owner != session.UserID {
		return Poll{}, ErrNotChangePoll("Can't change a poll from other user.")
	}

	fChange(&poll)

	return UpdatePoll(poll), nil
}

func CreateVote(w http.ResponseWriter, r *http.Request) {
	session, err := CheckSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	vars := mux.Vars(r)
	pollId := vars["id"]

	var data PollVoteData
	_ = json.NewDecoder(r.Body).Decode(&data)

	if !ExistsOption(pollId, data.Value) {
		message := fmt.Sprintf("There is no option %s for vote on this poll.", data.Value)
		http.Error(w, message, http.StatusConflict)
		return
	}

	if PollAlreadyVotedByUser(pollId, session.UserID) {
		message := fmt.Sprintf("You already voted in this poll.")
		http.Error(w, message, http.StatusConflict)
		return
	}

	vote := PollVote{
		ID:           uuid.New(),
		PoolID:       pollId,
		UserID:       session.UserID,
		ChosenOption: data.Value,
	}

	SaveVote(vote)

	voteResult := PollVoteResult{
		VoteID:       vote.ID,
		VoteCounting: CountVotes(pollId),
	}

	json.NewEncoder(w).Encode(voteResult)
}

func CountVotes(pollId string) map[string]float64 {
	poll := FindPollById(pollId)
	pollMap := votes[pollId]

	total := TotalVotes(pollId)
	result := make(map[string]float64)
	result["total"] = float64(total)

	for _, opt := range poll.Options {
		list, ok := pollMap[opt]

		if ok {
			countVote := len(list)
			perct := float64(countVote*100) / float64(total)
			result[opt] = math.Round(perct*100) / 100
		} else {
			result[opt] = 0
		}
	}

	return result
}

func TotalVotes(pollId string) int {
	sum := 0

	for _, list := range votes[pollId] {
		sum += len(list)
	}

	return sum
}

func GetPoll(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session Session) interface{} {
		vars := mux.Vars(r)
		return FindPollById(vars["id"])
	})
}

func CountingPollVotes(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session Session) interface{} {
		vars := mux.Vars(r)
		return CountVotes(vars["id"])
	})
}

func GetPolls(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session Session) interface{} {
		log.Println("Getting all")
		pollsSorted := make([]Poll, 0)
		for _, poll := range polls {
			pollsSorted = append(pollsSorted, poll)
		}

		return SortPolls(pollsSorted)
	})
}

func GetPollsMine(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session Session) interface{} {
		pollsMineSorted := make([]Poll, 0)
		for _, poll := range polls {
			if poll.Owner == session.UserID {
				pollsMineSorted = append(pollsMineSorted, poll)
			}
		}

		return SortPolls(pollsMineSorted)
	})
}

func SortPolls(polls []Poll) []Poll {
	sort.Slice(polls, func(i, j int) bool {
		return polls[i].CreatedAt.Before(polls[j].CreatedAt)
	})

	return polls
}

func ExecuteSessioned(w http.ResponseWriter, r *http.Request, f func(session Session) interface{}) {
	session, err := CheckSession(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	result := f(session)

	json.NewEncoder(w).Encode(result)
}

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

func CheckAuthentication(w http.ResponseWriter, r *http.Request) (Session, error) {
	session, error := CheckSession(w, r)

	if error != nil {
		return session, error
	}

	user := FindUserById(session.UserID)

	if user.Password == "" {
		return Session{}, ErrUserNotLogged("Must be logged to perform this action. Not authenticated.")
	}

	return session, nil
}

func CheckSession(w http.ResponseWriter, r *http.Request) (Session, error) {
	sessionID := r.Header.Get("sessionId")

	if sessionID == "" {
		return Session{}, ErrUserNotLogged("Must be logged to perform this action. Missing value.")
	}

	session, ok := FindSessionById(sessionID)

	if !ok {
		return Session{}, ErrUserNotLogged("Must be logged to perform this action. Session invalid.")
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
