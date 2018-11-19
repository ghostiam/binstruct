package binstruct

import (
	"io"

	"github.com/pkg/errors"
)

func IsEOF(err error) bool {
	return errors.Cause(err) == io.EOF
}

func IsUnexpectedEOF(err error) bool {
	return errors.Cause(err) == io.ErrUnexpectedEOF
}
