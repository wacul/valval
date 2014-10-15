package valval

import (
	"reflect"
	"testing"
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
	s := "def"
	t1 := struct {
		A string
		B *string
		C *string
		D int
	}{
		"abc",
		&s,
		nil,
		100,
	}
	exp1 := map[string]interface{}{
		"A": "abc",
		"B": "def",
		"C": nil,
		"D": 100,
	}
	act1, err := obj2Map(t1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(act1, exp1) {
		t.Errorf("%v must equal %v", act1, exp1)
	}
}
