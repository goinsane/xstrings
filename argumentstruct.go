package xstrings

import (
	"fmt"
	"reflect"
	"strings"
)

type ArgumentStruct struct {
	Unmarshaler              *Unmarshaler
	FieldNameBeginsLowerCase bool
	FieldNameFold            bool
	FieldTagKey              string
	FieldOffset              int
	ArgCountMin              int
	ArgCountMax              int
}

func (a *ArgumentStruct) Unmarshal(ifc interface{}, args ...string) error {
	return a.UnmarshalByValue(reflect.ValueOf(ifc), args...)
}

func (a *ArgumentStruct) UnmarshalByValue(val reflect.Value, args ...string) error {
	sizeArgs := len(args)
	if a.ArgCountMax > 0 && sizeArgs > a.ArgCountMax {
		return ErrArgumentCountExceeded
	}

	var err error
	argIdx := 0
	e := a.fieldsFunc(val, false, func(fieldName string, fieldVal reflect.Value) bool {
		fieldMinArgCount := getArgumentStructFieldMinArgCount(fieldVal.Type())
		if lastArgIdx := argIdx + fieldMinArgCount; lastArgIdx > sizeArgs {
			if argIdx < sizeArgs {
				err = &MissingArgumentError{fieldName}
				return true
			}
			if argIdx < a.ArgCountMin {
				err = &MissingArgumentError{fieldName}
				return true
			}
			if a.ArgCountMax > 0 && a.ArgCountMax <= argIdx {
				return true
			}

			fieldVal.Set(reflect.Zero(fieldVal.Type()))

			argIdx += fieldMinArgCount
			return false
		}

		var count int
		count, err = a.setFieldVal(fieldVal, fieldName, args[argIdx:]...)
		if err != nil {
			return true
		}
		argIdx += count
		return false
	})
	if e != nil {
		return e
	}
	return err
}

func (a *ArgumentStruct) Fields(ifc interface{}) (ArgumentStructFields, error) {
	return a.FieldsByValue(reflect.ValueOf(ifc))
}

