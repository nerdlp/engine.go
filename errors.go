package engine

import (
	"net/http"

	"github.com/pkg/errors"
)

type socketError struct {
	innerError error
	code       int
	message    string
}

func (err socketError) Error() string {
	if err.innerError == nil {
		return err.message
	}
	return errors.Wrap(err.innerError, err.message).Error()
}

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
