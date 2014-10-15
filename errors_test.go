package valval

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrors(t *testing.T) {

	Convey("Errors", t, func() {
		e1 := errors.New("error1")
		e2 := errors.New("error2")
		ve1 := ValueError([]error{
			e1, e2,
		})

		Convey("Value", func() {
			ed1 := Errors(&ve1, "")

			So(len(ed1), ShouldEqual, 2)
			So(ed1[0].Path, ShouldBeEmpty)
			So(ed1[1].Path, ShouldBeEmpty)
			So(ed1[0].Error, ShouldEqual, e1)
			So(ed1[1].Error, ShouldEqual, e2)
		})

		oe1 := ObjectError([]*ObjectFieldError{
			{Name: "f1", Err: e1},
			{Name: "f2", Err: e2},
		})

		oeNexted := ObjectError([]*ObjectFieldError{
			{Name: "f1", Err: e1},
			{Name: "nested", Err: &oe1},
		})

		Convey("Object", func() {

			ed1 := Errors(&oe1, "")
			So(len(ed1), ShouldEqual, 2)
			So(ed1[0].Path, ShouldEqual, "f1")
			So(ed1[1].Path, ShouldEqual, "f2")
			So(ed1[0].Error, ShouldEqual, e1)
			So(ed1[1].Error, ShouldEqual, e2)

			ed2 := Errors(&oe1, "base")
			So(len(ed2), ShouldEqual, 2)
			So(ed2[0].Path, ShouldEqual, "base.f1")
			So(ed2[1].Path, ShouldEqual, "base.f2")
			So(ed2[0].Error, ShouldEqual, e1)
			So(ed2[1].Error, ShouldEqual, e2)

			edNested := Errors(&oeNexted, "base")
			So(len(edNested), ShouldEqual, 3)
			So(edNested[0].Path, ShouldEqual, "base.f1")
			So(edNested[1].Path, ShouldEqual, "base.nested.f1")
			So(edNested[2].Path, ShouldEqual, "base.nested.f2")
			So(edNested[0].Error, ShouldEqual, e1)
			So(edNested[1].Error, ShouldEqual, e1)
			So(edNested[2].Error, ShouldEqual, e2)
		})

		se1 := SliceError([]*SliceElemError{
			{Index: 0, Err: e1},
			{Index: 1, Err: e2},
		})

		Convey("Slice", func() {
			ed1 := Errors(&se1, "")
			So(len(ed1), ShouldEqual, 2)
			So(ed1[0].Path, ShouldEqual, "[0]")
			So(ed1[1].Path, ShouldEqual, "[1]")
			So(ed1[0].Error, ShouldEqual, e1)
			So(ed1[1].Error, ShouldEqual, e2)

			ed2 := Errors(&se1, "base")
			So(len(ed2), ShouldEqual, 2)
			So(ed2[0].Path, ShouldEqual, "base[0]")
			So(ed2[1].Path, ShouldEqual, "base[1]")
			So(ed2[0].Error, ShouldEqual, e1)
			So(ed2[1].Error, ShouldEqual, e2)
		})

		Convey("Mixed", func() {
			eMixedInnerInner := ObjectError([]*ObjectFieldError{
				{Name: "es3", Err: e1},
			})
			se := SliceError([]*SliceElemError{
				{Index: 2, Err: &eMixedInnerInner},
			})
			eMixedInner := ObjectError([]*ObjectFieldError{
				{Name: "es1", Err: e1},
				{Name: "es2", Err: &ve1},
				{Name: "nested2", Err: &se1},
				{Name: "nested3", Err: &se},
			})
			eMixed := ObjectError([]*ObjectFieldError{
				{Name: "nested", Err: &eMixedInner},
			})
			e := Errors(&eMixed, "hoge.fuga")
			So(len(e), ShouldEqual, 6)
			So(e[0].Path, ShouldEqual, "hoge.fuga.nested.es1")
			So(e[1].Path, ShouldEqual, "hoge.fuga.nested.es2")
			So(e[2].Path, ShouldEqual, "hoge.fuga.nested.es2")
			So(e[3].Path, ShouldEqual, "hoge.fuga.nested.nested2[0]")
			So(e[4].Path, ShouldEqual, "hoge.fuga.nested.nested2[1]")
			So(e[5].Path, ShouldEqual, "hoge.fuga.nested.nested3[2].es3")
		})
	})
}
