package app

import (
	"fmt"
	"testing"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/chai2010/assert"
)

type ProcessErrorBox struct {
	Object       interface{}
	ErrorOcurred error
}

func helperMockProcessFuncBoxed(box *ProcessErrorBox) func(interface{}, ...ProcessingBlock) {
	return func(v interface{}, blocks ...ProcessingBlock) {
		var r interface{} = v
		var e error
		box.Object = r
		for _, f := range blocks {
			r, e = f(r)

			if e != nil {
				box.ErrorOcurred = e
				return
			}

			box.Object = r
		}
	}
}

func helperMockProcessFuncInputed(input interface{}) func(interface{}, ...ProcessingBlock) {
	return func(v interface{}, blocks ...ProcessingBlock) {
		var r interface{} = input
		var e error
		for _, f := range blocks {
			r, e = f(r)

			if e != nil {
				return
			}
		}
	}
}

func helperMockProcessFuncBoxedInputed(box *ProcessErrorBox, input interface{}) func(interface{}, ...ProcessingBlock) {
	return func(v interface{}, blocks ...ProcessingBlock) {
		var r interface{} = input
		var e error
		box.Object = r
		for _, f := range blocks {
			r, e = f(r)

			if e != nil {
				box.ErrorOcurred = e
				return
			}

			box.Object = r
		}
	}
}

func helperMockProcessFunc(v interface{}, blocks ...ProcessingBlock) {
	var r interface{} = v
	var e error
	for _, f := range blocks {
		r, e = f(r)

		if e != nil {
			return
		}
	}
}

func loggedUserID() kallax.ULID {
	ulid, _ := kallax.NewULIDFromText("c4dfaa18-103f-47a9-b6d3-ece7758832b5")
	return ulid
}

func createBasicHelperMock() *HTTPHelperMock {
	return &HTTPHelperMock{
		ProcessFunc: helperMockProcessFunc,
	}
}

func createAuthenticatedHelperMock() *HTTPHelperMock {
	return &HTTPHelperMock{
		ProcessFunc:          helperMockProcessFunc,
		ValidateSessionFunc:  func() error { return nil },
		IsRegisteredUserFunc: func() bool { return true },
		LoggedUserIDFunc:     loggedUserID,
	}
}

func getPollIDVarValue(string) string {
	return "01678ef4-3fd6-7e86-a52b-a1ed224aa249"
}

func createPollChangeHelperMock() *HTTPHelperMock {
	helperMock := createAuthenticatedHelperMock()

	helperMock.GetVarFunc = getPollIDVarValue

	return helperMock
}

func createPollChangeProcessBoxedHelperMock(box *ProcessErrorBox) *HTTPHelperMock {
	helperMock := createPollChangeHelperMock()

	helperMock.ProcessFunc = helperMockProcessFuncBoxed(box)

	return helperMock
}

func TestCreateUser(t *testing.T) {
	helperMock := createBasicHelperMock()
	handlerMock := &UserHandlerMock{
		CreateUserFromDataFunc: func(d *UserCreationData) (User, error) {
			return User{}, nil
		},
		SaveUserFunc: func(user User) User {
			return user
		},
	}

	CreateUser(helperMock, handlerMock)

	assert.AssertEqual(t, 1, len(handlerMock.CreateUserFromDataCalls()))
	assert.AssertEqual(t, 1, len(handlerMock.SaveUserCalls()))
}

func TestVisit(t *testing.T) {
	helperMock := createBasicHelperMock()
	userHandlerMock := &UserHandlerMock{
		CreateAnonUserFunc: func() User {
			return User{
				ID: kallax.NewULID(),
			}
		},
	}

	sessionHandlerMock := &SessionHandlerMock{
		CreateSessionFunc: func(ID kallax.ULID, flag bool) *Session {
			return &Session{}
		},
	}

	Visit(helperMock, userHandlerMock, sessionHandlerMock)

	assert.AssertEqual(t, 1, len(userHandlerMock.CreateAnonUserCalls()))
	assert.AssertEqual(t, 1, len(sessionHandlerMock.CreateSessionCalls()))
}

