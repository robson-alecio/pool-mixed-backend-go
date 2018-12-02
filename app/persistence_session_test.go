package app

import (
	"testing"

	"github.com/chai2010/assert"

	"gopkg.in/src-d/go-kallax.v1"
)

func TestCreateSession(t *testing.T) {
	store := &ISessionStoreMock{
		SaveFunc: func(session *Session) (bool, error) {
			return true, nil
		},
	}
	handler := SessionHandlerImpl{
		Store: store,
	}

	userID := kallax.NewULID()
	session := handler.CreateSession(userID)

	assert.AssertEqual(t, userID, session.UserID)
	assert.AssertEqual(t, 1, len(store.SaveCalls()))
}

func TestFindSessionByID(t *testing.T) {
	var sqlExecuted string

	store := &ISessionStoreMock{
		FindOneFunc: func(q *SessionQuery) (*Session, error) {
			sqlExecuted = q.String()
			return &Session{}, nil
		},
	}

	handler := SessionHandlerImpl{
		Store: store,
	}

	session, err := handler.FindSessionByID(kallax.NewULID())

	assert.AssertNotNil(t, session)
	assert.AssertNil(t, err)
	assert.AssertEqual(t, 1, len(store.FindOneCalls()))
	sqlExpected := "SELECT __session.id, __session.created_at, __session.updated_at, __session.user_id " +
		"FROM poll_session __session WHERE __session.id IN ($1)"
	assert.AssertEqual(t, sqlExpected, sqlExecuted)
}
