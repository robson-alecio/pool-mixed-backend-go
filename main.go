package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

type AuthenticatedFunction func(b Session, c interface{}) (interface{}, error)

/////// Repositories
var db *sql.DB
var userStore *UserStore
var users map[kallax.ULID]User = make(map[kallax.ULID]User)
var usersByLogin map[string]kallax.ULID = make(map[string]kallax.ULID)
var sessions map[kallax.ULID]Session = make(map[kallax.ULID]Session)
var polls map[kallax.ULID]Poll = make(map[kallax.ULID]Poll)
var pollsByUser map[kallax.ULID][]Poll = make(map[kallax.ULID][]Poll)
var votes map[kallax.ULID]map[string][]PollVote = make(map[kallax.ULID]map[string][]PollVote)
var pollVotedByUser map[kallax.ULID]map[kallax.ULID]bool = make(map[kallax.ULID]map[kallax.ULID]bool)

//SaveUser ...
func SaveUser(user User) User {
	log.Println("Saving User", user)

	userStore.Save(&user)

	usersByLogin[user.Login] = user.ID
	return user
}

//FindUserByLogin ...
func FindUserByLogin(login string) (kallax.ULID, bool) {
	id, ok := usersByLogin[login]
	return id, ok
}

//FindUserById ...
func FindUserById(id kallax.ULID) User {
	return users[id]
}

//SaveSession ...
func SaveSession(session Session) Session {
	log.Println("Saving Session", session)
	sessions[session.ID] = session
	return session
}

//FindSessionById ...
func FindSessionById(id kallax.ULID) (Session, bool) {
	session, ok := sessions[id]
	return session, ok
}

//SavePoll ...
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

//UpdatePoll ...
func UpdatePoll(poll Poll) Poll {
	log.Println("Updating Poll", poll)
	polls[poll.ID] = poll

	return poll
}

//FindPollById ...
func FindPollById(id kallax.ULID) Poll {
	return polls[id]
}

//SaveVote ...
func SaveVote(vote PollVote) PollVote {
	log.Println("Registering vote", vote)
	pollMap, pollVoted := votes[vote.PoolID]

	if !pollVoted {
		pollMap = make(map[string][]PollVote)
		pollVotedByUser[vote.PoolID] = make(map[kallax.ULID]bool)
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

//PollAlreadyVotedByUser ...
func PollAlreadyVotedByUser(pollID, userID kallax.ULID) bool {
	_, exists := pollVotedByUser[pollID][userID]

	return exists
}

//ExistsOption ...
func ExistsOption(pollID kallax.ULID, candidate string) bool {
	poll := FindPollById(pollID)

	for _, opt := range poll.Options {
		if opt == candidate {
			return true
		}
	}

	return false
}

//CreateUser ...
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

//CreateUserFromData ...
func CreateUserFromData(d UserCreationData) (User, error) {
	if d.Password != d.PasswordConfirm {
		return User{}, ErrPasswordDoNotMatch("Passwords don't match")
	}

	user := User{
		ID:       kallax.NewULID(),
		Login:    d.Login,
		Name:     d.Name,
		Password: d.Password,
	}
	return user, nil
}

//IsAnon ...
func IsAnon(u User) bool {
	return u.Password == ""
}

//Visit ...
func Visit(w http.ResponseWriter, r *http.Request) {
	user := CreateAnonUser()
	session := CreateSession(user)
	json.NewEncoder(w).Encode(session)
}

//CreateAnonUser ...
func CreateAnonUser() User {
	user := User{
		ID: kallax.NewULID(),
	}
	user.Name = "Anon" + user.ID.String()
	user.Login = user.Name

	SaveUser(user)

	return user
}

//Login ...
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

//Authenticate ...
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

//CreateSession ...
func CreateSession(u User) Session {
	return SaveSession(Session{
		ID:     kallax.NewULID(),
		UserID: u.ID,
	})
}

//StartCreatePoll ...
func StartCreatePoll(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session Session, protoData interface{}) (interface{}, error) {
		var data CreatePollData
		mapstructure.Decode(protoData, &data)

		return SavePoll(Poll{
			ID:      kallax.NewULID(),
			Name:    data.Name,
			Options: make([]string, 0),
			Owner:   session.UserID,
		}), nil
	})
}