func TestLogin(t *testing.T) {
	helperMock := createBasicHelperMock()
	userHandlerMock := &UserHandlerMock{
		FindUserByLoginAndPasswordFunc: func(login, password string) (*User, error) {
			return &User{
				ID: kallax.NewULID(),
			}, nil
		},
	}

	sessionHandlerMock := &SessionHandlerMock{
		CreateSessionFunc: func(ID kallax.ULID, flag bool) *Session {
			return &Session{}
		},
	}

	Login(helperMock, userHandlerMock, sessionHandlerMock)

	assert.AssertEqual(t, 1, len(userHandlerMock.FindUserByLoginAndPasswordCalls()))
	assert.AssertEqual(t, 1, len(sessionHandlerMock.CreateSessionCalls()))
}

func TestStartCreatePoll(t *testing.T) {
	var savedPoll Poll
	helperMock := createAuthenticatedHelperMock()
	pollHandlerMock := &PollHandlerMock{
		SavePollFunc: func(poll Poll) Poll {
			savedPoll = poll
			return poll
		},
	}

	StartCreatePoll(helperMock, pollHandlerMock)

	assert.AssertEqual(t, loggedUserID(), savedPoll.Owner)
	assert.AssertEqual(t, 1, len(pollHandlerMock.SavePollCalls()))
}

func TestShouldGetErrorForSessionInvalidOnCheckAuthentication(t *testing.T) {
	helperMock := &HTTPHelperMock{
		ValidateSessionFunc:  func() error { return fmt.Errorf("Dammit") },
		IsRegisteredUserFunc: func() bool { return true },
	}

	err := CheckAuthentication(helperMock)

	assert.AssertEqual(t, "Dammit", err.Error())
}

func TestShouldGetErrorForUserUnregisteredOnCheckAuthentication(t *testing.T) {
	helperMock := &HTTPHelperMock{
		ValidateSessionFunc:  func() error { return nil },
		IsRegisteredUserFunc: func() bool { return false },
	}

	err := CheckAuthentication(helperMock)

	assert.AssertEqual(t, "Must be logged to perform this action. Not authenticated.", err.Error())
}

func TestShouldForbidExecution(t *testing.T) {
	var errorMessage string
	helperMock := &HTTPHelperMock{
		ValidateSessionFunc: func() error { return fmt.Errorf("Invalid session") },
		ForbidFunc:          func(err error) { errorMessage = err.Error() },
	}

	ExecuteAuthenticated(helperMock, &LoginData{})

	assert.AssertEqual(t, 1, len(helperMock.ForbidCalls()))
	assert.AssertEqual(t, "Invalid session", errorMessage)
}

func TestShouldChangePollCryWhenNotExtractPollID(t *testing.T) {
	box := &ProcessErrorBox{}

	helperMock := createPollChangeHelperMock()
	helperMock.GetVarFunc = func(name string) string {
		return "avocado"
	}
	helperMock.ProcessFunc = helperMockProcessFuncBoxed(box)

	changePollOrCry(helperMock, &AddOptionData{}, nil, nil, nil)

	assert.AssertEqual(t, "uuid: UUID string too short: avocado", box.ErrorOcurred.Error())
}

func TestShouldChangeCryWhenNotFindPoll(t *testing.T) {
	box := &ProcessErrorBox{}
	helperMock := createPollChangeProcessBoxedHelperMock(box)

	pollHandlerMock := &PollHandlerMock{
		FindPollByIDFunc: func(ID kallax.ULID) (*Poll, error) {
			return nil, fmt.Errorf("Deadpoll")
		},
	}

	changePollOrCry(helperMock, &AddOptionData{}, pollHandlerMock, nil, nil)

	assert.AssertEqual(t, "Deadpoll", box.ErrorOcurred.Error())
}

