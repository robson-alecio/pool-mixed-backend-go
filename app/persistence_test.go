package app

import (
	"testing"

	"github.com/chai2010/assert"
)

func TestCreateUserFromData(t *testing.T) {
	handler := UserHandlerImpl{}

	data := &UserCreationData{
		Login:           "phineas@disney.com",
		Name:            "Phineas Flynn",
		Password:        "summer",
		PasswordConfirm: "summer",
	}

	user, err := handler.CreateUserFromData(data)

	assert.AssertNil(t, err)
	assert.AssertEqual(t, "phineas@disney.com", user.Login)
	assert.AssertEqual(t, "Phineas Flynn", user.Name)
	assert.AssertEqual(t, "summer", user.Password)
}

func TestShouldCreateUserWhenPasswordNotConfirmed(t *testing.T) {
	handler := UserHandlerImpl{}

	data := &UserCreationData{
		Login:           "phineas@disney.com",
		Name:            "Phineas Flynn",
		Password:        "summer",
		PasswordConfirm: "winter",
	}

	user, err := handler.CreateUserFromData(data)

	assert.AssertNotNil(t, err)
	assert.AssertEqual(t, "Passwords don't match", err.Error())
	assert.AssertEqual(t, "", user.Login)
	assert.AssertEqual(t, "", user.Name)
	assert.AssertEqual(t, "", user.Password)
}

func TestSave(t *testing.T) {
	var saved bool

	userStoreMock := &IUserStoreMock{
		SaveFunc: func(record *User) (bool, error) {
			saved = true
			return true, nil
		},
	}
	handler := UserHandlerImpl{
		Store: userStoreMock,
	}

	user := User{}
	savedUser := handler.SaveUser(user)

	assert.AssertEqual(t, user, savedUser)
	assert.AssertTrue(t, saved)
}
