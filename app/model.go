package app

import (
	"gopkg.in/src-d/go-kallax.v1"
)

//go:generate kallax gen -e main.go -e persistence_user.go -e persistence_session.go -e bis.go

//User ...
type User struct {
	kallax.Model `table:"poll_user"`
	kallax.Timestamps
	ID       kallax.ULID `pk:""`
	Login    string
	Name     string
	Password string
}

//IsRegistered ...
func (u *User) IsRegistered() bool {
	return u.Password != ""
}

//Session ...
type Session struct {
	kallax.Model `table:"poll_session"`
	kallax.Timestamps
	ID             kallax.ULID `pk:""`
	UserID         kallax.ULID
	RegisteredUser bool
}

//Poll ...
type Poll struct {
	kallax.Model
	kallax.Timestamps
	ID        kallax.ULID `pk:""`
	Name      string
	Options   []*PollOption
	Owner     kallax.ULID
	Published bool
}

// PollOption ...
type PollOption struct {
	kallax.Model
	ID      kallax.ULID `pk:""`
	Owner   *Poll       `fk:"poll_id,inverse"`
	Content string
}

//PollVote ...
type PollVote struct {
	kallax.Model
	kallax.Timestamps
	ID           kallax.ULID `pk:""`
	PollID       kallax.ULID
	UserID       kallax.ULID
	ChosenOption string
}
