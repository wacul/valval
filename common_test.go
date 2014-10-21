package valval

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnwrapPtr(t *testing.T) {
	var nilv interface{}
	np := &nilv
	v1 := "123"
	p1 := &v1
	s1 := []string{"abc", "def"}
	sp1 := &s1
	st1 := struct {
		A string
	}{"123"}
	stp1 := &st1

	table := [][]interface{}{
		{v1, v1},
		{&v1, "123"},
		{&p1, "123"},
		{1, 1},
		{s1, s1},
		{&s1, s1},
		{&sp1, s1},
		{st1, st1},
		{&st1, st1},
		{&stp1, st1},
		{np, nil},
		{nil, nil},
	}

	for _, pair := range table {
		uw := unwrapPtr(pair[0])
		expected := pair[1]
		if !reflect.DeepEqual(uw, expected) {
			t.Errorf("unwrap %v must be %v actural %v", pair[0], pair[1], uw)
		}
	}
}

func TestObj2Map(t *testing.T) {
	Convey("obj2map", t, func() {
		s := "def"
		t1 := struct {
			A string `tag1:"123" tag2:"abc"`
			B *string
			C *string
			D int
		}{
			"abc",
			&s,
			nil,
			100,
		}
		act, err := obj2Map(t1)
		if err != nil {
			t.Error(err)
		}
		So(act["A"].value, ShouldEqual, "abc")
		So(act["A"].tag.Get("tag1"), ShouldEqual, "123")
		So(act["A"].tag.Get("tag2"), ShouldEqual, "abc")
	})
}
