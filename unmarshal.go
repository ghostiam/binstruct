package binstruct

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type unmarshal struct {
	r Reader
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "binstruct: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "binstruct: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "binstruct: Unmarshal(nil " + e.Type.String() + ")"
}

func (u *unmarshal) Unmarshal(v interface{}) error {
	return u.unmarshal(v, nil)
}

func (u *unmarshal) unmarshal(v interface{}, parentStructValues []reflect.Value) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	structValue := rv.Elem()
	numField := structValue.NumField()

	valueType := structValue.Type()
	for i := 0; i < numField; i++ {
		fieldType := valueType.Field(i)
		tags, err := parseTag(fieldType.Tag.Get(tagName))
		if err != nil {
			return fmt.Errorf(`failed parseTag for field "%s": %w`, fieldType.Name, err)
		}

		fieldData, err := parseReadDataFromTags(structValue, tags)
		if err != nil {
			return fmt.Errorf(`failed parse ReadData from tags for field "%s": %w`, fieldType.Name, err)
		}

		fieldValue := structValue.Field(i)
		err = u.setValueToField(structValue, fieldValue, fieldData, parentStructValues)
		if err != nil {
			return fmt.Errorf(`failed set value to field "%s": %w`, fieldType.Name, err)
		}
	}

	return nil
}

func (u *unmarshal) setValueToField(structValue, fieldValue reflect.Value, fieldData *fieldReadData, parentStructValues []reflect.Value) error {
	if fieldData == nil {
		fieldData = &fieldReadData{}
	}

	if fieldData.Ignore {
		return nil
	}

	r := u.r
	if fieldData.Order != nil {
		r = r.WithOrder(fieldData.Order)
	}

	err := setOffset(r, fieldData)
	if err != nil {
		return fmt.Errorf("set offset: %w", err)
	}

	if fieldData.FuncName != "" {
		var okCallFunc bool
		okCallFunc, err = callFunc(r, fieldData.FuncName, structValue, fieldValue)
		if err != nil {
			return fmt.Errorf("call custom func(%s): %w", structValue.Type().Name(), err)
		}

		if !okCallFunc {
			// Try call function from parent structs
			for i := len(parentStructValues) - 1; i >= 0; i-- {
				sv := parentStructValues[i]
				okCallFunc, err = callFunc(r, fieldData.FuncName, sv, fieldValue)
				if err != nil {
					return fmt.Errorf("call custom func from parent(%s): %w", sv.Type().Name(), err)
				}

				if okCallFunc {
					return nil
				}
			}

			message := `
failed call method, expected methods:
	func (*{{Struct}}) {{MethodName}}(r binstruct.Reader) error {} 
or
	func (*{{Struct}}) {{MethodName}}(r binstruct.Reader) ({{FieldType}}, error) {}
`
			message = strings.NewReplacer(
				`{{Struct}}`, structValue.Type().Name(),
				`{{MethodName}}`, fieldData.FuncName,
				`{{FieldType}}`, fieldValue.Type().String(),
			).Replace(message)
			return errors.New(message)
		}

		return nil
	}

	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var value int64
		var err error

		if fieldData.Length != nil {
			value, err = r.ReadIntX(int(*fieldData.Length))
		} else {
			switch fieldValue.Kind() {
			case reflect.Int8:
				v, e := r.ReadInt8()
				value = int64(v)
				err = e
			case reflect.Int16:
				v, e := r.ReadInt16()
				value = int64(v)
				err = e
			case reflect.Int32:
				v, e := r.ReadInt32()
				value = int64(v)
				err = e
			case reflect.Int64:
				v, e := r.ReadInt64()
				value = v
				err = e
			default: // reflect.Int:
				return errors.New("need set tag with len or use int8/int16/int32/int64")
			}
		}

		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetInt(value)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var value uint64
		var err error

		if fieldData.Length != nil {
			value, err = r.ReadUintX(int(*fieldData.Length))
		} else {
			switch fieldValue.Kind() {
			case reflect.Uint8:
				v, e := r.ReadUint8()
				value = uint64(v)
				err = e
			case reflect.Uint16:
				v, e := r.ReadUint16()
				value = uint64(v)
				err = e
			case reflect.Uint32:
				v, e := r.ReadUint32()
				value = uint64(v)
				err = e
			case reflect.Uint64:
				v, e := r.ReadUint64()
				value = v
				err = e
			default: // reflect.Uint:
				return errors.New("need set tag with len or use uint8/uint16/uint32/uint64")
			}
		}

		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetUint(value)
		}
	case reflect.Float32:
		f, err := r.ReadFloat32()
		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetFloat(float64(f))
		}
	case reflect.Float64:
		f, err := r.ReadFloat64()
		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetFloat(f)
		}
	case reflect.Bool:
		b, err := r.ReadBool()
		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetBool(b)
		}
	case reflect.String:
		if fieldData.Length == nil {
			return errors.New("need set tag with len for string")
		}

		_, b, err := r.ReadBytes(int(*fieldData.Length))
		if err != nil {
			return err
		}

		if fieldValue.CanSet() {
			fieldValue.SetString(string(b))
		}
	case reflect.Slice:
		if fieldData.Length == nil {
			return errors.New("need set tag with len for slice")
		}

		for i := int64(0); i < *fieldData.Length; i++ {
			tmpV := reflect.New(fieldValue.Type().Elem()).Elem()
			err = u.setValueToField(structValue, tmpV, fieldData.ElemFieldData, parentStructValues)
			if err != nil {
				return err
			}
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.Append(fieldValue, tmpV))
			}
		}
	case reflect.Array:
		var arrLen int64

		if fieldData.Length != nil {
			arrLen = *fieldData.Length
		}

		if arrLen == 0 {
			arrLen = int64(fieldValue.Len())
		}

		for i := int64(0); i < arrLen; i++ {
			tmpV := reflect.New(fieldValue.Type().Elem()).Elem()
			err = u.setValueToField(structValue, tmpV, fieldData.ElemFieldData, parentStructValues)
			if err != nil {
				return err
			}
			if fieldValue.CanSet() {
				fieldValue.Index(int(i)).Set(tmpV)
			}
		}
	case reflect.Struct:
		err = u.unmarshal(fieldValue.Addr().Interface(), append(parentStructValues, structValue))
		if err != nil {
			return fmt.Errorf("unmarshal struct: %w", err)
		}
	default:
		return errors.New(`type "` + fieldValue.Kind().String() + `" not supported`)
	}

	return nil
}

