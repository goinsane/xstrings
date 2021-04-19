package command

import (
	"strings"

	"github.com/goinsane/xstrings"
)

type Handler struct {
	Unmarshaler              *xstrings.Unmarshaler
	FieldNameBeginsLowerCase bool
	FieldNameFold            bool
	FieldTagKey              string
}

func (h *Handler) Unmarshal(cmd Command, args ...string) error {
	if err := h.getArgumentStruct(cmd).Unmarshal(cmd, args...); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Find(cmds []Command, args ...string) (int, error) {
	if err := h.checkArgs(args...); err != nil {
		return -1, err
	}
	cmdName := args[0]
	for idx, cmd := range cmds {
		if cmd.Is(cmdName) {
			return idx, nil
		}
	}
	return -1, &UnknownCommandError{cmdName, nil}
}

func (h *Handler) FindCmd(cmds []Command, args ...string) (Command, error) {
	idx, err := h.Find(cmds, args...)
	if err != nil {
		return nil, err
	}
	return cmds[idx], nil
}

func (h *Handler) FindAndUnmarshal(cmds []Command, args ...string) (Command, error) {
	cmd, err := h.FindCmd(cmds, args...)
	if err != nil {
		return nil, err
	}
	return cmd, h.Unmarshal(cmd, args...)
}

func (h *Handler) Usage(cmd Command, cmdName string, prefix string) (string, error) {
	usage, err := h.ParameterUsage(cmd)
	if err != nil {
		return "", err
	}
	if cmdName != "" && !cmd.Is(cmdName) {
		return "", &UnknownCommandError{cmdName, nil}
	}
	result := ""
	for idx, cmdName2 := range cmd.CmdNames() {
		if cmdName != "" {
			if cmd.CmdNamesFold() {
				if !strings.EqualFold(cmdName, cmdName2) {
					continue
				}
			} else {
				if cmdName != cmdName2 {
					continue
				}
			}
		}
		nl := ""
		if idx > 0 {
			nl = "\n"
		}
		result += nl + prefix + cmdName2 + " " + usage
	}
	return result, nil
}

func (h *Handler) ParameterUsage(cmd Command) (string, error) {
	fields, err := h.getArgumentStruct(cmd).Fields(cmd)
	if err != nil {
		return "", err
	}
	if len(fields) > 0 {
		fields = fields[1:]
	}
	return fields.String(), nil
}

func (h *Handler) checkArgs(args ...string) error {
	sizeArgs := len(args)
	if sizeArgs <= 0 {
		return ErrCommandNotSet
	}
	if args[0] == "" {
		return ErrCommandNotSet
	}
	return nil
}

func (h *Handler) getArgumentStruct(cmd Command) *xstrings.ArgumentStruct {
	return &xstrings.ArgumentStruct{
		Unmarshaler:              h.Unmarshaler,
		FieldNameBeginsLowerCase: h.FieldNameBeginsLowerCase,
		FieldNameFold:            h.FieldNameFold,
		FieldTagKey:              h.FieldTagKey,
		FieldOffset:              cmd.FieldOffset(),
		ArgCountMin:              cmd.ArgCountMin(),
		ArgCountMax:              cmd.ArgCountMax(),
	}
}
