package valval

import (
	"errors"
	"reflect"
)

type Validator interface {
	Validate(val interface{}) error
}

type objField struct {
	value interface{}
	tag   reflect.StructTag
}

type ValidatorFunc func(val interface{}) error

func isNilValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr,
		reflect.Interface, reflect.Slice:
		if v.IsNil() {
			return true
		}
	}
	return false
}

func unwrapPtr(val interface{}) interface{} {
	if val == nil {
		return nil
	}
	v := reflect.ValueOf(val)

	if isNilValue(v) {
		return nil
	}

	if v.Kind() == reflect.Ptr {
		elm := v.Elem()
		return unwrapPtr(elm.Interface())
	}
	return val
}

func obj2Map(val interface{}) (map[string]objField, error) {
	if val == nil {
		return nil, nil
	}
	if m, ok := val.(map[string]interface{}); ok {
		return flattenMap(m), nil
	}
	uv := unwrapPtr(val)
	rv := reflect.ValueOf(uv)
	if rv.Kind() == reflect.Struct {
		return struct2Map(val), nil
	}

	return nil, errors.New("invalid type")
}

func fieldMap2objMap(in map[string]objField) map[string]interface{} {
	ret := map[string]interface{}{}
	for k, v := range in {
		ret[k] = v.value
	}
	return ret
}

func flattenMap(m map[string]interface{}) map[string]objField {
	ret := map[string]objField{}
	for k, v := range m {
		ret[k] = objField{
			value: unwrapPtr(v),
		}
	}
	return ret
}

func struct2Map(val interface{}) map[string]objField {
	rv := reflect.ValueOf(val)
	sv := reflect.TypeOf(val)
	ret := map[string]objField{}
	for i := 0; i < rv.NumField(); i++ {
		f := sv.Field(i)
		fv := rv.Field(i)
		ret[f.Name] = objField{
			value: unwrapPtr(fv.Interface()),
			tag:   f.Tag,
		}
	}

	return ret
}

func interfaceSlice(s interface{}) []interface{} {
	rv := reflect.ValueOf(s)
	length := rv.Len()
	ret := make([]interface{}, length)
	for i := 0; i < length; i++ {
		v := rv.Index(i)
		if isNilValue(v) {
			ret[i] = nil
			continue
		}

		ret[i] = unwrapPtr(v.Interface())
	}
	return ret
}

func typeMissmatchError(ts string) error {
	return errors.New("type missmatch " + ts)
}
