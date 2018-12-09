package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/chai2010/assert"
)

type FakeResponseWriter struct {
	FakeWriter *bytes.Buffer
	FakeHeader http.Header
	FakeStatus int
}

func (w FakeResponseWriter) Header() http.Header {
	return w.FakeHeader
}

func (w FakeResponseWriter) Write(x []byte) (int, error) {
	return w.FakeWriter.Write(x)
}

func (w FakeResponseWriter) WriteHeader(statusCode int) {
	w.FakeStatus = statusCode
}

type JSONReader struct {
	InnerReader io.Reader
}

func (jr JSONReader) Read(b []byte) (n int, err error) {
	return jr.InnerReader.Read(b)
}

func (jr JSONReader) Close() error {
	return nil
}

type FakeData struct {
	Color string `json:"color,omitempty"`
	Name  string `json:"name,omitempty"`
}

type FakeObject struct {
	Color string
	Name  string
}

func TestProcessHappyDay(t *testing.T) {
	result := bytes.NewBuffer(make([]byte, 0))
	writer := FakeResponseWriter{
		FakeHeader: make(http.Header, 0),
		FakeWriter: result,
	}

	reader := JSONReader{
		InnerReader: strings.NewReader(`{"color":"#0000ff","name":"Blue"}`),
	}
	helper := &HTTPHelperImpl{
		ResponseWriter: writer,
		Request: &http.Request{
			Body: reader,
		},
	}

	convert := func(v interface{}) (interface{}, error) {
		data := v.(*FakeData)

		return &FakeObject{
			Color: data.Color,
			Name:  data.Name,
		}, nil
	}

	stringfy := func(v interface{}) (interface{}, error) {
		object := v.(*FakeObject)

		message := fmt.Sprintf("My name is %s and my skin is %s", object.Name, object.Color)

		stringfyResult := map[string]string{
			"message": message,
		}

		return stringfyResult, nil
	}

	helper.Process(&FakeData{}, convert, stringfy)

	expected := `{"message":"My name is Blue and my skin is #0000ff"}`
	assert.AssertEqual(t, expected, strings.TrimSpace(result.String()))
}

func TestProcessErrorOnProcessing(t *testing.T) {
	result := bytes.NewBuffer(make([]byte, 0))
	writer := FakeResponseWriter{
		FakeHeader: make(http.Header, 0),
		FakeWriter: result,
	}

	reader := JSONReader{
		InnerReader: strings.NewReader(`{"color":"#0000ff","name":"Blue"}`),
	}
	helper := &HTTPHelperImpl{
		ResponseWriter: writer,
		Request: &http.Request{
			Body: reader,
		},
	}

	convert := func(v interface{}) (interface{}, error) {
		return nil, fmt.Errorf("Nothing good happens after 2:00 PM")
	}

	stringfy := func(v interface{}) (interface{}, error) {
		object := v.(*FakeObject)

		message := fmt.Sprintf("My name is %s and my skin is %s", object.Name, object.Color)

		stringfyResult := map[string]string{
			"message": message,
		}

		return stringfyResult, nil
	}

	helper.Process(&FakeData{}, convert, stringfy)

	expected := "Nothing good happens after 2:00 PM"
	assert.AssertEqual(t, expected, strings.TrimSpace(result.String()))
}

func TestProcessWithInvalidJSON(t *testing.T) {
	result := bytes.NewBuffer(make([]byte, 0))
	writer := FakeResponseWriter{
		FakeHeader: make(http.Header, 0),
		FakeWriter: result,
	}

	reader := JSONReader{
		InnerReader: strings.NewReader("Kabuuuuuuuummmmmmmmmmmm"),
	}
	helper := &HTTPHelperImpl{
		ResponseWriter: writer,
		Request: &http.Request{
			Body: reader,
		},
	}

	convert := func(v interface{}) (interface{}, error) {
		return &FakeObject{}, nil
	}

	helper.Process(&FakeData{}, convert)

	expected := "invalid character 'K' looking for beginning of value"
	assert.AssertEqual(t, expected, strings.TrimSpace(result.String()))
}

