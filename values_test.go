package valval

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValuValidators(t *testing.T) {
	checkValid := func(v Validator, val interface{}, expected bool) {
		err := v.Validate(val)
		if expected == true {
			if err != nil {
				t.Errorf("validation for %v ( %v ) should be valid", v, val)
			}
		} else {
			if err == nil {
				t.Errorf("validation for %v ( %v ) should be invalid", v, val)
			}
		}
	}

	i1 := NewIntValidator(func(v int64) error {
		if v >= 0 && v <= 100 {
			return nil
		}
		return errors.New("invalid value")
	})

	s1 := NewStringValidator(func(v string) error {
		if v == "abc" {
			return nil
		}
		return errors.New("invalid value")
	})

	Convey("Values validator", t, func() {
		Convey("int", func() {
			checkValid(Number(i1), 10, true)
			checkValid(Number(i1), 1000, false)
			checkValid(Number(i1), nil, true)
			checkValid(Number(i1), "string", false)
			checkValid(Number(i1), true, false)
			checkValid(Number(i1), struct{}{}, false)
		})

		Convey("string", func() {
			checkValid(String(s1), "abc", true)
			checkValid(String(s1), "abcd", false)
			checkValid(String(s1), 123, false)
		})

		Convey("bool", func() {
			checkValid(Bool(), "abc", false)
			checkValid(Bool(), 123, false)
			checkValid(Bool(), struct{}{}, false)
			checkValid(Bool(), nil, true)
			checkValid(Bool(), true, true)
			checkValid(Bool(), false, true)
		})

		Convey("any", func() {
			checkValid(Any(), "abc", true)
			checkValid(Any(), 123, true)
			checkValid(Any(), true, true)
			checkValid(Any(), struct{}{}, true)
			checkValid(Any(), nil, true)
		})
	})
}
