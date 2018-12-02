package app

import (
	"testing"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/chai2010/assert"
)

func HelperMockProcessFunc(v interface{}, blocks ...ProcessingBlock) {
	var r interface{} = v
	var e error
	for _, f := range blocks {
		r, e = f(r)

		if e != nil {
			return
		}
	}
}

func TestCreateUser(t *testing.T) {
	helperMock := &HTTPHelperMock{
		ProcessFunc: HelperMockProcessFunc,
	}
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
	helperMock := &HTTPHelperMock{
		ProcessFunc: HelperMockProcessFunc,
	}
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
	helperMock := &HTTPHelperMock{
		ProcessFunc: HelperMockProcessFunc,
	}
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
