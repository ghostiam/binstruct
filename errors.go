package binstruct

import (
	"io"

	"github.com/pkg/errors"
)

var (
	ErrTagUnbalanced = errors.New("unbalanced square bracket")
)

func IsEOF(err error) bool {
	return errors.Cause(err) == io.EOF
}

func IsUnexpectedEOF(err error) bool {
	return errors.Cause(err) == io.ErrUnexpectedEOF
}
