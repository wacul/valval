package valval

import "fmt"

type ErrorDescription struct {
	Path  string
	Error error
}

type valueError []error

func (ve *valueError) Error() string {
	return "invalid value"
}

type objectErrorField struct {
	Name string
	Err  error
}

type objectError []*objectErrorField

func (oe *objectError) Error() string {
	return "invalid object"
}

type sliceErrorElem struct {
	Index int
	Err   error
}

type sliceError []*sliceErrorElem

func (se *sliceError) Error() string {
	return "invalid slice"
}

func Errors(err error, basePath string) []ErrorDescription {
	ret := []ErrorDescription{}
	if err == nil {
		return ret
	}
	switch t := err.(type) {
	case *valueError:
		for _, e := range *t {
			ret = append(ret, ErrorDescription{
				Path:  basePath,
				Error: e,
			})
		}
	case *objectError:
		for _, ofe := range *t {
			nextBase := basePath
			if nextBase != "" {
				nextBase += "."
			}
			nextBase += ofe.Name
			ret = append(ret, Errors(ofe.Err, nextBase)...)
		}
	case *sliceError:
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
	return ret
}