func (a *ArgumentStruct) FieldsByValue(val reflect.Value) (ArgumentStructFields, error) {
	result := make(ArgumentStructFields, 0, 1024)
	argIdx := 0
	err := a.fieldsFunc(val, true, func(fieldName string, fieldVal reflect.Value) bool {
		if a.ArgCountMax > 0 && a.ArgCountMax <= argIdx {
			return true
		}
		typ := fieldVal.Type()
		typ2 := typ
		isPtr := typ2.Kind() == reflect.Ptr
		if isPtr {
			typ2 = typ2.Elem()
		}
		fieldMinArgCount := getArgumentStructFieldMinArgCount(typ2)
		result = append(result, ArgumentStructField{
			Name:        fieldName,
			Optional:    argIdx >= a.ArgCountMin,
			MinArgCount: fieldMinArgCount,
			Variadic:    typ2.Kind() == reflect.Slice,
		})
		argIdx += fieldMinArgCount
		return false
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *ArgumentStruct) GetField(ifc interface{}, name string) (interface{}, string, error) {
	fieldVal, name, err := a.GetFieldByValue(reflect.ValueOf(ifc), name)
	if err != nil {
		return nil, name, err
	}
	return fieldVal.Interface(), name, nil
}

func (a *ArgumentStruct) GetFieldByValue(val reflect.Value, name string) (reflect.Value, string, error) {
	fieldVal, name, err := a.find(val, true, name)
	if err != nil {
		return reflect.Value{}, name, err
	}

	result := reflect.New(fieldVal.Type()).Elem()
	result.Set(fieldVal)
	return result, name, nil
}

func (a *ArgumentStruct) SetField(ifc interface{}, name string, values ...string) (interface{}, string, error) {
	fieldVal, name, err := a.SetFieldByValue(reflect.ValueOf(ifc), name, values...)
	if err != nil {
		return nil, name, err
	}
	return fieldVal.Interface(), name, err
}

func (a *ArgumentStruct) SetFieldByValue(val reflect.Value, name string, values ...string) (reflect.Value, string, error) {
	fieldVal, name, err := a.find(val, false, name)
	if err != nil {
		return reflect.Value{}, name, err
	}

	_, err = a.setFieldVal(fieldVal, name, values...)
	if err != nil {
		return reflect.Value{}, name, err
	}

	result := reflect.New(fieldVal.Type()).Elem()
	result.Set(fieldVal)
	return result, name, nil
}

func (a *ArgumentStruct) setFieldVal(val reflect.Value, name string, values ...string) (count int, err error) {
	unmarshaler := a.Unmarshaler
	if unmarshaler == nil {
		unmarshaler = DefaultUnmarshaler
	}

	typ := val.Type()

	sizeValues := len(values)
	if sizeValues <= 0 {
		val.Set(reflect.Zero(typ))
		return 0, nil
	}

	typ2 := typ
	isPtr := typ2.Kind() == reflect.Ptr
	if isPtr {
		typ2 = typ2.Elem()
	}

	var av reflect.Value
	switch typ2.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		av = reflect.New(reflect.ArrayOf(sizeValues, typ2.Elem())).Elem()
		for i := 0; i < sizeValues; i++ {
			v, err := unmarshaler.ParseToValue(values[i], typ2.Elem())
			if err != nil {
				return 0, &ArgumentParseError{name, err.(*ParseError).Unwrap()}
			}
			av.Index(i).Set(v)
		}
	default:
	}

	count = getArgumentStructFieldMinArgCount(typ2)

	switch typ2.Kind() {
	case reflect.Array:
		if count > sizeValues {
			count = sizeValues
		}
		if isPtr {
			if sizeValues != typ2.Len() {
				val.Set(reflect.New(reflect.ArrayOf(typ2.Len(), typ2.Elem())))
				reflect.Copy(val.Elem(), av)
				break
			}
			val.Set(av.Addr())
			break
		}
		val.Set(reflect.Zero(typ))
		reflect.Copy(val, av)

	case reflect.Slice:
		count = sizeValues
		slc := av.Slice(0, av.Len())
		if isPtr {
			val.Set(reflect.New(reflect.SliceOf(typ2.Elem())))
			val.Elem().Set(slc)
			break
		}
		val.Set(slc)

	default:
		v, err := unmarshaler.ParseToValue(values[0], typ)
		if err != nil {
			return 0, &ArgumentParseError{name, err.(*ParseError).Unwrap()}
		}
		val.Set(v)

	}

	return count, nil
}

func (a *ArgumentStruct) find(val reflect.Value, readOnly bool, name string) (reflect.Value, string, error) {
	var result reflect.Value

	err := a.fieldsFunc(val, readOnly, func(fieldName string, fieldVal reflect.Value) bool {
		var ok bool
		if a.FieldNameFold {
			ok = strings.EqualFold(fieldName, name)
		} else {
			ok = fieldName == name
		}
		if ok {
			name = fieldName
			result = fieldVal
			return true
		}
		return false
	})
	if err != nil {
		return reflect.Value{}, name, err
	}

	if result.IsValid() {
		return result, name, nil
	}

	return reflect.Value{}, name, ErrArgumentStructFieldNotFound
}

func (a *ArgumentStruct) fieldsFunc(val reflect.Value, readOnly bool, f func(fieldName string, fieldVal reflect.Value) bool) error {
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

	if typ.Kind() != reflect.Struct {
		return ErrValueMustBeStruct
	}

	offset := a.FieldOffset
	if offset < 0 {
		offset = 0
	}

	for i, j := offset, typ.NumField(); i < j; i++ {
		sf := typ.Field(i)
		fieldVal := val.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		if sf.Anonymous && (sf.Type.Kind() == reflect.Struct ||
			(sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Struct) ||
			(sf.Type.Kind() == reflect.Interface && !fieldVal.IsNil() && fieldVal.Elem().Type().Kind() == reflect.Ptr)) {
			curFieldVal := fieldVal
			isPtrNil := sf.Type.Kind() == reflect.Ptr && fieldVal.IsNil()
			if isPtrNil {
				curFieldVal = reflect.New(sf.Type.Elem())
			}
			if sf.Type.Kind() == reflect.Interface {
				curFieldVal = fieldVal.Elem()
			}
			if err := a.fieldsFunc(curFieldVal, readOnly, f); err != nil {
				return err
			}
			if isPtrNil && !readOnly {
				fieldVal.Set(curFieldVal)
			}
			continue
		}
		fieldName := sf.Name
		if a.FieldNameBeginsLowerCase {
			fieldName = ToLowerBeginning(fieldName)
		}
		if a.FieldTagKey != "" {
			fieldTagFieldName := sf.Tag.Get(a.FieldTagKey)
			if idx := strings.Index(fieldTagFieldName, ","); idx >= 0 {
				fieldTagFieldName = fieldTagFieldName[:idx]
			}
			if fieldTagFieldName == "-" {
				continue
			}
			if fieldTagFieldName != "" {
				fieldName = fieldTagFieldName
			}
		}

		if f(fieldName, fieldVal) {
			break
		}

		if typ := fieldVal.Type(); typ.Kind() == reflect.Slice || (typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Slice) {
			break
		}
	}
	return nil
}

func getArgumentStructFieldMinArgCount(typ reflect.Type) int {
	typ2 := typ
	isPtr := typ2.Kind() == reflect.Ptr
	if isPtr {
		typ2 = typ2.Elem()
	}
	switch typ2.Kind() {
	case reflect.Array:
		return typ2.Len()
	case reflect.Slice:
		return 1
	default:
		return 1
	}
}

type ArgumentStructField struct {
	Name        string
	Optional    bool
	MinArgCount int
	Variadic    bool
}

type ArgumentStructFields []ArgumentStructField

func (a ArgumentStructFields) String() string {
	result := ""
	idx := 0

	str := ""
	for _, field := range a[idx:] {
		if field.Optional {
			break
		}
		idx++
		if str != "" {
			str += " "
		}
		vari := ""
		if field.Variadic {
			vari = "..."
		}
		if field.MinArgCount <= 1 {
			str += fmt.Sprintf("<%s>%s", field.Name, vari)
		} else {
			for i := 0; i < field.MinArgCount; i++ {
				if i > 0 {
					str += " "
				}
				str += fmt.Sprintf("<%s-%d>", field.Name, i+1)
			}
		}
	}
	result += str

	str = ""
	k := 0
	for _, field := range a[idx:] {
		if str != "" {
			str += " "
		}
		vari := ""
		if field.Variadic {
			vari = "..."
		}
		str += "["
		if field.MinArgCount <= 1 {
			str += fmt.Sprintf("<%s>%s", field.Name, vari)
		} else {
			for i := 0; i < field.MinArgCount; i++ {
				if i > 0 {
					str += " "
				}
				str += fmt.Sprintf("<%s-%d>", field.Name, i+1)
			}
		}
		k++
	}
	for i := 0; i < k; i++ {
		str += "]"
	}
	if result != "" {
		result += " "
	}
	result += str

	return result
}
