package valval

import "reflect"

func NewStringValidator(inner func(string) error) ValidatorFunc {
	return func(val interface{}) error {
		if val == nil {
			return nil
		}

		rv := reflect.ValueOf(val)
		if rv.Kind() == reflect.String {
			return inner(rv.String())
		}
		return typeMissmatchError("string")
	}
}

func NewFloatValidator(inner func(float64) error) ValidatorFunc {
	return func(val interface{}) error {
		if val == nil {
			return nil
		}

		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Float32, reflect.Float64:
			return inner(rv.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return inner(float64(rv.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return inner(float64(rv.Uint()))
		}
		return typeMissmatchError("number")
	}
}

func NewIntValidator(inner func(int64) error) ValidatorFunc {
	return func(val interface{}) error {
		if val == nil {
			return nil
		}

		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return inner(rv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return inner(int64(rv.Uint()))
		}
		return typeMissmatchError("integer")
	}
}

func NewBoolValidator(inner func(bool) error) ValidatorFunc {
	return func(val interface{}) error {
		if val == nil {
			return nil
		}

		rv := reflect.ValueOf(val)
		if rv.Kind() == reflect.Bool {
			return inner(rv.Bool())
		}
		return typeMissmatchError("integer")
	}
}
