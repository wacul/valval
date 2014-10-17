package valval

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func intptr(i int) *int {
	return &i
}

func TestSliceValidator(t *testing.T) {
	i1 := NewIntValidator(func(v int64) error {
		if v >= 0 && v <= 100 {
			return nil
		}
		return errors.New("invalid value")
	})

	Convey("Slice validator", t, func() {
		v1 := Slice(Number(i1))

		So(v1.Validate([]int{1, 2, 3}), ShouldBeNil)
		So(v1.Validate([]*int{intptr(1), intptr(2), intptr(3)}), ShouldBeNil)
		So(v1.Validate([]int{1, 2, 3, 101}), ShouldNotBeNil)
		So(v1.Validate(nil), ShouldBeNil)
		So(v1.Validate(1), ShouldNotBeNil)
		So(v1.Validate(""), ShouldNotBeNil)
		So(v1.Validate([]string{"a", "b"}), ShouldNotBeNil)

		vf := func(slice []interface{}) error {
			if slice[0].(int) != 10 || slice[1].(int) != 20 {
				return errors.New("invalid slice!")
			}
			return nil
		}
		v2 := v1.Self(vf)
		So(v1, ShouldNotEqual, v2)
		So(v2.Validate([]int{10, 20, 30}), ShouldBeNil)
		So(v2.Validate([]int{10, 10, 30}), ShouldNotBeNil)
	})
}