func TestShouldChangeCryWhenPollPublished(t *testing.T) {
	box := &ProcessErrorBox{}
	helperMock := createPollChangeProcessBoxedHelperMock(box)

	pollHandlerMock := &PollHandlerMock{
		FindPollByIDFunc: func(ID kallax.ULID) (*Poll, error) {
			return &Poll{
				Published: true,
			}, nil
		},
	}

	changePollOrCry(helperMock, &AddOptionData{}, pollHandlerMock, nil, nil)

	assert.AssertEqual(t, "Can't change a published poll.", box.ErrorOcurred.Error())
}

func TestShouldChangeCryWhenPollOwnedByOtherUser(t *testing.T) {
	box := &ProcessErrorBox{}
	helperMock := createPollChangeProcessBoxedHelperMock(box)

	pollHandlerMock := &PollHandlerMock{
		FindPollByIDFunc: func(ID kallax.ULID) (*Poll, error) {
			ownerID, _ := kallax.NewULIDFromText("b5ffeb10-712a-45ee-b939-e27033bf1db5")
			return &Poll{
				Owner: ownerID,
			}, nil
		},
	}

	changePollOrCry(helperMock, &AddOptionData{}, pollHandlerMock, nil, nil)

	assert.AssertEqual(t, "Can't change a poll from other user.", box.ErrorOcurred.Error())
}

func TestAddOption(t *testing.T) {
	helperMock := createPollChangeHelperMock()

	pollHandlerMock := &PollHandlerMock{
		FindPollByIDFunc: func(ID kallax.ULID) (*Poll, error) {
			return &Poll{
				Published: false,
				Owner:     loggedUserID(),
			}, nil
		},
		SavePollFunc: func(v Poll) Poll {
			return v
		},
	}

	pollOptionHandlerMock := &PollOptionHandlerMock{
		SavePollOptionFunc: func(v PollOption) PollOption {
			return v
		},
	}

	AddOption(helperMock, pollHandlerMock, pollOptionHandlerMock)

	assert.AssertEqual(t, 1, len(helperMock.GetVarCalls()))
	assert.AssertEqual(t, 1, len(pollHandlerMock.FindPollByIDCalls()))
	assert.AssertEqual(t, 1, len(pollHandlerMock.SavePollCalls()))
	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.SavePollOptionCalls()))
}

func TestCreatePollOptionFromData(t *testing.T) {
	poll := &Poll{}

	option := createPollOptionFrom(poll, &AddOptionData{"Opt"})

	assert.AssertEqual(t, poll, option.Owner)
	assert.AssertEqual(t, "Opt", option.Content)
}

func TestRemoveOption(t *testing.T) {
	helperMock := createPollChangeHelperMock()
	helperMock.ProcessFunc = helperMockProcessFuncInputed(&RemoveOptionData{
		Value: "9d627cdc-8e4a-435e-a2f7-c9bafaa41e45",
	})

	pollHandlerMock := &PollHandlerMock{
		FindPollByIDFunc: func(ID kallax.ULID) (*Poll, error) {
			return &Poll{
				Published: false,
				Owner:     loggedUserID(),
			}, nil
		},
		SavePollFunc: func(v Poll) Poll {
			return v
		},
	}

	pollOptionHandlerMock := &PollOptionHandlerMock{
		DeletePollOptionFunc: func(id kallax.ULID) error {
			return nil
		},
	}

	RemoveOption(helperMock, pollHandlerMock, pollOptionHandlerMock)

	assert.AssertEqual(t, 1, len(helperMock.GetVarCalls()))
	assert.AssertEqual(t, 1, len(pollHandlerMock.FindPollByIDCalls()))
	assert.AssertEqual(t, 1, len(pollHandlerMock.SavePollCalls()))
	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.DeletePollOptionCalls()))
}

