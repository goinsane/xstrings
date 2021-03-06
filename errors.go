package xstrings

import (
	"errors"
	"fmt"
)

var (
	ErrCanNotGetAddr               = errors.New("can not get address of value")
	ErrNilPointer                  = errors.New("nil pointer error")
	ErrValueMustBeStruct           = errors.New("value must be struct")
	ErrArgumentCountExceeded       = errors.New("argument count exceeded")
	ErrArgumentStructFieldNotFound = errors.New("argument struct field not found")
)

// ParseError is type of error
type ParseError struct {
	err error
}

// newParseError wraps err into ParseError
func newParseError(err error) error {
	return &ParseError{
		err: err,
	}
}

// Error is implementation of error
func (e *ParseError) Error() string {
	str := "parse error"
	if e.err == nil || e.err.Error() == "" {
		return str
	}
	return fmt.Sprintf("%s: %v", str, e.err)
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
	return &FormatError{
		err: err,
	}
}

// Error is implementation of error
func (e *FormatError) Error() string {
	str := "format error"
	if e.err == nil || e.err.Error() == "" {
		return str
	}
	return fmt.Sprintf("%s: %v", str, e.err)
}

// Unwrap returns wrapped error
func (e *FormatError) Unwrap() error {
	return e.err
}

type MissingArgumentError struct {
	name string
	err  error
}

func (e *MissingArgumentError) Error() string {
	str := "missing argument"
	if e.name != "" {
		str = fmt.Sprintf("%s <%s>", str, e.name)
	}
	if e.err == nil || e.err.Error() == "" {
		return str
	}
	return fmt.Sprintf("%s: %v", str, e.err)
}

func (e *MissingArgumentError) Unwrap() error {
	return e.err
}

func (e *MissingArgumentError) Name() string {
	return e.name
}

type ArgumentParseError struct {
	name string
	err  error
}

func (e *ArgumentParseError) Error() string {
	str := "argument"
	if e.name != "" {
		str = fmt.Sprintf("%s <%s>", str, e.name)
	}
	str = fmt.Sprintf("%s parse error", str)
	if e.err == nil || e.err.Error() == "" {
		return str
	}
	return fmt.Sprintf("%s: %v", str, e.err)
}

func (e *ArgumentParseError) Unwrap() error {
	return e.err
}

func (e *ArgumentParseError) Name() string {
	return e.name
}
