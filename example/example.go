package main

import (
	"fmt"
	"regexp"

	"github.com/wacul/valval"
)

func main() {
	personValidator := valval.Object(valval.M{
		"Name": valval.String(
			valval.MaxLength(20),
			valval.Regexp(regexp.MustCompile(`^[a-z ]+$`)),
		),
		"Attr": valval.Object(valval.M{
			"Age": valval.Number(
				valval.Min(0),
				valval.Max(120),
			),
			"Gender": valval.String(valval.In("male", "female")),
			"Tags": valval.Slice(
				valval.String(
					valval.MinLength(1),
					valval.MaxLength(10),
				),
			).Self(valval.MaxSliceLength(10)),
		}),
	})

	// Validate struct
	type PersonAttr struct {
		Age    int
		Gender string
		Tags   []string
	}
	type Person struct {
		Name string
		Attr PersonAttr
	}
	person := &Person{
		Name: "John!",
		Attr: PersonAttr{
			Age:    200,
			Gender: "otoko",
			Tags:   []string{"", "abcdefghijklmn"},
		},
	}
	if err := personValidator.Validate(person); err != nil {
		errs := valval.Errors(err)
		for _, errInfo := range errs {
			fmt.Printf("%s : %v\n", errInfo.Path, errInfo.Error)
		}
	}

	// validate map
	personMap := map[string]interface{}{
		"Name": "John!",
		"Attr": map[string]interface{}{
			"Age":    200,
			"Gender": "otoko",
			"Tags":   []string{"", "abcdefghijklmn"},
		},
	}
	if err := personValidator.Validate(personMap); err != nil {
		errs := valval.Errors(err)
		for _, errInfo := range errs {
			fmt.Printf("%s : %v\n", errInfo.Path, errInfo.Error)
		}
	}
}
