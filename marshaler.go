package xstrings

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Marshaler struct {
	IntBase int

	TimeLayout string

	FloatFmt  byte
	FloatPrec int

	ComplexFmt  byte
	ComplexPrec int

	Indent          string
	MultiLinePrefix string

	FuncFormatBool     func(v bool) string
	FuncFormatInt      func(v int64) string
	FuncFormatUint     func(v uint64) string
	FuncFormatFloat    func(v float64) string
	FuncFormatComplex  func(v complex128) string
	FuncFormatTime     func(v time.Time) string
	FuncFormatDuration func(v time.Duration) string
	FuncMarshalData    func(v interface{}) (string, error)
}

func NewMarshaler() *Marshaler {
	r := initialMarshaler
	return &r
}

func (m *Marshaler) Marshal(ifc interface{}) (string, error) {
	return m.MarshalByValue(reflect.ValueOf(ifc))
}

func (m *Marshaler) MarshalByValue(val reflect.Value) (string, error) {
	var err error
	var str string

	intBase := m.IntBase
	if intBase < 0 {
		intBase = DefaultIntBase
	}

	timeLayout := m.TimeLayout
	if timeLayout == "" {
		timeLayout = DefaultTimeLayout
	}

	floatFmt := m.FloatFmt
	if floatFmt == 0 {
		floatFmt = DefaultFloatFmt
	}

	floatPrec := m.FloatPrec
	if floatPrec < -1 {
		floatPrec = DefaultFloatPrec
	}

	complexFmt := m.ComplexFmt
	if complexFmt == 0 {
		complexFmt = DefaultComplexFmt
	}

	complexPrec := m.ComplexPrec
	if complexPrec < -1 {
		complexPrec = DefaultComplexPrec
	}

	indent := m.Indent
	if indent == "" {
		indent = DefaultIndent
	}

	multiLinePrefix := m.MultiLinePrefix
	if multiLinePrefix == "" {
		multiLinePrefix = DefaultMultiLinePrefix
	}

	typ := val.Type()
	ifcOrig := val.Interface()
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", nil
		}
		val = val.Elem()
		typ = typ.Elem()
	}
	ifc := val.Interface()

	if t, ok := ifc.(time.Time); ok {
		if m.FuncFormatTime != nil {
			str = m.FuncFormatTime(t)
		} else {
			str = t.Format(timeLayout)
		}
		return str, nil
	}

	if t, ok := ifc.(time.Duration); ok {
		if m.FuncFormatDuration != nil {
			str = m.FuncFormatDuration(t)
		} else {
			str = t.String()
		}
		return str, nil
	}

	if t, ok := ifcOrig.(encoding.TextMarshaler); ok {
		var data []byte
		data, err = t.MarshalText()
		if err != nil {
			return "", newFormatError(err)
		}
		return string(data), nil
	}

	if t, ok := ifcOrig.(error); ok {
		return t.Error(), nil
	}

	/*if t, ok := ifcOrig.(fmt.Stringer); ok {
		return t.String(), nil
	}*/

	tryFmtPrint := false
	var boolVal bool
	var intVal int64
	var uintVal uint64
	var floatVal float64
	var complexVal complex128
	var stringVal string
	var dataVal interface{}
	switch typ.Kind() {
	case reflect.Bool:
		boolVal = ifc.(bool)
	case reflect.Int:
		intVal = int64(ifc.(int))
	case reflect.Int8:
		intVal = int64(ifc.(int8))
	case reflect.Int16:
		intVal = int64(ifc.(int16))
	case reflect.Int32:
		intVal = int64(ifc.(int32))
	case reflect.Int64:
		intVal = ifc.(int64)
	case reflect.Uint:
		uintVal = uint64(ifc.(uint))
	case reflect.Uint8:
		uintVal = uint64(ifc.(uint8))
	case reflect.Uint16:
		uintVal = uint64(ifc.(uint16))
	case reflect.Uint32:
		uintVal = uint64(ifc.(uint32))
	case reflect.Uint64:
		uintVal = ifc.(uint64)
	case reflect.Uintptr:
		uintVal = uint64(ifc.(uintptr))
	case reflect.Float32:
		floatVal = float64(ifc.(float32))
	case reflect.Float64:
		floatVal = ifc.(float64)
	case reflect.Complex64:
		complexVal = complex128(ifc.(complex64))
	case reflect.Complex128:
		complexVal = ifc.(complex128)
	case reflect.String:
		stringVal = ifc.(string)

	case reflect.Array:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Struct:
		dataVal = ifcOrig

	default:
		tryFmtPrint = true
	}

	kind := typ.Kind()
	if tryFmtPrint {
		kind = reflect.Invalid
	}
	switch kind {
	case reflect.Bool:
		if m.FuncFormatBool != nil {
			str = m.FuncFormatBool(boolVal)
		} else {
			str = strconv.FormatBool(boolVal)
		}

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		if m.FuncFormatInt != nil {
			str = m.FuncFormatInt(intVal)
		} else {
			str = strconv.FormatInt(intVal, intBase)
		}

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
		if m.FuncFormatUint != nil {
			str = m.FuncFormatUint(uintVal)
		} else {
			str = strconv.FormatUint(uintVal, intBase)
		}

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		if m.FuncFormatFloat != nil {
			str = m.FuncFormatFloat(floatVal)
		} else {
			str = strconv.FormatFloat(floatVal, floatFmt, floatPrec, 64)
		}

	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		if m.FuncFormatComplex != nil {
			str = m.FuncFormatComplex(complexVal)
		} else {
			if formatComplex == nil {
				tryFmtPrint = true
				break
			}
			str = formatComplex(complexVal, complexFmt, complexPrec, 128)
		}

	case reflect.String:
		str = stringVal

	case reflect.Array:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Struct:
		if m.FuncMarshalData != nil {
			str, err = m.FuncMarshalData(dataVal)
		} else {
			var data []byte
			data, err = json.Marshal(dataVal)
			if err != nil {
				break
			}
			if indent != "" {
				buf := bytes.NewBuffer(make([]byte, 0, 4*len(data)))
				err = json.Indent(buf, data, "", indent)
				if err != nil {
					break
				}
				data = buf.Bytes()
			}
			str = string(data)
		}

	case reflect.Invalid:
		fallthrough

	default:
		tryFmtPrint = true

	}

	if tryFmtPrint {
		str = fmt.Sprintf("%v", ifc)
		err = nil
	}

	if err != nil {
		return "", newFormatError(err)
	}

	newStr := ""
	for idx, line := range strings.Split(str, "\n") {
		nl := ""
		prefix := ""
		if idx > 0 {
			nl = "\n"
			prefix = multiLinePrefix
		}
		newStr += nl + prefix + line
	}

	return newStr, nil
}

var (
	formatComplex func(c complex128, fmt byte, prec, bitSize int) string
)
