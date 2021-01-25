package xstrings

import (
	"reflect"
	"strings"
	"unicode"
)

type ArgumentStruct struct {
	Unmarshaler  *Unmarshaler
	StructTagKey string
}

func (a *ArgumentStruct) Unmarshal(ifc interface{}, offset, countMin, countMax int, args ...string) error {
	return a.UnmarshalByValue(reflect.ValueOf(ifc), offset, countMin, countMax, args...)
}

func (a *ArgumentStruct) UnmarshalByValue(val reflect.Value, offset, countMin, countMax int, args ...string) error {
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

	if offset < 0 {
		offset = 0
	}

	sizeArgs := len(args)
	if countMax > 0 && sizeArgs > countMax {
		return ErrArgumentCountExceeded
	}

	for i, j, k := offset, 0, typ.NumField(); i < k; i++ {
		sf := typ.Field(i)
		fieldName := sf.Name
		if fieldName == "" || !unicode.IsUpper([]rune(fieldName)[0]) {
			continue
		}
		fieldName = ToLowerBeginning(fieldName)
		if a.StructTagKey != "" {
			fieldName = sf.Tag.Get(a.StructTagKey)
			if idx := strings.Index(fieldName, ","); idx >= 0 {
				fieldName = fieldName[:idx]
			}
			if fieldName == "" || fieldName == "-" {
				continue
			}
		}

		if j >= sizeArgs {
			if j < countMin {
				return &MissingArgumentError{fieldName}
			}
			if countMax > 0 && countMax <= j {
				break
			}
			if f := val.Field(i); f.CanSet() {
				f.Set(reflect.Zero(sf.Type))
			}
			if kind := sf.Type.Kind(); kind == reflect.Ptr || kind == reflect.Slice {
				if kind == reflect.Array || kind == reflect.Slice {
					break
				}
				if kind := sf.Type.Elem().Kind(); kind == reflect.Slice {
					break
				}
			}
			continue
		}

		if f := val.Field(i); f.CanSet() {
			count, err := a.setFieldVal(f, fieldName, args[j:]...)
			if err != nil {
				return err
			}
			j += count
		}

	}

	return nil
}

func (a *ArgumentStruct) SetField(ifc interface{}, offset int, name string, values ...string) error {
	return a.SetFieldByValue(reflect.ValueOf(ifc), offset, name, values...)
}

func (a *ArgumentStruct) SetFieldByValue(val reflect.Value, offset int, name string, values ...string) error {
	fieldVal, err := a.find(val, offset, name)
	if err != nil {
		return err
	}

	if fieldVal.CanSet() {
		_, err = a.setFieldVal(fieldVal, name, values...)
		return err
	}

	return nil
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
				return 0, &ArgumentParseError{name, err}
			}
			av.Index(i).Set(v)
		}
		count = sizeValues
		if count > typ2.Len() {
			count = typ2.Len()
		}
	default:
		count = 1
	}

	switch typ2.Kind() {
	case reflect.Array:
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
			return 0, &ArgumentParseError{name, err}
		}
		val.Set(v)

	}

	return count, nil
}

func (a *ArgumentStruct) find(val reflect.Value, offset int, name string) (reflect.Value, error) {
	if val.Type().Kind() != reflect.Ptr {
		if !val.CanAddr() {
			return reflect.Value{}, ErrCanNotGetAddr
		}
		val = val.Addr()
	}
	if val.IsNil() {
		return reflect.Value{}, ErrNilPointer
	}

	v := val
	val = v.Elem()
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return reflect.Value{}, ErrValueMustBeStruct
	}

	if offset < 0 {
		offset = 0
	}

	for i, _, k := offset, 0, typ.NumField(); i < k; i++ {
		sf := typ.Field(i)
		fieldName := sf.Name
		if fieldName == "" || !unicode.IsUpper([]rune(fieldName)[0]) {
			continue
		}
		fieldName = ToLowerBeginning(fieldName)
		if a.StructTagKey != "" {
			fieldName = sf.Tag.Get(a.StructTagKey)
			if idx := strings.Index(fieldName, ","); idx >= 0 {
				fieldName = fieldName[:idx]
			}
			if fieldName == "" || fieldName == "-" {
				continue
			}
		}

		if fieldName == name {
			return val.Field(i), nil
		}

	}

	return reflect.Value{}, ErrArgumentStructFieldNotFound
}
