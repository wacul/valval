package valval

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestObjectValidator(t *testing.T) {
	type t1 struct {
		A string
		B int
		C bool
		D *string
		E *int
		F *bool
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

	Convey("Object validator", t, func() {
		v1 := Object(M{
			"A": String(s1),
			"B": Number(i1),
			"C": Bool(),
			"D": String(s1),
			"E": Number(i1),
			"F": Bool(),
		})
		st1 := t1{
			A: "abc",
			B: 1,
			C: false,
		}
		So(v1.Validate(st1), ShouldBeNil)

		st2 := t1{
			A: "abc",
			B: 1,
			C: false,
			D: &st1.A,
			E: &st1.B,
			F: &st1.C,
		}
		So(v1.Validate(st2), ShouldBeNil)

		st3 := t1{
			A: "abcd",
			B: 1,
			C: false,
			D: &st1.A,
			E: &st1.B,
			F: &st1.C,
		}

		Convey("object", func() {
			So(v1.Validate(1), ShouldNotBeNil)
			So(v1.Validate(nil), ShouldBeNil)
			So(v1.Validate("aa"), ShouldNotBeNil)
			So(v1.Validate(true), ShouldNotBeNil)
			So(v1.Validate(st3), ShouldNotBeNil)
		})

		Convey("map", func() {
			m1 := map[string]interface{}{
				"A": "abc",
				"B": 1,
				"C": false,
			}
			So(v1.Validate(m1), ShouldBeNil)

			m2 := map[string]interface{}{
				"A": "abc",
				"B": 1,
				"C": false,
				"D": &st1.A,
				"E": &st1.B,
				"F": &st1.C,
			}
			So(v1.Validate(m2), ShouldBeNil)

			v2 := v1.Self(func(content map[string]interface{}) error {
				if content["D"] == nil && content["E"] == nil {
					return errors.New("D and E needed")
				}
				return nil
			})
			So(v1, ShouldNotEqual, v2)

			st4 := t1{
				A: "abc",
				B: 1,
				C: false,
				D: nil,
				E: &st1.B,
				F: &st1.C,
			}
			So(v2.Validate(st4), ShouldBeNil)

			st5 := t1{
				A: "abc",
				B: 1,
				C: false,
				D: nil,
				E: nil,
				F: &st1.C,
			}
			So(v2.Validate(st5), ShouldNotBeNil)

		})

		Convey("Pointer nil", func() {
			type t2 struct {
				P *t1
			}
			v := Object(M{
				"P": Object(M{}),
			})
			o := t1{}
			So(v.Validate(o), ShouldBeNil)
		})
	})

}
