package engine

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	contentType = "Content-Type"
)

const (
	contentTypeApplicationJson = "application/json"
)

func handleError(w http.ResponseWriter, err error) {
	serr, ok := err.(socketError)
	if !ok {
		serr = socketError{
			innerError: err,
			code:       http.StatusInternalServerError,
			message:    "internal server error",
		}
	}

	w.Header().Add(contentType, contentTypeApplicationJson)
	w.WriteHeader(serr.code)
	w.Write([]byte(err.Error()))
}

type response struct {
	status int
	body   any
}

func handleResponse(w http.ResponseWriter, response *response) {
	body, err := json.Marshal(response.body)
	if err != nil {
		err = socketError{
			innerError: err,
			code:       http.StatusInternalServerError,
			message:    "fail to marshal response",
		}
		slog.Error(err.Error())
		handleError(w, err)
		return
	}

	w.Header().Add(contentType, contentTypeApplicationJson)
	w.WriteHeader(response.status)
	w.Write(body)
}
