package app

import (
	"fmt"
	"math"

	kallax "gopkg.in/src-d/go-kallax.v1"
)

//CreateUser ...
func CreateUser(helper HTTPHelper, handler UserHandler) {
	createUser := func(v interface{}) (interface{}, error) {
		return handler.CreateUserFromData(v.(*UserCreationData))
	}

	saveUser := func(v interface{}) (interface{}, error) {
		return handler.SaveUser(v.(User)), nil
	}

	helper.Process(&UserCreationData{}, createUser, saveUser)
}

//Visit ...
func Visit(helper HTTPHelper, userHandler UserHandler, sessionHandler SessionHandler) {
	createAnonUser := func(v interface{}) (interface{}, error) {
		return userHandler.CreateAnonUser(), nil
	}

	createSession := func(v interface{}) (interface{}, error) {
		user := v.(User)
		return sessionHandler.CreateSession(user.ID, user.IsRegistered()), nil
	}

	helper.Process(nil, createAnonUser, createSession)
}

//Login ...
func Login(helper HTTPHelper, userHandler UserHandler, sessionHandler SessionHandler) {
	findUser := func(v interface{}) (interface{}, error) {
		loginData := v.(*LoginData)
		return userHandler.FindUserByLoginAndPassword(loginData.Login, loginData.Password)
	}

	createSession := func(v interface{}) (interface{}, error) {
		user := v.(*User)
		return sessionHandler.CreateSession(user.ID, user.IsRegistered()), nil
	}

	helper.Process(&LoginData{}, findUser, createSession)
}

//StartCreatePoll ...
func StartCreatePoll(helper HTTPHelper, pollHandler PollHandler) {
	createPoll := func(v interface{}) (interface{}, error) {
		data := v.(*CreatePollData)
		return pollHandler.SavePoll(Poll{
			ID:      kallax.NewULID(),
			Name:    data.Name,
			Options: make([]*PollOption, 0),
			Owner:   helper.LoggedUserID(),
		}), nil
	}
	ExecuteAuthenticated(helper, &CreatePollData{}, createPoll)
}

//ChangePollDataPack ...
type ChangePollDataPack struct {
	PollID     kallax.ULID
	PollTarget *Poll
	Data       interface{}
}

//AddOption ...
func AddOption(helper HTTPHelper, pollHandler PollHandler, pollOptionHandler PollOptionHandler) {
	effectiveChange := func(v interface{}) (interface{}, error) {
		pack := v.(*ChangePollDataPack)

		pollOption := createPollOptionFrom(pack.PollTarget, pack.Data.(*AddOptionData))
		pollOptionHandler.SavePollOption(*pollOption)

		return pack.PollTarget, nil
	}

	changePollOrCry(helper, &AddOptionData{}, pollHandler, pollOptionHandler, effectiveChange)
}

func createPollOptionFrom(poll *Poll, data *AddOptionData) *PollOption {
	return &PollOption{
		ID:      kallax.NewULID(),
		Owner:   poll,
		Content: data.Value,
	}
}

//RemoveOption ...
func RemoveOption(helper HTTPHelper, pollHandler PollHandler, pollOptionHandler PollOptionHandler) {
	effectiveChange := func(v interface{}) (interface{}, error) {
		pack := v.(*ChangePollDataPack)
		data := pack.Data.(*RemoveOptionData)

		id, err := kallax.NewULIDFromText(data.Value)

		if err != nil {
			return nil, err
		}

		pollOptionHandler.DeletePollOption(id)

		return pack.PollTarget, nil
	}

	changePollOrCry(helper, &RemoveOptionData{}, pollHandler, pollOptionHandler, effectiveChange)
}

//Publish ...
func Publish(helper HTTPHelper, pollHandler PollHandler, pollOptionHandler PollOptionHandler) {
	effectiveChange := func(v interface{}) (interface{}, error) {
		pack := v.(*ChangePollDataPack)

		pack.PollTarget.Published = true

		return pack.PollTarget, nil
	}

	changePollOrCry(helper, new(interface{}), pollHandler, pollOptionHandler, effectiveChange)
}

