package valval

import (
	"fmt"
	"reflect"
	"strings"
)

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
	Tag  reflect.StructTag
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

func Errors(err error) []ErrorDescription {
	return getErrors(err, "", "")
}

func ErrorsBase(err error, basePath string) []ErrorDescription {
	return getErrors(err, basePath, "")
}

func JSONErrors(err error) []ErrorDescription {
	return getErrors(err, "", "json")
}

func JSONErrorsBase(err error, basePath string) []ErrorDescription {
	return getErrors(err, basePath, "json")
}

func getErrors(err error, basePath, structTag string) []ErrorDescription {
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
			nextBase += fieldNameForTag(ofe, structTag)
			ret = append(ret, getErrors(ofe.Err, nextBase, structTag)...)
		}
	case *sliceError:
		for _, see := range *t {
			nextBase := basePath
			nextBase += fmt.Sprintf("[%d]", see.Index)
			ret = append(ret, getErrors(see.Err, nextBase, structTag)...)
		}
	default:
		ret = append(ret, ErrorDescription{
			Path:  basePath,
			Error: t,
		})
	}
	return ret
}

func fieldNameForTag(oef *objectErrorField, structTag string) string {
	if structTag == "" {
		return oef.Name
	}
	tag := oef.Tag.Get(structTag)
	if idx := strings.Index(tag, ","); idx != -1 {
		tag = tag[:idx]
	}
	if tag == "" {
		return oef.Name
	}
	return tag
}
