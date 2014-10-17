package main

import (
	"github.com/tutuming/valval"
)

var _ valval.Validator

func main() {
	v:= valval.Struct(valval.M{
		"test1" : valval.String(
			valval.NotNull,
			valval.MaxLength(10),
			valval.Regex(regep.MustCompile(`^[a-z]+$`)),
		),
		"test2" : valval.Struct(valval.M{
			"hoge1" : valval.Number(
				valval.NotNull,
				valval.Min(10),
				valval.Max(100),
				valval.MultiplyOf(100),
			),
		}).Self(

		)
	})
}
