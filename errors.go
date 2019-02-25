package binstruct

import (
	"io"

	"github.com/pkg/errors"
)

// IsEOF checks that the error is EOF
func IsEOF(err error) bool {
	return errors.Cause(err) == io.EOF
}

// IsUnexpectedEOF checks that the error is Unexpected EOF
func IsUnexpectedEOF(err error) bool {
	return errors.Cause(err) == io.ErrUnexpectedEOF
}