func TestRemoveOptionWhenIdDoesNotExists(t *testing.T) {
	helperMock := createPollChangeHelperMock()

	pollHandlerMock := &PollHandlerMock{
		FindPollByIDFunc: func(ID kallax.ULID) (*Poll, error) {
			return &Poll{
				Published: false,
				Owner:     loggedUserID(),
			}, nil
		},
		SavePollFunc: func(v Poll) Poll {
			return v
		},
	}

	pollOptionHandlerMock := &PollOptionHandlerMock{
		DeletePollOptionFunc: func(id kallax.ULID) error {
			return nil
		},
	}

	RemoveOption(helperMock, pollHandlerMock, pollOptionHandlerMock)

	assert.AssertEqual(t, 1, len(helperMock.GetVarCalls()))
	assert.AssertEqual(t, 1, len(pollHandlerMock.FindPollByIDCalls()))
	assert.AssertEqual(t, 0, len(pollHandlerMock.SavePollCalls()))
	assert.AssertEqual(t, 0, len(pollOptionHandlerMock.DeletePollOptionCalls()))
}

func TestPublish(t *testing.T) {
	helperMock := createPollChangeHelperMock()

	poll := &Poll{
		Published: false,
		Owner:     loggedUserID(),
	}

	pollHandlerMock := &PollHandlerMock{
		FindPollByIDFunc: func(ID kallax.ULID) (*Poll, error) {
			return poll, nil
		},
		SavePollFunc: func(v Poll) Poll {
			return v
		},
	}

	pollOptionHandlerMock := &PollOptionHandlerMock{}

	Publish(helperMock, pollHandlerMock, pollOptionHandlerMock)

	assert.AssertEqual(t, 1, len(helperMock.GetVarCalls()))
	assert.AssertEqual(t, 1, len(pollHandlerMock.FindPollByIDCalls()))
	assert.AssertEqual(t, 1, len(pollHandlerMock.SavePollCalls()))
	assert.AssertTrue(t, poll.Published)
}

func TestCreateVote(t *testing.T) {
	helperMock := createAuthenticatedHelperMock()
	helperMock.GetVarFunc = func(name string) string {
		return "c5c1827e-2649-49ee-b960-cd04ac34c1a8"
	}
	helperMock.ProcessFunc = helperMockProcessFuncInputed(&PollVoteData{
		Value: "Terceira",
	})

	pollOptionHandlerMock := &PollOptionHandlerMock{
		ExistsOptionFunc: func(pollID kallax.ULID, candidate string) (bool, error) {
			return true, nil
		},
		FindPollOptionsFunc: func(pollID kallax.ULID) ([]*PollOption, error) {
			return []*PollOption{
				&PollOption{Content: "A"},
				&PollOption{Content: "B"},
				&PollOption{Content: "C"},
			}, nil
		},
	}

	pollVoteHandlerMock := &PollVoteHandlerMock{
		PollAlreadyVotedByUserFunc: func(pollID kallax.ULID, userID kallax.ULID) (bool, error) {
			return false, nil
		},
		SaveVoteFunc: func(vote PollVote) PollVote {
			return vote
		},
		VotesForFunc: func(pollID kallax.ULID, option string) int64 {
			return 1
		},
	}

	CreateVote(helperMock, pollOptionHandlerMock, pollVoteHandlerMock)

	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.ExistsOptionCalls()))
	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.FindPollOptionsCalls()))
	assert.AssertEqual(t, 1, len(pollVoteHandlerMock.PollAlreadyVotedByUserCalls()))
	assert.AssertEqual(t, 1, len(pollVoteHandlerMock.SaveVoteCalls()))
	assert.AssertEqual(t, 3, len(pollVoteHandlerMock.VotesForCalls()))
}

