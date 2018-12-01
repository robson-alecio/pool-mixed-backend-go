package app

import (
	"database/sql"
	"log"

	kallax "gopkg.in/src-d/go-kallax.v1"
)

//SessionHandler ...
//go:generate moq -out sessionhandler_moq.go . SessionHandler
type SessionHandler interface {
	CreateSession(userID kallax.ULID) *Session
	SaveSession(session Session) Session
}

//ISessionStore ...
//go:generate moq -out isessionstore_moq.go . ISessionStore
type ISessionStore interface {
	Save(record *Session) (updated bool, err error)
	FindOne(q *SessionQuery) (*Session, error)
}

//SessionHandlerImpl ...
type SessionHandlerImpl struct {
	Store ISessionStore
}

//NewSessionHandler ...
func NewSessionHandler(db *sql.DB) *SessionHandlerImpl {
	return &SessionHandlerImpl{
		Store: NewSessionStore(db),
	}
}

//CreateSession ...
func (h SessionHandlerImpl) CreateSession(userID kallax.ULID) *Session {
	session := Session{
		ID:     kallax.NewULID(),
		UserID: userID,
	}
	h.SaveSession(session)

	return &session
}

//SaveSession ...
func (h SessionHandlerImpl) SaveSession(session Session) Session {
	log.Println("Saving Session", session)

	h.Store.Save(&session)
	return session
}

// FindSessionByID ...
func (h SessionHandlerImpl) FindSessionByID(id kallax.ULID) (*Session, error) {
	query := NewSessionQuery().FindByID(id)
	return h.Store.FindOne(query)
}
