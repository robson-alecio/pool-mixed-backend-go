// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package app

import (
	"gopkg.in/src-d/go-kallax.v1"
	"sync"
)

var (
	lockSessionHandlerMockCreateSession sync.RWMutex
	lockSessionHandlerMockSaveSession   sync.RWMutex
)

// SessionHandlerMock is a mock implementation of SessionHandler.
//
//     func TestSomethingThatUsesSessionHandler(t *testing.T) {
//
//         // make and configure a mocked SessionHandler
//         mockedSessionHandler := &SessionHandlerMock{
//             CreateSessionFunc: func(userID kallax.ULID) *Session {
// 	               panic("mock out the CreateSession method")
//             },
//             SaveSessionFunc: func(session Session) Session {
// 	               panic("mock out the SaveSession method")
//             },
//         }
//
//         // use mockedSessionHandler in code that requires SessionHandler
//         // and then make assertions.
//
//     }
type SessionHandlerMock struct {
	// CreateSessionFunc mocks the CreateSession method.
	CreateSessionFunc func(userID kallax.ULID) *Session

	// SaveSessionFunc mocks the SaveSession method.
	SaveSessionFunc func(session Session) Session

	// calls tracks calls to the methods.
	calls struct {
		// CreateSession holds details about calls to the CreateSession method.
		CreateSession []struct {
			// UserID is the userID argument value.
			UserID kallax.ULID
		}
		// SaveSession holds details about calls to the SaveSession method.
		SaveSession []struct {
			// Session is the session argument value.
			Session Session
		}
	}
}

// CreateSession calls CreateSessionFunc.
func (mock *SessionHandlerMock) CreateSession(userID kallax.ULID) *Session {
	if mock.CreateSessionFunc == nil {
		panic("SessionHandlerMock.CreateSessionFunc: method is nil but SessionHandler.CreateSession was just called")
	}
	callInfo := struct {
		UserID kallax.ULID
	}{
		UserID: userID,
	}
	lockSessionHandlerMockCreateSession.Lock()
	mock.calls.CreateSession = append(mock.calls.CreateSession, callInfo)
	lockSessionHandlerMockCreateSession.Unlock()
	return mock.CreateSessionFunc(userID)
}

// CreateSessionCalls gets all the calls that were made to CreateSession.
// Check the length with:
//     len(mockedSessionHandler.CreateSessionCalls())
func (mock *SessionHandlerMock) CreateSessionCalls() []struct {
	UserID kallax.ULID
} {
	var calls []struct {
		UserID kallax.ULID
	}
	lockSessionHandlerMockCreateSession.RLock()
	calls = mock.calls.CreateSession
	lockSessionHandlerMockCreateSession.RUnlock()
	return calls
}

// SaveSession calls SaveSessionFunc.
func (mock *SessionHandlerMock) SaveSession(session Session) Session {
	if mock.SaveSessionFunc == nil {
		panic("SessionHandlerMock.SaveSessionFunc: method is nil but SessionHandler.SaveSession was just called")
	}
	callInfo := struct {
		Session Session
	}{
		Session: session,
	}
	lockSessionHandlerMockSaveSession.Lock()
	mock.calls.SaveSession = append(mock.calls.SaveSession, callInfo)
	lockSessionHandlerMockSaveSession.Unlock()
	return mock.SaveSessionFunc(session)
}

// SaveSessionCalls gets all the calls that were made to SaveSession.
// Check the length with:
//     len(mockedSessionHandler.SaveSessionCalls())
func (mock *SessionHandlerMock) SaveSessionCalls() []struct {
	Session Session
} {
	var calls []struct {
		Session Session
	}
	lockSessionHandlerMockSaveSession.RLock()
	calls = mock.calls.SaveSession
	lockSessionHandlerMockSaveSession.RUnlock()
	return calls
}
