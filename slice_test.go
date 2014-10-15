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
		So(v1.Validate(1), ShouldNotBeNil)
		So(v1.Validate(nil), ShouldNotBeNil)
		So(v1.Validate(""), ShouldNotBeNil)
		So(v1.Validate([]string{"a", "b"}), ShouldNotBeNil)
	})
}
