package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/src-d/go-kallax.v1"
)

//ProcessingBlock ...
type ProcessingBlock func(v interface{}) (interface{}, error)

//HTTPHelper ...
//go:generate moq -out httphelper_moq.go . HTTPHelper
type HTTPHelper interface {
	Process(interface{}, ...ProcessingBlock)
	ValidateSession() error
	GetRequestSessionID() (string, error)
	IsRegisteredUser() bool
	Forbid(error)
	LoggedUserID() kallax.ULID
	GetVar(string) string
}

//HTTPHelperImpl ...
type HTTPHelperImpl struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	CheckSession   func(ID string) error
	Session        *Session
}

//NewHTTPHelper ...
func NewHTTPHelper(w http.ResponseWriter, r *http.Request) HTTPHelperImpl {
	return HTTPHelperImpl{
		Request:        r,
		ResponseWriter: w,
	}
}

//Process ...
func (h HTTPHelperImpl) Process(v interface{}, blocks ...ProcessingBlock) {
	err := json.NewDecoder(h.Request.Body).Decode(&v)

	if err != nil {
		http.Error(h.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	var result interface{} = v
	var aErr error

	for _, f := range blocks {
		result, aErr = f(result)

		if aErr != nil {
			http.Error(h.ResponseWriter, aErr.Error(), http.StatusConflict)
			return
		}
	}

	json.NewEncoder(h.ResponseWriter).Encode(result)
}

//ValidateSession ...
func (h HTTPHelperImpl) ValidateSession() error {
	ID, err := h.GetRequestSessionID()
	if err != nil {
		return err
	}

	return h.CheckSession(ID)
}

//GetRequestSessionID ...
func (h HTTPHelperImpl) GetRequestSessionID() (string, error) {
	sessionID := h.Request.Header.Get("sessionId")

	if sessionID == "" {
		return "", ErrUserNotLogged("Must be logged to perform this action. Missing value.")
	}

	return sessionID, nil
}

//IsRegisteredUser ...
func (h HTTPHelperImpl) IsRegisteredUser() bool {
	return h.Session != nil && h.Session.RegisteredUser
}

//Forbid ...
func (h HTTPHelperImpl) Forbid(err error) {
	http.Error(h.ResponseWriter, err.Error(), http.StatusForbidden)
}

//LoggedUserID ...
func (h HTTPHelperImpl) LoggedUserID() kallax.ULID {
	return h.Session.UserID
}

//GetVar ...
func (h HTTPHelperImpl) GetVar(name string) string {
	return mux.Vars(h.Request)["id"]
}
