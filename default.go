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

	DefaultPrefix = ""
	DefaultIndent = ""
)