func TestShouldNotCreateVoteWithoutPollID(t *testing.T) {
	box := &ProcessErrorBox{}

	helperMock := createAuthenticatedHelperMock()
	helperMock.GetVarFunc = func(name string) string {
		return "no-uuid"
	}
	helperMock.ProcessFunc = helperMockProcessFuncBoxed(box)
	// helperMock.ProcessFunc = helperMockProcessFuncInputed(&PollVoteData{
	// 	Value: "Terceira",
	// })

	pollOptionHandlerMock := &PollOptionHandlerMock{
		// ExistsOptionFunc: func(pollID kallax.ULID, candidate string) (bool, error) {
		// 	return true, nil
		// },
		// FindPollOptionsFunc: func(pollID kallax.ULID) ([]*PollOption, error) {
		// 	return []*PollOption{
		// 		&PollOption{Content: "A"},
		// 		&PollOption{Content: "B"},
		// 		&PollOption{Content: "C"},
		// 	}, nil
		// },
	}

	pollVoteHandlerMock := &PollVoteHandlerMock{
		// PollAlreadyVotedByUserFunc: func(pollID kallax.ULID, userID kallax.ULID) (bool, error) {
		// 	return false, nil
		// },
		// SaveVoteFunc: func(vote PollVote) PollVote {
		// 	return vote
		// },
		// VotesForFunc: func(pollID kallax.ULID, option string) int64 {
		// 	return 1
		// },
	}

	CreateVote(helperMock, pollOptionHandlerMock, pollVoteHandlerMock)

	assert.AssertEqual(t, 0, len(pollOptionHandlerMock.ExistsOptionCalls()))
	assert.AssertEqual(t, 0, len(pollOptionHandlerMock.FindPollOptionsCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.PollAlreadyVotedByUserCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.SaveVoteCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.VotesForCalls()))

	assert.AssertEqual(t, "uuid: UUID string too short: no-uuid", box.ErrorOcurred.Error())
}

func TestShouldCreateVoteFailWhenOptionExistsFail(t *testing.T) {
	box := &ProcessErrorBox{}

	helperMock := createAuthenticatedHelperMock()
	helperMock.GetVarFunc = func(name string) string {
		return "c5c1827e-2649-49ee-b960-cd04ac34c1a8"
	}
	helperMock.ProcessFunc = helperMockProcessFuncBoxedInputed(box, &PollVoteData{
		Value: "Terceira",
	})

	pollOptionHandlerMock := &PollOptionHandlerMock{
		ExistsOptionFunc: func(pollID kallax.ULID, candidate string) (bool, error) {
			return false, fmt.Errorf("Fail")
		},
		// FindPollOptionsFunc: func(pollID kallax.ULID) ([]*PollOption, error) {
		// 	return []*PollOption{
		// 		&PollOption{Content: "A"},
		// 		&PollOption{Content: "B"},
		// 		&PollOption{Content: "C"},
		// 	}, nil
		// },
	}

	pollVoteHandlerMock := &PollVoteHandlerMock{
		// PollAlreadyVotedByUserFunc: func(pollID kallax.ULID, userID kallax.ULID) (bool, error) {
		// 	return false, nil
		// },
		// SaveVoteFunc: func(vote PollVote) PollVote {
		// 	return vote
		// },
		// VotesForFunc: func(pollID kallax.ULID, option string) int64 {
		// 	return 1
		// },
	}

	CreateVote(helperMock, pollOptionHandlerMock, pollVoteHandlerMock)

	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.ExistsOptionCalls()))
	assert.AssertEqual(t, 0, len(pollOptionHandlerMock.FindPollOptionsCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.PollAlreadyVotedByUserCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.SaveVoteCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.VotesForCalls()))

	assert.AssertEqual(t, "Fail", box.ErrorOcurred.Error())
}

