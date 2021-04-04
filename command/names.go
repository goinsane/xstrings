package command

import (
	"strings"
)

type Names interface {
	CmdNames() []string
	CmdNamesFold() bool
	Is(cmdName string) bool
}

func NewNames(cmdNamesFold bool, cmdNames ...string) Names {
	n := &namesStruct{
		cmdNames:     make([]string, len(cmdNames)),
		cmdNamesFold: cmdNamesFold,
	}
	copy(n.cmdNames, cmdNames)
	return n
}

type namesStruct struct {
	cmdNames     []string
	cmdNamesFold bool
}

func (n *namesStruct) CmdNames() []string {
	result := make([]string, len(n.cmdNames))
	copy(result, n.cmdNames)
	return result
}

func (n *namesStruct) CmdNamesFold() bool {
	return n.cmdNamesFold
}

func (n *namesStruct) Is(cmdName string) bool {
	for i := range n.cmdNames {
		if n.cmdNamesFold {
			if strings.EqualFold(cmdName, n.cmdNames[i]) {
				return true
			}
		} else {
			if cmdName == n.cmdNames[i] {
				return true
			}
		}
	}
	return false
}
