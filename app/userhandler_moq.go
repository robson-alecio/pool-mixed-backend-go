// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package app

import (
	"gopkg.in/src-d/go-kallax.v1"
	"sync"
)

var (
	lockUserHandlerMockCreateAnonUser     sync.RWMutex
	lockUserHandlerMockCreateUserFromData sync.RWMutex
	lockUserHandlerMockFindUserByID       sync.RWMutex
	lockUserHandlerMockFindUserByLogin    sync.RWMutex
	lockUserHandlerMockSaveUser           sync.RWMutex
)

// UserHandlerMock is a mock implementation of UserHandler.
//
//     func TestSomethingThatUsesUserHandler(t *testing.T) {
//
//         // make and configure a mocked UserHandler
//         mockedUserHandler := &UserHandlerMock{
//             CreateAnonUserFunc: func() User {
// 	               panic("mock out the CreateAnonUser method")
//             },
//             CreateUserFromDataFunc: func(d *UserCreationData) (User, error) {
// 	               panic("mock out the CreateUserFromData method")
//             },
//             FindUserByIDFunc: func(ID kallax.ULID) (*User, error) {
// 	               panic("mock out the FindUserByID method")
//             },
//             FindUserByLoginFunc: func(login string) (*User, error) {
// 	               panic("mock out the FindUserByLogin method")
//             },
//             SaveUserFunc: func(user User) User {
// 	               panic("mock out the SaveUser method")
//             },
//         }
//
//         // use mockedUserHandler in code that requires UserHandler
//         // and then make assertions.
//
//     }
type UserHandlerMock struct {
	// CreateAnonUserFunc mocks the CreateAnonUser method.
	CreateAnonUserFunc func() User

	// CreateUserFromDataFunc mocks the CreateUserFromData method.
	CreateUserFromDataFunc func(d *UserCreationData) (User, error)

	// FindUserByIDFunc mocks the FindUserByID method.
	FindUserByIDFunc func(ID kallax.ULID) (*User, error)

	// FindUserByLoginFunc mocks the FindUserByLogin method.
	FindUserByLoginFunc func(login string) (*User, error)

	// SaveUserFunc mocks the SaveUser method.
	SaveUserFunc func(user User) User

	// calls tracks calls to the methods.
	calls struct {
		// CreateAnonUser holds details about calls to the CreateAnonUser method.
		CreateAnonUser []struct {
		}
		// CreateUserFromData holds details about calls to the CreateUserFromData method.
		CreateUserFromData []struct {
			// D is the d argument value.
			D *UserCreationData
		}
		// FindUserByID holds details about calls to the FindUserByID method.
		FindUserByID []struct {
			// ID is the ID argument value.
			ID kallax.ULID
		}
		// FindUserByLogin holds details about calls to the FindUserByLogin method.
		FindUserByLogin []struct {
			// Login is the login argument value.
			Login string
		}
		// SaveUser holds details about calls to the SaveUser method.
		SaveUser []struct {
			// User is the user argument value.
			User User
		}
	}
}

// CreateAnonUser calls CreateAnonUserFunc.
func (mock *UserHandlerMock) CreateAnonUser() User {
	if mock.CreateAnonUserFunc == nil {
		panic("UserHandlerMock.CreateAnonUserFunc: method is nil but UserHandler.CreateAnonUser was just called")
	}
	callInfo := struct {
	}{}
	lockUserHandlerMockCreateAnonUser.Lock()
	mock.calls.CreateAnonUser = append(mock.calls.CreateAnonUser, callInfo)
	lockUserHandlerMockCreateAnonUser.Unlock()
	return mock.CreateAnonUserFunc()
}

// CreateAnonUserCalls gets all the calls that were made to CreateAnonUser.
// Check the length with:
//     len(mockedUserHandler.CreateAnonUserCalls())
func (mock *UserHandlerMock) CreateAnonUserCalls() []struct {
} {
	var calls []struct {
	}
	lockUserHandlerMockCreateAnonUser.RLock()
	calls = mock.calls.CreateAnonUser
	lockUserHandlerMockCreateAnonUser.RUnlock()
	return calls
}