//AddOption ...
func AddOption(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session Session, protoData interface{}) (interface{}, error) {
		return ChangePollOrCry(w, r, session, func(poll *Poll) {
			var data AddOptionData
			mapstructure.Decode(protoData, &data)

			poll.Options = append(poll.Options, data.Value)
		})
	})
}

//RemoveOption ...
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

//Publish ...
func Publish(w http.ResponseWriter, r *http.Request) {
	ExecuteAuthenticated(w, r, func(session Session, protoData interface{}) (interface{}, error) {
		return ChangePollOrCry(w, r, session, func(poll *Poll) {
			poll.Published = true
		})
	})
}

//ChangePollOrCry ...
func ChangePollOrCry(w http.ResponseWriter, r *http.Request, session Session, fChange func(aPoll *Poll)) (Poll, error) {
	vars := mux.Vars(r)
	id, err := kallax.NewULIDFromText(vars["id"])
	if err != nil {
		return Poll{}, ErrNotChangePoll(err.Error())
	}

	poll := FindPollById(id)

	if poll.Published {
		return Poll{}, ErrNotChangePoll("Can't change a published poll.")
	}

	if poll.Owner != session.UserID {
		return Poll{}, ErrNotChangePoll("Can't change a poll from other user.")
	}

	fChange(&poll)

	return UpdatePoll(poll), nil
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

	if !ExistsOption(pollID, data.Value) {
		message := fmt.Sprintf("There is no option %s for vote on this poll.", data.Value)
		http.Error(w, message, http.StatusConflict)
		return
	}

	if PollAlreadyVotedByUser(pollID, session.UserID) {
		message := fmt.Sprintf("You already voted in this poll.")
		http.Error(w, message, http.StatusConflict)
		return
	}

	vote := PollVote{
		ID:           kallax.NewULID(),
		PoolID:       pollID,
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
	poll := FindPollById(pollID)
	pollMap := votes[pollID]

	total := TotalVotes(pollID)
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

//TotalVotes ...
func TotalVotes(pollID kallax.ULID) int {
	sum := 0

	for _, list := range votes[pollID] {
		sum += len(list)
	}

	return sum
}

//GetPoll ...
func GetPoll(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session Session) interface{} {
		vars := mux.Vars(r)
		id, err := kallax.NewULIDFromText(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return Poll{}
		}

		return FindPollById(id)
	})
}

//CountingPollVotes ...
func CountingPollVotes(w http.ResponseWriter, r *http.Request) {
	ExecuteSessioned(w, r, func(session Session) interface{} {
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
	ExecuteSessioned(w, r, func(session Session) interface{} {
		log.Println("Getting all")
		pollsSorted := make([]Poll, 0)
		for _, poll := range polls {
			pollsSorted = append(pollsSorted, poll)
		}

		return SortPolls(pollsSorted)
	})
}

//GetPollsMine ...
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

//SortPolls ...
func SortPolls(polls []Poll) []Poll {
	sort.Slice(polls, func(i, j int) bool {
		return polls[i].CreatedAt.Before(polls[j].CreatedAt)
	})

	return polls
}

//ExecuteSessioned ...
func ExecuteSessioned(w http.ResponseWriter, r *http.Request, f func(session Session) interface{}) {
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

//CheckSession ...
func CheckSession(w http.ResponseWriter, r *http.Request) (Session, error) {
	sessionID := r.Header.Get("sessionId")

	if sessionID == "" {
		return Session{}, ErrUserNotLogged("Must be logged to perform this action. Missing value.")
	}

	ID, err := kallax.NewULIDFromText(sessionID)
	if err != nil {
		return Session{}, err
	}

	session, ok := FindSessionById(ID)

	if !ok {
		return Session{}, ErrUserNotLogged("Must be logged to perform this action. Session invalid.")
	}

	return session, nil
}

//ConnectToDatabase ...
func ConnectToDatabase() {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=poll password=poll dbname=poll sslmode=disable")

	if err != nil {
		panic(err)
	}

	userStore = NewUserStore(db)

	log.Println("Successfuly connected!")
}

func ConfigStartServer() {
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

//main ...
func main() {
	ConnectToDatabase()
	ConfigStartServer()
}
