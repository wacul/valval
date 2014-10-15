package valval

import (
	"errors"
	"reflect"
)

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

func validateByFuncs(vfs []ValidatorFunc, val interface{}) error {
	errs := []error{}
	for _, vf := range vfs {
		err := vf(val)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		// TODO Detail
		return errors.New("error")
	}
	return nil
}

func obj2Map(val interface{}) (map[string]interface{}, error) {
	if m, ok := val.(map[string]interface{}); ok {
		return flattenMap(m), nil
	}

	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Struct {
		return struct2Map(val), nil
	}

	return nil, errors.New("invalid type")
}

func flattenMap(m map[string]interface{}) map[string]interface{} {
	ret := map[string]interface{}{}
	for k, v := range m {
		ret[k] = unwrapPtr(v)
	}
	return ret
}

func struct2Map(val interface{}) map[string]interface{} {
	rv := reflect.ValueOf(val)
	sv := reflect.TypeOf(val)
	ret := map[string]interface{}{}
	for i := 0; i < rv.NumField(); i++ {
		f := sv.Field(i)
		fv := rv.Field(i)
		ret[f.Name] = unwrapPtr(fv.Interface())
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
