package valval

import "fmt"

type ErrorDescription struct {
	Path  string
	Error error
}

type ValueError []error

func (ve *ValueError) Error() string {
	return "invalid value"
}

type ObjectFieldError struct {
	Name string
	Err  error
}

func (ofe *ObjectFieldError) Error() string {
	return "invalid field value"
}

type ObjectError []*ObjectFieldError

func (oe *ObjectError) Error() string {
	return "invalid object"
}

type SliceElemError struct {
	Index int
	Err   error
}

func (see *SliceElemError) Error() string {
	return "invalid slice value"
}

type SliceError []*SliceElemError

func (se *SliceError) Error() string {
	return "invalid slice"
}

func Errors(err error, basePath string) []ErrorDescription {
	ret := []ErrorDescription{}
	switch t := err.(type) {
	case *ValueError:
		for _, e := range *t {
			ret = append(ret, ErrorDescription{
				Path:  basePath,
				Error: e,
			})
		}
	case *ObjectError:
		for _, ofe := range *t {
			nextBase := basePath
			if nextBase != "" {
				nextBase += "."
			}
			nextBase += ofe.Name
			ret = append(ret, Errors(ofe.Err, nextBase)...)
		}
	case *SliceError:
		for _, see := range *t {
			nextBase := basePath
			nextBase += fmt.Sprintf("[%d]", see.Index)
			ret = append(ret, Errors(see.Err, nextBase)...)
		}
	default:
		ret = append(ret, ErrorDescription{
			Path:  basePath,
			Error: t,
		})
	}

	// if rp, ok := err.(ErrorReporter); ok {
	// 	errs := rp.Errors()
	// 	ret := make([]*ErrorDescription, len(errs))
	// 	for i, e := range errs {
	// 		//copy
	// 		newE := *e
	// 		if basePath == "" {
	// 		} else {
	// 		}
	// 		ret[i] = &newE
	// 	}
	// }
	return ret
}
