package xstrings

import (
	"time"
)

var (
	DefaultUnmarshaler = NewUnmarshaler()
	DefaultMarshaler   = NewMarshaler()

	DefaultIntBase = 10

	DefaultTimeLayout = time.RFC3339

	DefaultFloatFmt  = byte('f')
	DefaultFloatPrec = -1

	DefaultComplexFmt  = byte('f')
	DefaultComplexPrec = -1

	DefaultPrefix = ""
	DefaultIndent = ""
)

var (
	initialDefaultUnmarshaler = Unmarshaler{
		IntBase: -1,
	}
	initialDefaultMarshaler = Marshaler{
		IntBase:     -1,
		FloatPrec:   -2,
		ComplexPrec: -2,
	}
)