func changePollOrCry(helper HTTPHelper, data interface{}, pollHandler PollHandler,
	pollOptionHandler PollOptionHandler, effectiveChange ProcessingBlock) {
	getPollID := func(v interface{}) (interface{}, error) {
		ID, err := kallax.NewULIDFromText(helper.GetVar("id"))
		if err != nil {
			return nil, ErrNotChangePoll(err.Error())
		}

		return &ChangePollDataPack{
			PollID: ID,
			Data:   v,
		}, nil
	}

	getPoll := func(v interface{}) (interface{}, error) {
		pack := v.(*ChangePollDataPack)

		poll, errFind := pollHandler.FindPollByID(pack.PollID)

		if errFind != nil {
			return nil, errFind
		}

		pack.PollTarget = poll

		return pack, nil
	}

	checkPublished := func(v interface{}) (interface{}, error) {
		pack := v.(*ChangePollDataPack)

		if pack.PollTarget.Published {
			return nil, ErrNotChangePoll("Can't change a published poll.")
		}

		return pack, nil
	}

	checkOwner := func(v interface{}) (interface{}, error) {
		pack := v.(*ChangePollDataPack)

		if pack.PollTarget.Owner != helper.LoggedUserID() {
			return nil, ErrNotChangePoll("Can't change a poll from other user.")
		}

		return pack, nil
	}

	savePoll := func(v interface{}) (interface{}, error) {
		poll := v.(*Poll)

		pollHandler.SavePoll(*poll)

		return poll, nil
	}

	ExecuteAuthenticated(helper, data,
		getPollID, getPoll, checkPublished, checkOwner, effectiveChange, savePoll)
}

//CreateVoteDataPack ...
type CreateVoteDataPack struct {
	PollID      kallax.ULID
	Data        *PollVoteData
	VoteCreated *PollVote
}

//CreateVote ...
func CreateVote(helper HTTPHelper, pollOptionHandler PollOptionHandler, pollVoteHandler PollVoteHandler) {
	makeCreateVoteDataPack := func(v interface{}) (interface{}, error) {
		IDValue := helper.GetVar("id")
		pollID, err := kallax.NewULIDFromText(IDValue)
		if err != nil {
			return nil, err
		}

		return &CreateVoteDataPack{
			PollID: pollID,
			Data:   v.(*PollVoteData),
		}, nil
	}

	validateOption := func(v interface{}) (interface{}, error) {
		pack := v.(*CreateVoteDataPack)
		exists, err := pollOptionHandler.ExistsOption(pack.PollID, pack.Data.Value)

		if err != nil {
			return nil, err
		}

		if !exists {
			err := fmt.Errorf("There is no option %s for vote on this poll", pack.Data.Value)
			return nil, err
		}

		return v, nil
	}

	validateVoted := func(v interface{}) (interface{}, error) {
		pack := v.(*CreateVoteDataPack)

		voted, errVoted := pollVoteHandler.PollAlreadyVotedByUser(pack.PollID, helper.LoggedUserID())

		if errVoted != nil {
			return nil, errVoted
		}

		if voted {
			err := fmt.Errorf("You already voted in this poll")
			return nil, err
		}

		return v, nil
	}

	createVote := func(v interface{}) (interface{}, error) {
		pack := v.(*CreateVoteDataPack)

		pack.VoteCreated = &PollVote{
			ID:           kallax.NewULID(),
			PollID:       pack.PollID,
			UserID:       helper.LoggedUserID(),
			ChosenOption: pack.Data.Value,
		}

		pollVoteHandler.SaveVote(*(pack.VoteCreated))

		return pack, nil
	}

	mountResult := func(v interface{}) (interface{}, error) {
		pack := v.(*CreateVoteDataPack)

		result := PollVoteResult{
			VoteID:       pack.VoteCreated.ID.String(),
			VoteCounting: CountVotes(pack.PollID, pollOptionHandler, pollVoteHandler),
		}

		return result, nil
	}

	ExecuteSessioned(helper, &PollVoteData{}, makeCreateVoteDataPack, validateOption, validateVoted, createVote,
		mountResult)
}

//CountVotes ...
func CountVotes(pollID kallax.ULID, pollOptionHandler PollOptionHandler, pollVoteHandler PollVoteHandler) map[string]float64 {
	options, err := pollOptionHandler.FindPollOptions(pollID)

	if err != nil {
		return map[string]float64{
			err.Error(): -1.0,
		}
	}

	count := make(map[string]int64)

	for _, opt := range options {
		count[opt.Content] = pollVoteHandler.VotesFor(pollID, opt.Content)
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

//ExecuteSessioned ...
func ExecuteSessioned(helper HTTPHelper, v interface{}, blocks ...ProcessingBlock) {
	errCheck := helper.ValidateSession()

	if errCheck != nil {
		helper.Forbid(errCheck)
		return
	}

	helper.Process(v, blocks...)
}

//ExecuteAuthenticated ...
func ExecuteAuthenticated(helper HTTPHelper, v interface{}, blocks ...ProcessingBlock) {
	err := CheckAuthentication(helper)

	if err != nil {
		helper.Forbid(err)
		return
	}

	helper.Process(v, blocks...)
}

//CheckAuthentication ...
func CheckAuthentication(helper HTTPHelper) error {
	errCheck := helper.ValidateSession()

	if errCheck != nil {
		return errCheck
	}

	if !helper.IsRegisteredUser() {
		return ErrUserNotLogged("Must be logged to perform this action. Not authenticated.")
	}

	return nil
}