func TestShouldCreateVoteFailWhenOptionNotExists(t *testing.T) {
	box := &ProcessErrorBox{}

	helperMock := createAuthenticatedHelperMock()
	helperMock.GetVarFunc = func(name string) string {
		return "c5c1827e-2649-49ee-b960-cd04ac34c1a8"
	}
	helperMock.ProcessFunc = helperMockProcessFuncBoxedInputed(box, &PollVoteData{
		Value: "Terceira",
	})

	pollOptionHandlerMock := &PollOptionHandlerMock{
		ExistsOptionFunc: func(pollID kallax.ULID, candidate string) (bool, error) {
			return false, nil
		},
		// FindPollOptionsFunc: func(pollID kallax.ULID) ([]*PollOption, error) {
		// 	return []*PollOption{
		// 		&PollOption{Content: "A"},
		// 		&PollOption{Content: "B"},
		// 		&PollOption{Content: "C"},
		// 	}, nil
		// },
	}

	pollVoteHandlerMock := &PollVoteHandlerMock{
		// PollAlreadyVotedByUserFunc: func(pollID kallax.ULID, userID kallax.ULID) (bool, error) {
		// 	return false, nil
		// },
		// SaveVoteFunc: func(vote PollVote) PollVote {
		// 	return vote
		// },
		// VotesForFunc: func(pollID kallax.ULID, option string) int64 {
		// 	return 1
		// },
	}

	CreateVote(helperMock, pollOptionHandlerMock, pollVoteHandlerMock)

	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.ExistsOptionCalls()))
	assert.AssertEqual(t, 0, len(pollOptionHandlerMock.FindPollOptionsCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.PollAlreadyVotedByUserCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.SaveVoteCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.VotesForCalls()))

	assert.AssertEqual(t, "There is no option Terceira for vote on this poll", box.ErrorOcurred.Error())
}

func TestShouldCreateVoteFailWhenVerifyVoteExistsFail(t *testing.T) {
	box := &ProcessErrorBox{}

	helperMock := createAuthenticatedHelperMock()
	helperMock.GetVarFunc = func(name string) string {
		return "c5c1827e-2649-49ee-b960-cd04ac34c1a8"
	}
	helperMock.ProcessFunc = helperMockProcessFuncBoxedInputed(box, &PollVoteData{
		Value: "Terceira",
	})

	pollOptionHandlerMock := &PollOptionHandlerMock{
		ExistsOptionFunc: func(pollID kallax.ULID, candidate string) (bool, error) {
			return true, nil
		},
		// FindPollOptionsFunc: func(pollID kallax.ULID) ([]*PollOption, error) {
		// 	return []*PollOption{
		// 		&PollOption{Content: "A"},
		// 		&PollOption{Content: "B"},
		// 		&PollOption{Content: "C"},
		// 	}, nil
		// },
	}

	pollVoteHandlerMock := &PollVoteHandlerMock{
		PollAlreadyVotedByUserFunc: func(pollID kallax.ULID, userID kallax.ULID) (bool, error) {
			return false, fmt.Errorf("Fail")
		},
		// SaveVoteFunc: func(vote PollVote) PollVote {
		// 	return vote
		// },
		// VotesForFunc: func(pollID kallax.ULID, option string) int64 {
		// 	return 1
		// },
	}

	CreateVote(helperMock, pollOptionHandlerMock, pollVoteHandlerMock)

	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.ExistsOptionCalls()))
	assert.AssertEqual(t, 0, len(pollOptionHandlerMock.FindPollOptionsCalls()))
	assert.AssertEqual(t, 1, len(pollVoteHandlerMock.PollAlreadyVotedByUserCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.SaveVoteCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.VotesForCalls()))

	assert.AssertEqual(t, "Fail", box.ErrorOcurred.Error())
}

func TestShouldCreateVoteFailWhenAlreadyVoted(t *testing.T) {
	box := &ProcessErrorBox{}

	helperMock := createAuthenticatedHelperMock()
	helperMock.GetVarFunc = func(name string) string {
		return "c5c1827e-2649-49ee-b960-cd04ac34c1a8"
	}
	helperMock.ProcessFunc = helperMockProcessFuncBoxedInputed(box, &PollVoteData{
		Value: "Terceira",
	})

	pollOptionHandlerMock := &PollOptionHandlerMock{
		ExistsOptionFunc: func(pollID kallax.ULID, candidate string) (bool, error) {
			return true, nil
		},
		// FindPollOptionsFunc: func(pollID kallax.ULID) ([]*PollOption, error) {
		// 	return []*PollOption{
		// 		&PollOption{Content: "A"},
		// 		&PollOption{Content: "B"},
		// 		&PollOption{Content: "C"},
		// 	}, nil
		// },
	}

	pollVoteHandlerMock := &PollVoteHandlerMock{
		PollAlreadyVotedByUserFunc: func(pollID kallax.ULID, userID kallax.ULID) (bool, error) {
			return true, nil
		},
		// SaveVoteFunc: func(vote PollVote) PollVote {
		// 	return vote
		// },
		// VotesForFunc: func(pollID kallax.ULID, option string) int64 {
		// 	return 1
		// },
	}

	CreateVote(helperMock, pollOptionHandlerMock, pollVoteHandlerMock)

	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.ExistsOptionCalls()))
	assert.AssertEqual(t, 0, len(pollOptionHandlerMock.FindPollOptionsCalls()))
	assert.AssertEqual(t, 1, len(pollVoteHandlerMock.PollAlreadyVotedByUserCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.SaveVoteCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.VotesForCalls()))

	assert.AssertEqual(t, "You already voted in this poll", box.ErrorOcurred.Error())
}

