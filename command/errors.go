package command

import (
	"errors"
	"fmt"
)

var (
	ErrCommandNotSet = errors.New("command not set")
)

type UnknownCommandError struct {
	name string
	err  error
}

func (e *UnknownCommandError) Error() string {
	str := "unknown command"
	if e.name != "" {
		str = fmt.Sprintf("%s %q", str, e.name)
	}
	if e.err == nil || e.err.Error() == "" {
		return str
	}
	return fmt.Sprintf("%s: %v", str, e.err)
}

func (e *UnknownCommandError) Unwrap() error {
	return e.err
}

func (e *UnknownCommandError) Name() string {
	return e.name
}
