# Valval [ ![Codeship Status for wcl48/valval](https://www.codeship.io/projects/59f9c260-37f9-0132-e8f4-461894a5379e/status)](https://www.codeship.io/projects/41854) [![Coverage Status](https://coveralls.io/repos/wcl48/valval/badge.png?branch=master)](https://coveralls.io/r/wcl48/valval?branch=master)

Valval is simple validation library for Go.

Defines a validator

```go
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
```

validate a struct

```go
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
/* outputs
Attr.Age : must be 200 or less
Attr.Gender : invalid value. allowed are [male female]
Attr.Tags[0] : length must be 1 or greater
Attr.Tags[1] : length must be 10 or less
Name : must be match to the regexp ^[a-z ]+$
*/

```

and validate a map

```go
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
// outputs same as struct
```

## Install

```
go get github.com/wcl48/valval
```
