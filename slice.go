package valval

import "reflect"

type SliceValidatorFunc func(slice []interface{}) error

type SliceValidator struct {
	inner      Validator
	selfVfuncs []SliceValidatorFunc
}

func (sv *SliceValidator) Validate(slice interface{}) error {
	if slice == nil {
		return nil
	}
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

func (sv *SliceValidator) checkInner(is []interface{}) error {
	var errs []*sliceErrorElem

	for i, v := range is {
		err := sv.inner.Validate(v)
		if err != nil {
			errs = append(errs, &sliceErrorElem{
				Index: i,
				Err:   err,
			})
		}
	}
	if errs != nil {
		ret := sliceError(errs)
		return &ret
	}
	return nil
}

func (sv *SliceValidator) checkSelf(is []interface{}) error {
	errs := []error{}
	for _, svf := range sv.selfVfuncs {
		err := svf(is)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		ret := valueError(errs)
		return &ret
	}
	return nil
}

func (sv *SliceValidator) Self(vfs ...SliceValidatorFunc) *SliceValidator {
	// copy
	newSv := *sv
	newSv.selfVfuncs = vfs
	return &newSv
}

func Slice(inner Validator) *SliceValidator {
	return &SliceValidator{
		inner: inner,
	}
}
