package xstrings

import (
	"errors"
	"fmt"
)

// Error is type of error
type Error struct {
	err error
}

// newError wraps err into Error
func newError(err error) error {
	if err == nil {
		return nil
	}
	return &Error{
		err: err,
	}
}

// Error is implementation of error
func (e *Error) Error() string {
	if e.err == nil {
		return "string error"
	}
	return e.err.Error()
}

// Unwrap returns wrapped error
func (e *Error) Unwrap() error {
	return e.err
}

// ParseError is type of error
type ParseError struct {
	err error
}

// newParseError wraps err into ParseError
func newParseError(err error) error {
	if err == nil {
		return nil
	}
	return newError(&ParseError{
		err: err,
	})
}

// Error is implementation of error
func (e *ParseError) Error() string {
	s := "parse error"
	if e.err == nil {
		return s
	}
	return fmt.Sprintf("%s: %v", s, e.err)
}

// Unwrap returns wrapped error
func (e *ParseError) Unwrap() error {
	return e.err
}

// FormatError is type of error
type FormatError struct {
	err error
}

// newFormatError wraps err into FormatError
func newFormatError(err error) error {
	if err == nil {
		return nil
	}
	return newError(&FormatError{
		err: err,
	})
}

// Error is implementation of error
func (e *FormatError) Error() string {
	s := "format error"
	if e.err == nil {
		return s
	}
	return fmt.Sprintf("%s: %v", s, e.err)
}

// Unwrap returns wrapped error
func (e *FormatError) Unwrap() error {
	return e.err
}

var (
	ErrCanNotGetAddr = errors.New("can not get address of value")
	ErrNilPointer    = errors.New("nil pointer error")
)
