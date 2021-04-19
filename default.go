package xstrings

import (
	"time"
)

var (
	DefaultIntBase = 10

	DefaultTimeLayout = time.RFC3339

	DefaultFloatFmt  = byte('f')
	DefaultFloatPrec = -1

	DefaultComplexFmt  = byte('f')
	DefaultComplexPrec = -1

	DefaultIndent     = ""
	DefaultLinePrefix = ""
)

var (
	initialUnmarshaler = Unmarshaler{
		IntBase: -1,
	}
	initialMarshaler = Marshaler{
		IntBase:     -1,
		FloatPrec:   -2,
		ComplexPrec: -2,
	}
)
