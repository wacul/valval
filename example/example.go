package main

import (
	"regexp"

	"github.com/tutuming/valval"
)

func main() {
	v := valval.Object(valval.M{
		"test1": valval.String(
			valval.MaxLength(10),
			valval.Regexp(regexp.MustCompile(`^[a-z]+$`)),
		),
		"test2": valval.Object(valval.M{
			"hoge1": valval.Number(
				valval.Min(10),
				valval.Max(100),
			),
		}),
	})
	_ = v
}
