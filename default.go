package xstrings

import (
	"time"
)

var (
	DefaultIntBase = 10

	DefaultTimeLayout   = "2006-01-02T15:04:05"
	DefaultTimeLocation = time.Local

	DefaultFloatFmt  = byte('f')
	DefaultFloatPrec = -1

	DefaultComplexFmt  = byte('f')
	DefaultComplexPrec = -1

	DefaultPrefix = ""
	DefaultIndent = "  "
)
