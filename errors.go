package engine

import "github.com/pkg/errors"

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
