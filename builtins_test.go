package valval

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuiltins(t *testing.T) {
	Convey("Builtin validator funcs", t, func() {
		Convey("GreaterThan", func() {
			gt := GreaterThan(10)
			So(gt(2), ShouldNotBeNil)
			So(gt(10), ShouldNotBeNil)
			So(gt(11), ShouldBeNil)

			So(gt("abc"), ShouldNotBeNil)
			So(gt("1"), ShouldNotBeNil)
			So(gt("10"), ShouldNotBeNil)
			So(gt("11"), ShouldNotBeNil) // don't convert implicitly
		})

		Convey("Min", func() {
			min := Min(10)
			So(min(2), ShouldNotBeNil)
			So(min(10), ShouldBeNil)
			So(min(11), ShouldBeNil)

			So(min("abc"), ShouldNotBeNil)
			So(min("1"), ShouldNotBeNil)
			So(min("10"), ShouldNotBeNil)
			So(min("11"), ShouldNotBeNil) // don't convert implicitly
		})

		Convey("LessThan", func() {
			lt := LessThan(10)
			So(lt(12), ShouldNotBeNil)
			So(lt(10), ShouldNotBeNil)
			So(lt(9), ShouldBeNil)

			So(lt("abc"), ShouldNotBeNil)
			So(lt("1"), ShouldNotBeNil)
			So(lt("10"), ShouldNotBeNil)
			So(lt("11"), ShouldNotBeNil) // don't convert implicitly
		})

		Convey("Max", func() {
			max := Max(10)
			So(max(12), ShouldNotBeNil)
			So(max(10), ShouldBeNil)
			So(max(9), ShouldBeNil)

			So(max("abc"), ShouldNotBeNil)
			So(max("1"), ShouldNotBeNil)
			So(max("10"), ShouldNotBeNil)
			So(max("11"), ShouldNotBeNil) // don't convert implicitly
		})

		Convey("MinLength", func() {
			ml := MinLength(3)
			So(ml(2), ShouldNotBeNil)

			So(ml("a"), ShouldNotBeNil)
			So(ml("ab"), ShouldNotBeNil)
			So(ml("abc"), ShouldBeNil)
			So(ml("abcd"), ShouldBeNil)
		})

		Convey("MaxLength", func() {
			ml := MaxLength(3)
			So(ml(2), ShouldNotBeNil)

			So(ml("a"), ShouldBeNil)
			So(ml("ab"), ShouldBeNil)
			So(ml("abc"), ShouldBeNil)
			So(ml("abcd"), ShouldNotBeNil)
		})

		Convey("MinSliceLength", func() {
			ml := MinSliceLength(3)
			So(ml([]interface{}{}), ShouldNotBeNil)
			So(ml([]interface{}{1, 2, 3}), ShouldBeNil)
			So(ml([]interface{}{1, 2, 3, 4}), ShouldBeNil)
		})

		Convey("MaxSliceLength", func() {
			ml := MaxSliceLength(3)
			So(ml([]interface{}{}), ShouldBeNil)
			So(ml([]interface{}{1, 2, 3}), ShouldBeNil)
			So(ml([]interface{}{1, 2, 3, 4}), ShouldNotBeNil)
		})

		Convey("Regexp", func() {
			re := Regexp(regexp.MustCompile(`abc[0-9]{1,3}`))
			So(re("abc"), ShouldNotBeNil)
			So(re("abc123"), ShouldBeNil)
		})

		Convey("In", func() {
			in := In("a", 1, 2.0)
			So(in("b"), ShouldNotBeNil)
			So(in(2), ShouldNotBeNil)
			So(in(3), ShouldNotBeNil)

			So(in("a"), ShouldBeNil)
			So(in(1), ShouldBeNil)
			So(in(2.0), ShouldBeNil)
		})

		Convey("And", func() {
			v1 := Min(1)
			v2 := Max(10)
			and := And(v1, v2)

			So(and(0), ShouldNotBeNil)
			So(and(11), ShouldNotBeNil)

			So(and(1), ShouldBeNil)
			So(and(1.0), ShouldBeNil)
			So(and(10), ShouldBeNil)
			So(and(10.0), ShouldBeNil)
		})

		Convey("Or", func() {
			v1 := Regexp(regexp.MustCompile(`^abc$`))
			v2 := Regexp(regexp.MustCompile(`^def$`))
			or := Or(v1, v2)

			So(or("cde"), ShouldNotBeNil)
			So(or("abc"), ShouldBeNil)
			So(or("def"), ShouldBeNil)
		})

		Convey("RequiredFields", func() {
			req := RequiredFields("A", "B", "C")
			So(
				req(map[string]interface{}{}),
				ShouldNotBeNil,
			)
			So(
				req(map[string]interface{}{
					"A": 1,
					"B": 2,
					"C": 3,
				}),
				ShouldBeNil,
			)
			So(
				req(map[string]interface{}{
					"A": 1,
					"B": 2,
					"D": 3,
				}),
				ShouldNotBeNil,
			)
		})
	})
}