func TestValidateSession(t *testing.T) {
	result := bytes.NewBuffer(make([]byte, 0))
	writer := FakeResponseWriter{
		FakeHeader: make(http.Header, 0),
		FakeWriter: result,
	}

	header := http.Header{
		textproto.CanonicalMIMEHeaderKey("sessionId"): []string{"7d97abb1-2f1b-4542-8173-67e78a590ab9"},
	}

	helper := &HTTPHelperImpl{
		ResponseWriter: writer,
		Request: &http.Request{
			Header: header,
		},
		CheckSession: func(ID string) error {
			return nil
		},
	}

	err := helper.ValidateSession()
	assert.AssertTrue(t, err == nil)
}

func TestValidateSessionFailWhenNotHaveSessionId(t *testing.T) {
	header := http.Header{
		textproto.CanonicalMIMEHeaderKey("sessionId"): []string{},
	}

	helper := &HTTPHelperImpl{
		Request: &http.Request{
			Header: header,
		},
		CheckSession: func(ID string) error {
			return nil
		},
	}

	err := helper.ValidateSession()
	assert.AssertEqual(t, "Must be logged to perform this action. Missing value.", err.Error())
}

func TestIsRegisteredUser(t *testing.T) {
	helper := &HTTPHelperImpl{
		Session: &Session{
			RegisteredUser: true,
		},
	}

	result := helper.IsRegisteredUser()
	assert.AssertTrue(t, result)
}

func TestIsNotRegisteredUser(t *testing.T) {
	helper := &HTTPHelperImpl{
		Session: &Session{
			RegisteredUser: false,
		},
	}

	result := helper.IsRegisteredUser()
	assert.AssertFalse(t, result)
}
func TestIsNotRegisteredUserWithoutSession(t *testing.T) {
	helper := &HTTPHelperImpl{}

	result := helper.IsRegisteredUser()
	assert.AssertFalse(t, result)
}

func TestForbid(t *testing.T) {
	result := bytes.NewBuffer(make([]byte, 0))
	writer := FakeResponseWriter{
		FakeHeader: make(http.Header, 0),
		FakeWriter: result,
	}

	helper := &HTTPHelperImpl{
		ResponseWriter: writer,
	}

	helper.Forbid(fmt.Errorf("Baba Yaga"))

	expected := "Baba Yaga"
	assert.AssertEqual(t, expected, strings.TrimSpace(result.String()))
}

func TestLoggedUserID(t *testing.T) {
	userID := kallax.NewULID()
	helper := &HTTPHelperImpl{
		Session: &Session{
			UserID: userID,
		},
	}

	result := helper.LoggedUserID()
	assert.AssertEqual(t, userID, result)
}
func TestLoggedUserIDWithoutSession(t *testing.T) {
	helper := &HTTPHelperImpl{}

	first := helper.LoggedUserID()
	second := helper.LoggedUserID()
	assert.AssertNotEqual(t, first, second)
}

type FakeContext struct {
	Values map[string]string
}

func (f FakeContext) Deadline() (deadline time.Time, ok bool) {
	return time.Now().Add(time.Duration(6 * time.Hour)), true
}

func (f FakeContext) Done() <-chan struct{} {
	return make(chan struct{})
}

func (f FakeContext) Err() error {
	return nil
}

func (f FakeContext) Value(key interface{}) interface{} {
	return f.Values
}

type contextKey int

const (
	varsKey contextKey = iota
)

func TestGetVar(t *testing.T) {
	fakeContext := FakeContext{
		Values: map[string]string{
			"something": "I would like a dinner reservation for midnight",
		},
	}

	request := (&http.Request{}).WithContext(fakeContext)

	helper := &HTTPHelperImpl{
		Request: request,
	}

	result := helper.GetVar("something")
	assert.AssertEqual(t, "I would like a dinner reservation for midnight", result)
}