func TestShouldCountVotes(t *testing.T) {
	pollOptionHandlerMock := &PollOptionHandlerMock{
		FindPollOptionsFunc: func(pollID kallax.ULID) ([]*PollOption, error) {
			return []*PollOption{
				&PollOption{Content: "A"},
				&PollOption{Content: "B"},
				&PollOption{Content: "C"},
			}, nil
		},
	}

	pollVoteHandlerMock := &PollVoteHandlerMock{
		VotesForFunc: func(pollID kallax.ULID, option string) int64 {
			return 1
		},
	}

	votes := CountVotes(kallax.NewULID(), pollOptionHandlerMock, pollVoteHandlerMock)

	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.FindPollOptionsCalls()))
	assert.AssertEqual(t, 3, len(pollVoteHandlerMock.VotesForCalls()))

	assert.AssertEqual(t, 33.33, votes["A"])
	assert.AssertEqual(t, 33.33, votes["B"])
	assert.AssertEqual(t, 33.34, votes["C"])
	assert.AssertEqual(t, 3, votes["total"])
}

func TestShouldCountVotesWithoutRounding(t *testing.T) {
	pollOptionHandlerMock := &PollOptionHandlerMock{
		FindPollOptionsFunc: func(pollID kallax.ULID) ([]*PollOption, error) {
			return []*PollOption{
				&PollOption{Content: "A"},
				&PollOption{Content: "B"},
				&PollOption{Content: "C"},
			}, nil
		},
	}

	counts := map[string]int64{
		"A": 1,
		"B": 2,
		"C": 1,
	}
	pollVoteHandlerMock := &PollVoteHandlerMock{
		VotesForFunc: func(pollID kallax.ULID, option string) int64 {
			return counts[option]
		},
	}

	votes := CountVotes(kallax.NewULID(), pollOptionHandlerMock, pollVoteHandlerMock)

	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.FindPollOptionsCalls()))
	assert.AssertEqual(t, 3, len(pollVoteHandlerMock.VotesForCalls()))

	assert.AssertEqual(t, 25.0, votes["A"])
	assert.AssertEqual(t, 50.0, votes["B"])
	assert.AssertEqual(t, 25.0, votes["C"])
	assert.AssertEqual(t, 4, votes["total"])
}

func TestShouldCountVotesFailWhenFindExistentOptionsFail(t *testing.T) {
	errorMsg := "Unable to take values."
	pollOptionHandlerMock := &PollOptionHandlerMock{
		FindPollOptionsFunc: func(pollID kallax.ULID) ([]*PollOption, error) {
			return nil, fmt.Errorf(errorMsg)
		},
	}

	pollVoteHandlerMock := &PollVoteHandlerMock{}

	votes := CountVotes(kallax.NewULID(), pollOptionHandlerMock, pollVoteHandlerMock)

	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.FindPollOptionsCalls()))
	assert.AssertEqual(t, 0, len(pollVoteHandlerMock.VotesForCalls()))

	count, existsKey := votes[errorMsg]
	assert.AssertEqual(t, true, existsKey)
	assert.AssertEqual(t, -1, count)
	assert.AssertEqual(t, 1, len(votes))
}
