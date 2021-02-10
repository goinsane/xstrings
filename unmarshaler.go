package xstrings

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Unmarshaler struct {
	IntBase int

	TimeLayout string

	FuncParseBool     func(str string) (bool, error)
	FuncParseInt      func(str string) (int64, error)
	FuncParseUint     func(str string) (uint64, error)
	FuncParseFloat    func(str string) (float64, error)
	FuncParseComplex  func(str string) (complex128, error)
	FuncParseTime     func(str string) (time.Time, error)
	FuncParseDuration func(str string) (time.Duration, error)
	FuncUnmarshalData func(str string, ifc interface{}) error
}

func NewUnmarshaler() *Unmarshaler {
	r := initialDefaultUnmarshaler
	return &r
}

func (u *Unmarshaler) Unmarshal(str string, ifc interface{}) error {
	return u.UnmarshalByValue(str, reflect.ValueOf(ifc))
}

func (u *Unmarshaler) UnmarshalByValue(str string, val reflect.Value) (err error) {
	intBase := u.IntBase
	if intBase < 0 {
		intBase = DefaultIntBase
	}

	timeLayout := u.TimeLayout
	if timeLayout == "" {
		timeLayout = DefaultTimeLayout
	}

	if val.Type().Kind() != reflect.Ptr {
		if !val.CanAddr() {
			return ErrCanNotGetAddr
		}
		val = val.Addr()
	}
	if val.IsNil() {
		return ErrNilPointer
	}

	v := val
	val = v.Elem()
	typ := val.Type()

	if typ.Kind() == reflect.Ptr {
		if str == "" {
			val.Set(reflect.Zero(typ))
			return nil
		}
		v = reflect.New(typ.Elem())
		defer func(val reflect.Value) {
			if err != nil {
				return
			}
			val.Set(v)
		}(val)
		val = v.Elem()
		typ = val.Type()
	}

	ifc := v.Interface()

	if t, ok := ifc.(*time.Time); ok {
		var t2 time.Time
		if u.FuncParseTime != nil {
			t2, err = u.FuncParseTime(str)
		} else {
			t2, err = time.ParseInLocation(timeLayout, str, time.Local)
		}
		if err != nil {
			return newParseError(err)
		}
		*t = t2
		return nil
	}

	if t, ok := ifc.(*time.Duration); ok {
		var t2 time.Duration
		if u.FuncParseDuration != nil {
			t2, err = u.FuncParseDuration(str)
		} else {
			t2, err = time.ParseDuration(str)
		}
		if err != nil {
			return newParseError(err)
		}
		*t = t2
		return nil
	}

	if t, ok := ifc.(encoding.TextUnmarshaler); ok {
		err = t.UnmarshalText([]byte(str))
		if err != nil {
			return newParseError(err)
		}
		return nil
	}

	if t, ok := ifc.(*error); ok {
		*t = errors.New(str)
		return nil
	}

	var tryFmtScan bool
	switch typ.Kind() {
	case reflect.Bool:
		var x bool
		if u.FuncParseBool != nil {
			x, err = u.FuncParseBool(str)
		} else {
			x, err = strconv.ParseBool(str)
		}
		if err != nil {
			break
		}
		val.SetBool(x)

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		var x int64
		if u.FuncParseInt != nil {
			x, err = u.FuncParseInt(str)
		} else {
			x, err = strconv.ParseInt(str, intBase, 64)
		}
		if err != nil {
			break
		}
		val.SetInt(x)

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		var x uint64
		if u.FuncParseUint != nil {
			x, err = u.FuncParseUint(str)
		} else {
			x, err = strconv.ParseUint(str, 10, 64)
		}
		if err != nil {
			break
		}
		val.SetUint(x)

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		var x float64
		if u.FuncParseFloat != nil {
			x, err = u.FuncParseFloat(str)
		} else {
			x, err = strconv.ParseFloat(str, 64)
		}
		if err != nil {
			break
		}
		val.SetFloat(x)

	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		var x complex128
		if u.FuncParseComplex != nil {
			x, err = u.FuncParseComplex(str)
		} else {
			if parseComplex == nil {
				tryFmtScan = true
				break
			}
			x, err = parseComplex(str, 128)
		}
		if err != nil {
			break
		}
		val.SetComplex(x)

	case reflect.String:
		val.SetString(str)

	case reflect.Array:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Struct:
		if u.FuncUnmarshalData != nil {
			err = u.FuncUnmarshalData(str, ifc)
		} else {
			err = json.Unmarshal([]byte(str), ifc)
		}

	default:
		tryFmtScan = true

	}

	if tryFmtScan {
		_, err = fmt.Sscanf(str, "%v", ifc)
	}

	if err != nil {
		return newParseError(err)
	}
	return nil
}

func (u *Unmarshaler) Parse(str string, typ reflect.Type) (interface{}, error) {
	val, err := u.ParseToValue(str, typ)
	if err != nil {
		return nil, err
	}
	if typ.Kind() == reflect.Ptr && val.IsNil() {
		return nil, nil
	}
	return val.Interface(), nil
}

func (u *Unmarshaler) ParseToValue(str string, typ reflect.Type) (reflect.Value, error) {
	val := reflect.New(typ)
	err := u.UnmarshalByValue(str, val)
	if err != nil {
		return reflect.Value{}, err
	}
	return val.Elem(), nil
}

var (
	parseComplex func(s string, bitSize int) (complex128, error)
)
