package valval

type M map[string]Validator

type ObjectValidatorFunc func(content map[string]interface{}) error

type ObjectValidator interface {
	Validator
	Self(...ObjectValidatorFunc) ObjectValidator
}

type objectValidator struct {
	vMap       M
	selfVfuncs []ObjectValidatorFunc
}

func (ov *objectValidator) Validate(val interface{}) error {
	uw := unwrapPtr(val)
	if uw == nil {
		return nil
	}
	valMap, err := obj2Map(uw)
	if err != nil {
		return typeMissmatchError("object")
	}

	// inner
	if err := ov.checkInner(valMap); err != nil {
		return err
	}

	// self
	if err := ov.checkSelf(valMap); err != nil {
		return err
	}
	return nil
}

func (ov *objectValidator) checkInner(valMap map[string]interface{}) error {
	var errs []*objectErrorField
	for k, fv := range ov.vMap {
		fValue := valMap[k]
		err := fv.Validate(fValue)
		if err != nil {
			errs = append(errs, &objectErrorField{
				Name: k,
				Err:  err,
			})
		}
	}

	if errs != nil {
		ret := objectError(errs)
		return &ret
	}
	return nil
}

func (ov *objectValidator) checkSelf(valMap map[string]interface{}) error {
	errs := []error{}
	for _, svf := range ov.selfVfuncs {
		err := svf(valMap)
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

func (ov *objectValidator) Self(vfs ...ObjectValidatorFunc) ObjectValidator {
	// copy
	newOv := *ov
	newOv.selfVfuncs = vfs
	return &newOv
}

func Object(m M) ObjectValidator {
	return &objectValidator{
		vMap: m,
	}
}
