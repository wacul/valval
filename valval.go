package valval

type Validator interface {
	Validate(val interface{}) error
}

type ValidatorFunc func(val interface{}) error

type ValidationError struct {
}
