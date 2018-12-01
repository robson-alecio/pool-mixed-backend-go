package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

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
