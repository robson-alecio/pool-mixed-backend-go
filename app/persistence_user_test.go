package app

import (
	"testing"

	"gopkg.in/src-d/go-kallax.v1"

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
	userStoreMock := &IUserStoreMock{
		SaveFunc: func(record *User) (bool, error) {
			return true, nil
		},
	}
	handler := UserHandlerImpl{
		Store: userStoreMock,
	}

	user := User{}
	savedUser := handler.SaveUser(user)

	assert.AssertEqual(t, user, savedUser)
	assert.AssertEqual(t, 1, userStoreMock.SaveCalls())
}

func TestFindUserByLogin(t *testing.T) {
	var sqlExecuted string

	userStoreMock := &IUserStoreMock{
		FindOneFunc: func(q *UserQuery) (*User, error) {
			sqlExecuted = q.String()
			return &User{}, nil
		},
	}
	handler := UserHandlerImpl{
		Store: userStoreMock,
	}

	result, err := handler.FindUserByLogin("fulano@detal.com")

	assert.AssertNotNil(t, result)
	assert.AssertNil(t, err)
	sql := "SELECT __user.id, __user.created_at, __user.updated_at, __user.login, __user.name, __user.password " +
		"FROM poll_user __user WHERE __user.login = $1"
	assert.AssertEqual(t, sql, sqlExecuted)

}
func TestFindUserByID(t *testing.T) {
	var sqlExecuted string

	userStoreMock := &IUserStoreMock{
		FindOneFunc: func(q *UserQuery) (*User, error) {
			sqlExecuted = q.String()
			return &User{}, nil
		},
	}
	handler := UserHandlerImpl{
		Store: userStoreMock,
	}

	result, err := handler.FindUserByID(kallax.NewULID())

	assert.AssertNotNil(t, result)
	assert.AssertNil(t, err)
	sql := "SELECT __user.id, __user.created_at, __user.updated_at, __user.login, __user.name, __user.password " +
		"FROM poll_user __user WHERE __user.id IN ($1)"
	assert.AssertEqual(t, sql, sqlExecuted)
}

func TestCreateAnonUser(t *testing.T) {
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

	savedUser := handler.CreateAnonUser()

	assert.AssertTrue(t, saved)
	assert.AssertNotNil(t, savedUser)
	assert.AssertMatchString(t, "Anon[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}", savedUser.Login)
}
