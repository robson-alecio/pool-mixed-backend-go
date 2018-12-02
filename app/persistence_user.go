package app

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	kallax "gopkg.in/src-d/go-kallax.v1"
)

//UserHandler ...
//go:generate moq -out userhandler_moq.go . UserHandler
type UserHandler interface {
	CreateUserFromData(d *UserCreationData) (User, error)
	SaveUser(user User) User
	FindUserByLogin(login string) (*User, error)
	FindUserByID(ID kallax.ULID) (*User, error)
	CreateAnonUser() User
	FindUserByLoginAndPassword(login, password string) (*User, error)
}

//IUserStore ...
//go:generate moq -out iuserstore_moq.go . IUserStore
type IUserStore interface {
	Save(record *User) (updated bool, err error)
	FindOne(q *UserQuery) (*User, error)
}

//UserHandlerImpl ...
type UserHandlerImpl struct {
	Store IUserStore
}

//NewUserHandler ...
func NewUserHandler(db *sql.DB) *UserHandlerImpl {
	return &UserHandlerImpl{
		Store: NewUserStore(db),
	}
}

//CreateUserFromData ...
func (handler *UserHandlerImpl) CreateUserFromData(d *UserCreationData) (User, error) {
	if d.Password != d.PasswordConfirm {
		return User{}, ErrPasswordDoNotMatch("Passwords don't match")
	}

	hasher := md5.New()
	hasher.Write([]byte(d.Password))
	encryptedPassword := hex.EncodeToString(hasher.Sum(nil))

	user := User{
		ID:       kallax.NewULID(),
		Login:    d.Login,
		Name:     d.Name,
		Password: encryptedPassword,
	}
	return user, nil
}

//SaveUser ...
func (handler *UserHandlerImpl) SaveUser(user User) User {
	log.Println("Saving User", user)

	handler.Store.Save(&user)

	return user
}

//FindUserByLogin ...
func (handler *UserHandlerImpl) FindUserByLogin(login string) (*User, error) {
	query := NewUserQuery().FindByLogin(login)
	return handler.Store.FindOne(query)
}

//FindUserByID ...
func (handler *UserHandlerImpl) FindUserByID(ID kallax.ULID) (*User, error) {
	query := NewUserQuery().FindByID(ID)
	return handler.Store.FindOne(query)
}

//CreateAnonUser ...
func (handler *UserHandlerImpl) CreateAnonUser() User {
	user := User{
		ID: kallax.NewULID(),
	}
	user.Name = "Anon" + user.ID.String()
	user.Login = user.Name

	handler.SaveUser(user)

	return user
}

//FindUserByLoginAndPassword ...
func (handler *UserHandlerImpl) FindUserByLoginAndPassword(login, password string) (*User, error) {
	query := NewUserQuery().
		FindByLogin(login).
		FindByPassword(password)

	user, err := handler.Store.FindOne(query)

	if err != nil {
		return nil, fmt.Errorf("User and password invalid")
	}

	return user, nil
}
