// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package app

import (
	"gopkg.in/src-d/go-kallax.v1"
	"sync"
)

var (
	lockHTTPHelperMockForbid              sync.RWMutex
	lockHTTPHelperMockGetRequestSessionID sync.RWMutex
	lockHTTPHelperMockIsRegisteredUser    sync.RWMutex
	lockHTTPHelperMockLoggedUserID        sync.RWMutex
	lockHTTPHelperMockProcess             sync.RWMutex
	lockHTTPHelperMockValidateSession     sync.RWMutex
)

// HTTPHelperMock is a mock implementation of HTTPHelper.
//
//     func TestSomethingThatUsesHTTPHelper(t *testing.T) {
//
//         // make and configure a mocked HTTPHelper
//         mockedHTTPHelper := &HTTPHelperMock{
//             ForbidFunc: func(in1 error)  {
// 	               panic("mock out the Forbid method")
//             },
//             GetRequestSessionIDFunc: func() (string, error) {
// 	               panic("mock out the GetRequestSessionID method")
//             },
//             IsRegisteredUserFunc: func() bool {
// 	               panic("mock out the IsRegisteredUser method")
//             },
//             LoggedUserIDFunc: func() kallax.ULID {
// 	               panic("mock out the LoggedUserID method")
//             },
//             ProcessFunc: func(in1 interface{}, in2 ...ProcessingBlock)  {
// 	               panic("mock out the Process method")
//             },
//             ValidateSessionFunc: func() error {
// 	               panic("mock out the ValidateSession method")
//             },
//         }
//
//         // use mockedHTTPHelper in code that requires HTTPHelper
//         // and then make assertions.
//
//     }
type HTTPHelperMock struct {
	// ForbidFunc mocks the Forbid method.
	ForbidFunc func(in1 error)

	// GetRequestSessionIDFunc mocks the GetRequestSessionID method.
	GetRequestSessionIDFunc func() (string, error)

	// IsRegisteredUserFunc mocks the IsRegisteredUser method.
	IsRegisteredUserFunc func() bool

	// LoggedUserIDFunc mocks the LoggedUserID method.
	LoggedUserIDFunc func() kallax.ULID

	// ProcessFunc mocks the Process method.
	ProcessFunc func(in1 interface{}, in2 ...ProcessingBlock)

	// ValidateSessionFunc mocks the ValidateSession method.
	ValidateSessionFunc func() error

	// calls tracks calls to the methods.
	calls struct {
		// Forbid holds details about calls to the Forbid method.
		Forbid []struct {
			// In1 is the in1 argument value.
			In1 error
		}
		// GetRequestSessionID holds details about calls to the GetRequestSessionID method.
		GetRequestSessionID []struct {
		}
		// IsRegisteredUser holds details about calls to the IsRegisteredUser method.
		IsRegisteredUser []struct {
		}
		// LoggedUserID holds details about calls to the LoggedUserID method.
		LoggedUserID []struct {
		}
		// Process holds details about calls to the Process method.
		Process []struct {
			// In1 is the in1 argument value.
			In1 interface{}
			// In2 is the in2 argument value.
			In2 []ProcessingBlock
		}
		// ValidateSession holds details about calls to the ValidateSession method.
		ValidateSession []struct {
		}
	}
}

// Forbid calls ForbidFunc.
func (mock *HTTPHelperMock) Forbid(in1 error) {
	if mock.ForbidFunc == nil {
		panic("HTTPHelperMock.ForbidFunc: method is nil but HTTPHelper.Forbid was just called")
	}
	callInfo := struct {
		In1 error
	}{
		In1: in1,
	}
	lockHTTPHelperMockForbid.Lock()
	mock.calls.Forbid = append(mock.calls.Forbid, callInfo)
	lockHTTPHelperMockForbid.Unlock()
	mock.ForbidFunc(in1)
}

// ForbidCalls gets all the calls that were made to Forbid.
// Check the length with:
//     len(mockedHTTPHelper.ForbidCalls())
func (mock *HTTPHelperMock) ForbidCalls() []struct {
	In1 error
} {
	var calls []struct {
		In1 error
	}
	lockHTTPHelperMockForbid.RLock()
	calls = mock.calls.Forbid
	lockHTTPHelperMockForbid.RUnlock()
	return calls
}

// GetRequestSessionID calls GetRequestSessionIDFunc.
func (mock *HTTPHelperMock) GetRequestSessionID() (string, error) {
	if mock.GetRequestSessionIDFunc == nil {
		panic("HTTPHelperMock.GetRequestSessionIDFunc: method is nil but HTTPHelper.GetRequestSessionID was just called")
	}
	callInfo := struct {
	}{}
	lockHTTPHelperMockGetRequestSessionID.Lock()
	mock.calls.GetRequestSessionID = append(mock.calls.GetRequestSessionID, callInfo)
	lockHTTPHelperMockGetRequestSessionID.Unlock()
	return mock.GetRequestSessionIDFunc()
}