func callFunc(r Reader, funcName string, structValue, fieldValue reflect.Value) (bool, error) {
	// Call methods
	m := structValue.Addr().MethodByName(funcName)

	readerType := reflect.TypeOf((*Reader)(nil)).Elem()
	if m.IsValid() && m.Type().NumIn() == 1 && m.Type().In(0) == readerType {
		ret := m.Call([]reflect.Value{reflect.ValueOf(r)})

		errorType := reflect.TypeOf((*error)(nil)).Elem()

		// Method(r binstruct.Reader) error
		if len(ret) == 1 && ret[0].Type() == errorType {
			if !ret[0].IsNil() {
				return true, ret[0].Interface().(error)
			}

			return true, nil
		}

		// Method(r binstruct.Reader) (FieldType, error)
		if len(ret) == 2 && ret[0].Type() == fieldValue.Type() && ret[1].Type() == errorType {
			if !ret[1].IsNil() {
				return true, ret[1].Interface().(error)
			}

			if fieldValue.CanSet() {
				fieldValue.Set(ret[0])
			}
			return true, nil
		}
	}

	return false, nil
}

func setOffset(r Reader, fieldData *fieldReadData) error {
	for _, v := range fieldData.Offsets {
		_, err := r.Seek(v.Offset, v.Whence)
		if err != nil {
			return fmt.Errorf("seek: %w", err)
		}
	}

	return nil
}
