package app

import (
	"fmt"
	"testing"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/chai2010/assert"
)

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
	t.Fail()
}

func TestShouldChangeCryWhenNotFindPoll(t *testing.T) {
	t.Fail()
}

func TestShouldChangeCryWhenPollPublished(t *testing.T) {
	t.Fail()
}

func TestShouldChangeCryWhenPollOwnedByOtherUser(t *testing.T) {
	t.Fail()
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
	assert.AssertEqual(t, 1, len(pollOptionHandlerMock.SavePollOptionCalls()))
}

func TestCreatePollOptionFromData(t *testing.T) {
	t.Fail()
}
