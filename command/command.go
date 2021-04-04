package command

type Command interface {
	Runnable
	Names
	Info
}

func New(runnable Runnable, fieldOffset, argCountMin, argCountMax int, cmdNamesFold bool, cmdNames ...string) Command {
	return &commandStruct{
		Runnable:    runnable,
		Names:       NewNames(cmdNamesFold, cmdNames...),
		fieldOffset: fieldOffset,
		argCountMin: argCountMin,
		argCountMax: argCountMax,
	}
}

func NewWithRunFunc(unmarshalInterface UnmarshalInterface, f RunFunc, fieldOffset, argCountMin, argCountMax int, cmdNamesFold bool, cmdNames ...string) Command {
	return New(NewRunnable(unmarshalInterface, f), fieldOffset, argCountMin, argCountMax, cmdNamesFold, cmdNames...)
}

type commandStruct struct {
	Runnable
	Names
	fieldOffset int
	argCountMin int
	argCountMax int
}

func (c *commandStruct) FieldOffset() int {
	return c.fieldOffset
}

func (c *commandStruct) ArgCountMin() int {
	return c.argCountMin
}

func (c *commandStruct) ArgCountMax() int {
	return c.argCountMax
}
