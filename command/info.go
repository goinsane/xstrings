package command

type Info interface {
	FieldOffset() int
	ArgCountMin() int
	ArgCountMax() int
}