// GetRequestSessionIDCalls gets all the calls that were made to GetRequestSessionID.
// Check the length with:
//     len(mockedHTTPHelper.GetRequestSessionIDCalls())
func (mock *HTTPHelperMock) GetRequestSessionIDCalls() []struct {
} {
	var calls []struct {
	}
	lockHTTPHelperMockGetRequestSessionID.RLock()
	calls = mock.calls.GetRequestSessionID
	lockHTTPHelperMockGetRequestSessionID.RUnlock()
	return calls
}

// IsRegisteredUser calls IsRegisteredUserFunc.
func (mock *HTTPHelperMock) IsRegisteredUser() bool {
	if mock.IsRegisteredUserFunc == nil {
		panic("HTTPHelperMock.IsRegisteredUserFunc: method is nil but HTTPHelper.IsRegisteredUser was just called")
	}
	callInfo := struct {
	}{}
	lockHTTPHelperMockIsRegisteredUser.Lock()
	mock.calls.IsRegisteredUser = append(mock.calls.IsRegisteredUser, callInfo)
	lockHTTPHelperMockIsRegisteredUser.Unlock()
	return mock.IsRegisteredUserFunc()
}

// IsRegisteredUserCalls gets all the calls that were made to IsRegisteredUser.
// Check the length with:
//     len(mockedHTTPHelper.IsRegisteredUserCalls())
func (mock *HTTPHelperMock) IsRegisteredUserCalls() []struct {
} {
	var calls []struct {
	}
	lockHTTPHelperMockIsRegisteredUser.RLock()
	calls = mock.calls.IsRegisteredUser
	lockHTTPHelperMockIsRegisteredUser.RUnlock()
	return calls
}

// LoggedUserID calls LoggedUserIDFunc.
func (mock *HTTPHelperMock) LoggedUserID() kallax.ULID {
	if mock.LoggedUserIDFunc == nil {
		panic("HTTPHelperMock.LoggedUserIDFunc: method is nil but HTTPHelper.LoggedUserID was just called")
	}
	callInfo := struct {
	}{}
	lockHTTPHelperMockLoggedUserID.Lock()
	mock.calls.LoggedUserID = append(mock.calls.LoggedUserID, callInfo)
	lockHTTPHelperMockLoggedUserID.Unlock()
	return mock.LoggedUserIDFunc()
}

// LoggedUserIDCalls gets all the calls that were made to LoggedUserID.
// Check the length with:
//     len(mockedHTTPHelper.LoggedUserIDCalls())
func (mock *HTTPHelperMock) LoggedUserIDCalls() []struct {
} {
	var calls []struct {
	}
	lockHTTPHelperMockLoggedUserID.RLock()
	calls = mock.calls.LoggedUserID
	lockHTTPHelperMockLoggedUserID.RUnlock()
	return calls
}

// Process calls ProcessFunc.
func (mock *HTTPHelperMock) Process(in1 interface{}, in2 ...ProcessingBlock) {
	if mock.ProcessFunc == nil {
		panic("HTTPHelperMock.ProcessFunc: method is nil but HTTPHelper.Process was just called")
	}
	callInfo := struct {
		In1 interface{}
		In2 []ProcessingBlock
	}{
		In1: in1,
		In2: in2,
	}
	lockHTTPHelperMockProcess.Lock()
	mock.calls.Process = append(mock.calls.Process, callInfo)
	lockHTTPHelperMockProcess.Unlock()
	mock.ProcessFunc(in1, in2...)
}

// ProcessCalls gets all the calls that were made to Process.
// Check the length with:
//     len(mockedHTTPHelper.ProcessCalls())
func (mock *HTTPHelperMock) ProcessCalls() []struct {
	In1 interface{}
	In2 []ProcessingBlock
} {
	var calls []struct {
		In1 interface{}
		In2 []ProcessingBlock
	}
	lockHTTPHelperMockProcess.RLock()
	calls = mock.calls.Process
	lockHTTPHelperMockProcess.RUnlock()
	return calls
}

// ValidateSession calls ValidateSessionFunc.
func (mock *HTTPHelperMock) ValidateSession() error {
	if mock.ValidateSessionFunc == nil {
		panic("HTTPHelperMock.ValidateSessionFunc: method is nil but HTTPHelper.ValidateSession was just called")
	}
	callInfo := struct {
	}{}
	lockHTTPHelperMockValidateSession.Lock()
	mock.calls.ValidateSession = append(mock.calls.ValidateSession, callInfo)
	lockHTTPHelperMockValidateSession.Unlock()
	return mock.ValidateSessionFunc()
}

// ValidateSessionCalls gets all the calls that were made to ValidateSession.
// Check the length with:
//     len(mockedHTTPHelper.ValidateSessionCalls())
func (mock *HTTPHelperMock) ValidateSessionCalls() []struct {
} {
	var calls []struct {
	}
	lockHTTPHelperMockValidateSession.RLock()
	calls = mock.calls.ValidateSession
	lockHTTPHelperMockValidateSession.RUnlock()
	return calls
}
