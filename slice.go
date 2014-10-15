package valval

import (
	"errors"
	"reflect"
)

type SliceValidator interface {
	Validator
	Self(...SliceValidatorFunc) SliceValidator
}

type SliceValidatorFunc func(slice []interface{}) error

type sliceValidator struct {
	inner      Validator
	selfVfuncs []SliceValidatorFunc
}

func (sv *sliceValidator) Validate(slice interface{}) error {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return typeMissmatchError("slice")
	}
	is := interfaceSlice(slice)

	if err := sv.checkInner(is); err != nil {
		return err
	}

	// self
	if err := sv.checkSelf(is); err != nil {
		return err
	}

	return nil
}

func (sv *sliceValidator) checkInner(is []interface{}) error {
	var errs []SliceElemError

	for i, v := range is {
		err := sv.inner.Validate(v)
		if err != nil {
			errs = append(errs, SliceElemError{
				Index: i,
				Err:   err,
			})
		}
	}
	if errs != nil {
		return SliceError(errs)
	}
	return nil
}

func (sv *sliceValidator) checkSelf(is []interface{}) error {
	errs := []error{}
	for _, svf := range sv.selfVfuncs {
		err := svf(is)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		// TODO Detail
		return errors.New("error")
	}
	return nil
}

func (sv *sliceValidator) Self(vfs ...SliceValidatorFunc) SliceValidator {
	// copy
	newSv := *sv
	newSv.selfVfuncs = vfs
	return &newSv
}

func Slice(inner Validator) SliceValidator {
	return &sliceValidator{
		inner: inner,
	}
}
