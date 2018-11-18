package main

import (
	"gopkg.in/src-d/go-kallax.v1"
)

//go:generate kallax gen -e main.go

//User ...
type User struct {
	kallax.Model
	kallax.Timestamps
	ID       kallax.ULID `pk:""`
	Login    string
	Name     string
	Password string
}

//Session ...
type Session struct {
	kallax.Model
	kallax.Timestamps
	ID     kallax.ULID `pk:""`
	UserID kallax.ULID
}

//Poll ...
type Poll struct {
	kallax.Model
	kallax.Timestamps
	ID        kallax.ULID `pk:""`
	Name      string
	Options   []string
	Owner     kallax.ULID
	Published bool
}

//PollVote ...
type PollVote struct {
	kallax.Model
	kallax.Timestamps
	ID           kallax.ULID `pk:""`
	PoolID       kallax.ULID
	UserID       kallax.ULID
	ChosenOption string
}
