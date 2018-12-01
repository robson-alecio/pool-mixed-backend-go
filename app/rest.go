package app

import (
	"encoding/json"
	"net/http"
)

//ProcessingBlock ...
type ProcessingBlock func(v interface{}) (interface{}, error)

//HTTPHelper ...
//go:generate moq -out httphelper_moq.go . HTTPHelper
type HTTPHelper interface {
	Process(v interface{}, blocks ...ProcessingBlock)
}

//HTTPHelperImpl ...
type HTTPHelperImpl struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
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
