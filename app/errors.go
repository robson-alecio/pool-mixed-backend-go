package app

import "fmt"

//ErrPasswordDoNotMatch ...
type ErrPasswordDoNotMatch string

func (e ErrPasswordDoNotMatch) Error() string {
	return string(e)
}

//ErrUserNotFound ...
type ErrUserNotFound struct {
	Login string
}

func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("User %s not found.", e.Login)
}

//ErrUserNotLogged ...
type ErrUserNotLogged string

func (e ErrUserNotLogged) Error() string {
	return string(e)
}

//ErrNotChangePoll ...
type ErrNotChangePoll string

func (e ErrNotChangePoll) Error() string {
	return string(e)
}
