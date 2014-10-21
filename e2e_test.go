package valval

import (
	"sort"
	"testing"

	"bytes"

	. "github.com/smartystreets/goconvey/convey"
)

type errorsByKey []ErrorDescription

func (errs errorsByKey) Len() int      { return len(errs) }
func (errs errorsByKey) Swap(i, j int) { errs[i], errs[j] = errs[j], errs[i] }
func (errs errorsByKey) Less(i, j int) bool {
	return bytes.Compare(
		[]byte(errs[i].Path),
		[]byte(errs[j].Path),
	) < 0
}

func TestE2E(t *testing.T) {
	type PersonAttr struct {
		IsProgrammer bool
		Height       float64
	}
	type Person struct {
		Name   string
		Age    int
		Gender string
		Attr   *PersonAttr
		Tags   []string
	}
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
		"Age":    Number(Min(0), Max(120)),
		"Gender": String(In("male", "female")),
		"Attr":   vPersonAttr,
		"Tags": Slice(String(
			MinLength(1),
			MaxLength(10),
		)),
	})

	Convey("E2E", t, func() {
		Convey("valid struct", func() {
			p := Person{
				Name:   "John",
				Age:    10,
				Gender: "male",
				Attr: &PersonAttr{
					IsProgrammer: false,
					Height:       170.0,
				},
				Tags: []string{"good", "great"},
			}
			So(vPerson.Validate(p), ShouldBeNil)
		})

		Convey("invalid struct", func() {
			p := Person{
				Name:   "John abcdefghijklmn",
				Age:    -10,
				Gender: "male2",
				Attr: &PersonAttr{
					IsProgrammer: false,
					Height:       500.0,
				},
				Tags: []string{"good", "great", "", "greatgreatgreatgreatgreatgreatgreat"},
			}
			err := vPerson.Validate(p)
			So(err, ShouldNotBeNil)
			errs := Errors(err)
			sort.Sort(errorsByKey(errs))

			So(errs[0].Path, ShouldEqual, "Age")
			So(errs[1].Path, ShouldEqual, "Attr.Height")
			So(errs[2].Path, ShouldEqual, "Gender")
			So(errs[3].Path, ShouldEqual, "Name")
			So(errs[4].Path, ShouldEqual, "Tags[2]")
			So(errs[5].Path, ShouldEqual, "Tags[3]")
		})
	})
}
