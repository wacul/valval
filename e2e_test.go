package valval

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestE2E(t *testing.T) {
	type PersonAttr struct {
		IsProgrammer bool
		Height       float64
	}
	type Person struct {
		Name string
		Age  int
		Type string
		Attr *PersonAttr
	}
	Convey("E2E", t, func() {
		Convey("empty struct", func() {
			vPersonAttr := Object(M{
				"IsProgrammer": Bool(),
				"Height": Number(
					Min(100.0),
					Max(250.0),
				),
			})

			vPerson := Object(M{
				"Name": String(),
				"Attr": vPersonAttr,
			})

			p := Person{}

			So(vPerson.Validate(p), ShouldBeNil)
		})
		Convey("empty struct", func() {
			vPersonAttr := Object(M{
				"IsProgrammer": Bool(),
				"Height": Number(
					Min(100.0),
					Max(250.0),
				),
				"Foo": Slice(Number()),
			})

			vPerson := Object(M{
				"Name": String(
					MaxLength(10),
				),
				"Attr": vPersonAttr,
			})

			p := Person{}

			So(vPerson.Validate(p), ShouldBeNil)
		})
	})
}