// CreateUserFromData calls CreateUserFromDataFunc.
func (mock *UserHandlerMock) CreateUserFromData(d *UserCreationData) (User, error) {
	if mock.CreateUserFromDataFunc == nil {
		panic("UserHandlerMock.CreateUserFromDataFunc: method is nil but UserHandler.CreateUserFromData was just called")
	}
	callInfo := struct {
		D *UserCreationData
	}{
		D: d,
	}
	lockUserHandlerMockCreateUserFromData.Lock()
	mock.calls.CreateUserFromData = append(mock.calls.CreateUserFromData, callInfo)
	lockUserHandlerMockCreateUserFromData.Unlock()
	return mock.CreateUserFromDataFunc(d)
}

// CreateUserFromDataCalls gets all the calls that were made to CreateUserFromData.
// Check the length with:
//     len(mockedUserHandler.CreateUserFromDataCalls())
func (mock *UserHandlerMock) CreateUserFromDataCalls() []struct {
	D *UserCreationData
} {
	var calls []struct {
		D *UserCreationData
	}
	lockUserHandlerMockCreateUserFromData.RLock()
	calls = mock.calls.CreateUserFromData
	lockUserHandlerMockCreateUserFromData.RUnlock()
	return calls
}

// FindUserByID calls FindUserByIDFunc.
func (mock *UserHandlerMock) FindUserByID(ID kallax.ULID) (*User, error) {
	if mock.FindUserByIDFunc == nil {
		panic("UserHandlerMock.FindUserByIDFunc: method is nil but UserHandler.FindUserByID was just called")
	}
	callInfo := struct {
		ID kallax.ULID
	}{
		ID: ID,
	}
	lockUserHandlerMockFindUserByID.Lock()
	mock.calls.FindUserByID = append(mock.calls.FindUserByID, callInfo)
	lockUserHandlerMockFindUserByID.Unlock()
	return mock.FindUserByIDFunc(ID)
}

// FindUserByIDCalls gets all the calls that were made to FindUserByID.
// Check the length with:
//     len(mockedUserHandler.FindUserByIDCalls())
func (mock *UserHandlerMock) FindUserByIDCalls() []struct {
	ID kallax.ULID
} {
	var calls []struct {
		ID kallax.ULID
	}
	lockUserHandlerMockFindUserByID.RLock()
	calls = mock.calls.FindUserByID
	lockUserHandlerMockFindUserByID.RUnlock()
	return calls
}

// FindUserByLogin calls FindUserByLoginFunc.
func (mock *UserHandlerMock) FindUserByLogin(login string) (*User, error) {
	if mock.FindUserByLoginFunc == nil {
		panic("UserHandlerMock.FindUserByLoginFunc: method is nil but UserHandler.FindUserByLogin was just called")
	}
	callInfo := struct {
		Login string
	}{
		Login: login,
	}
	lockUserHandlerMockFindUserByLogin.Lock()
	mock.calls.FindUserByLogin = append(mock.calls.FindUserByLogin, callInfo)
	lockUserHandlerMockFindUserByLogin.Unlock()
	return mock.FindUserByLoginFunc(login)
}

// FindUserByLoginCalls gets all the calls that were made to FindUserByLogin.
// Check the length with:
//     len(mockedUserHandler.FindUserByLoginCalls())
func (mock *UserHandlerMock) FindUserByLoginCalls() []struct {
	Login string
} {
	var calls []struct {
		Login string
	}
	lockUserHandlerMockFindUserByLogin.RLock()
	calls = mock.calls.FindUserByLogin
	lockUserHandlerMockFindUserByLogin.RUnlock()
	return calls
}

// SaveUser calls SaveUserFunc.
func (mock *UserHandlerMock) SaveUser(user User) User {
	if mock.SaveUserFunc == nil {
		panic("UserHandlerMock.SaveUserFunc: method is nil but UserHandler.SaveUser was just called")
	}
	callInfo := struct {
		User User
	}{
		User: user,
	}
	lockUserHandlerMockSaveUser.Lock()
	mock.calls.SaveUser = append(mock.calls.SaveUser, callInfo)
	lockUserHandlerMockSaveUser.Unlock()
	return mock.SaveUserFunc(user)
}

// SaveUserCalls gets all the calls that were made to SaveUser.
// Check the length with:
//     len(mockedUserHandler.SaveUserCalls())
func (mock *UserHandlerMock) SaveUserCalls() []struct {
	User User
} {
	var calls []struct {
		User User
	}
	lockUserHandlerMockSaveUser.RLock()
	calls = mock.calls.SaveUser
	lockUserHandlerMockSaveUser.RUnlock()
	return calls
}