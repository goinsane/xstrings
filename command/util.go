package command

import (
	"strings"
	"time"

	"github.com/goinsane/xstrings"
)

var (
	unmarshaler = xstrings.NewUnmarshaler()
)

func init() {
	unmarshaler.FuncParseTime = func(str string) (time.Time, error) {
		layout := "2006-01-02T15:04"
		str = strings.ToUpper(str)
		if strings.IndexRune(str, 'T') < 0 {
			str += "T00:00"
		}
		return time.ParseInLocation(layout, str, time.Local)
	}
}

func getArgumentStruct(cmd Command) *xstrings.ArgumentStruct {
	return &xstrings.ArgumentStruct{
		Unmarshaler:              unmarshaler,
		FieldNameBeginsLowerCase: true,
		FieldNameFold:            true,
		FieldTagKey:              "field",
		FieldOffset:              cmd.FieldOffset(),
		ArgCountMin:              cmd.ArgCountMin(),
		ArgCountMax:              cmd.ArgCountMax(),
	}
}
