package command

import (
	"strings"
)

func Unmarshal(cmd Command, args ...string) error {
	if err := getArgumentStruct(cmd).Unmarshal(cmd, args...); err != nil {
		return err
	}
	return nil
}

func Find(cmds []Command, args ...string) (int, error) {
	if err := checkArgs(args...); err != nil {
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

func FindCmd(cmds []Command, args ...string) (Command, error) {
	idx, err := Find(cmds, args...)
	if err != nil {
		return nil, err
	}
	return cmds[idx], nil
}

func FindAndUnmarshal(cmds []Command, args ...string) (Command, error) {
	cmd, err := FindCmd(cmds, args...)
	if err != nil {
		return nil, err
	}
	return cmd, Unmarshal(cmd, args...)
}

func Usage(cmd Command, cmdName string, prefix string) (string, error) {
	usage, err := ParameterUsage(cmd)
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

func ParameterUsage(cmd Command) (string, error) {
	fields, err := getArgumentStruct(cmd).Fields(cmd)
	if err != nil {
		return "", err
	}
	if len(fields) > 0 {
		fields = fields[1:]
	}
	return fields.String(), nil
}

func checkArgs(args ...string) error {
	sizeArgs := len(args)
	if sizeArgs <= 0 {
		return ErrCommandNotSet
	}
	if args[0] == "" {
		return ErrCommandNotSet
	}
	return nil
}
