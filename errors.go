package valval

type PathError struct {
	Path string
	Err  error
}

type ErrorReporter interface {
	Errors() []PathError
}

type ValueError []error

func (ve *ValueError) Errors() []PathError {
	ret := make([]PathError, len(*ve))
	for i, e := range *ve {
		ret[i].Err = e
	}
	return ret
}

func (ve ValueError) Error() string {
	return "invalid value"
}

type ObjectFieldError struct {
	Name string
	Err  error
}

func (ofe ObjectFieldError) Error() string {
	return "invalid field value"
}

type ObjectError []ObjectFieldError

func (oe ObjectError) Error() string {
	return "invalid object"
}

type SliceElemError struct {
	Index int
	Err   error
}

func (see SliceElemError) Error() string {
	return "invalid slice value"
}

type SliceError []SliceElemError

func (se SliceError) Error() string {
	return "invalid slice"
}
