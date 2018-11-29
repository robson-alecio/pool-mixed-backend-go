package app

import (
	"fmt"
	"testing"
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
	var userCreated bool
	var userSaved bool

	helperMock := &HTTPHelperMock{
		ProcessFunc: HelperMockProcessFunc,
	}
	handlerMock := &UserHandlerMock{
		CreateUserFromDataFunc: func(d *UserCreationData) (User, error) {
			userCreated = true
			return User{}, nil
		},
		SaveUserFunc: func(user User) User {
			userSaved = true
			return user
		},
	}

	CreateUser(helperMock, handlerMock)

	if !userCreated {
		t.Error("User was not created.")
	}
	if !userSaved {
		t.Error("User was not saved.")
	}
}

func TestCreateButDontSaveUser(t *testing.T) {
	var userCreated bool
	var userSaved bool

	helperMock := &HTTPHelperMock{
		ProcessFunc: HelperMockProcessFunc,
	}
	handlerMock := &UserHandlerMock{
		CreateUserFromDataFunc: func(d *UserCreationData) (User, error) {
			userCreated = true
			return User{}, fmt.Errorf("I don't wanna save")
		},
		SaveUserFunc: func(user User) User {
			userSaved = true
			return user
		},
	}

	CreateUser(helperMock, handlerMock)

	if !userCreated {
		t.Error("User was not created.")
	}
	if userSaved {
		t.Error("User SHOULD NOT be saved.")
	}
}
