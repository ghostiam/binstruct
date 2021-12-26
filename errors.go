package binstruct

import (
	"errors"
	"io"
)

// Deprecated: use errors.Is(err, io.EOF)
// IsEOF checks that the error is EOF
func IsEOF(err error) bool {
	return errors.Is(err, io.EOF)
}

// Deprecated: use errors.Is(err, io.ErrUnexpectedEOF)
// IsUnexpectedEOF checks that the error is Unexpected EOF
func IsUnexpectedEOF(err error) bool {
	return errors.Is(err, io.ErrUnexpectedEOF)
}
