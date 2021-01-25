package xstrings

import (
	"fmt"
	"reflect"
	"strings"
)

type StructArgs struct {
	Unmarshaler  *Unmarshaler
	StructTagKey string
}

func (s *StructArgs) Unmarshal(ifc interface{}, offset, countMin, countMax int, args ...string) error {
	return s.UnmarshalByValue(reflect.ValueOf(ifc), offset, countMin, countMax, args...)
}

func (s *StructArgs) UnmarshalByValue(val reflect.Value, offset, countMin, countMax int, args ...string) error {
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

	unmarshaler := s.Unmarshaler
	if unmarshaler == nil {
		unmarshaler = DefaultUnmarshaler
	}

	sizeArgs := len(args)

	if offset < 0 {
		offset = 0
	}

	if countMax > 0 && sizeArgs > countMax {
		return ErrArgumentCountExceeded
	}

	for i, j, k := offset, 0, typ.NumField(); i < k; i++ {
		var name string
		sf := typ.Field(i)
		if s.StructTagKey != "" {
			name = sf.Tag.Get(s.StructTagKey)
			if idx := strings.Index(name, ","); idx >= 0 {
				name = name[:idx]
			}
			if name == "-" {
				continue
			}
		}
		if name == "" {
			name = ToLowerBeginning(sf.Name)
		}

		if j >= sizeArgs {
			if j < countMin {
				return &MissingArgumentError{name}
			}
			if countMax > 0 && countMax <= j {
				break
			}
			if f := val.Field(i); f.CanSet() {
				f.Set(reflect.Zero(sf.Type))
			}
			if kind := sf.Type.Kind(); kind == reflect.Ptr || kind == reflect.Array || kind == reflect.Slice {
				if kind == reflect.Array || kind == reflect.Slice {
					break
				}
				if kind := sf.Type.Elem().Kind(); kind == reflect.Array || kind == reflect.Slice {
					break
				}
			}
			continue
		}

		if f := val.Field(i); f.CanSet() {
			count, err := s.set(name, f, args[j:]...)
			if err != nil {
				return err
			}
			j += count
		}

	}

	return nil
}

func (s *StructArgs) Set(name string, args ...string) {

}

func (s *StructArgs) set(name string, val reflect.Value, values ...string) (count int, err error) {
	unmarshaler := s.Unmarshaler
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
	default:
		count = 1
	}

	switch typ2.Kind() {
	case reflect.Array:
		if sizeValues != typ2.Len() {
			return 0, fmt.Errorf("value count must be equal to %d", typ2.Len())
		}
		if isPtr {
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
