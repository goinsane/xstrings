package xstrings

import (
	"errors"
	"fmt"
)

// ParseError is type of error
type ParseError struct {
	err error
}

// newParseError wraps err into ParseError
func newParseError(err error) error {
	if err == nil {
		return nil
	}
	return &ParseError{
		err: err,
	}
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
	return &FormatError{
		err: err,
	}
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

type MissingArgumentError struct {
	name string
}

func (e *MissingArgumentError) Error() string {
	return fmt.Sprintf("missing argument of <%s>", e.name)
}

func (e *MissingArgumentError) Unwrap() error {
	return nil
}

func (e *MissingArgumentError) ArgName() string {
	return e.name
}

type ArgumentParseError struct {
	name string
	err  error
}

func (e *ArgumentParseError) Error() string {
	return fmt.Sprintf("argument of <%s> parse error: %v", e.name, e.err)
}

func (e *ArgumentParseError) Unwrap() error {
	return e.err
}

func (e *ArgumentParseError) ArgName() string {
	return e.name
}

var (
	ErrCanNotGetAddr         = errors.New("can not get address of value")
	ErrNilPointer            = errors.New("nil pointer error")
	ErrValueMustBeStruct     = errors.New("value must be struct")
	ErrArgumentCountExceeded = errors.New("argument count exceeded")
)
