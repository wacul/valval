package valval

type M map[string]Validator

type ObjectValidatorFunc func(content map[string]interface{}) error

type ObjectValidator struct {
	vMap       M
	selfVfuncs []ObjectValidatorFunc
}

func (ov *ObjectValidator) Validate(val interface{}) error {
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

func (ov *ObjectValidator) checkInner(valMap map[string]objField) error {
	var errs []*objectErrorField
	for k, fv := range ov.vMap {
		fValue := valMap[k]
		err := fv.Validate(fValue.value)
		if err != nil {
			errs = append(errs, &objectErrorField{
				Name: k,
				Err:  err,
				Tag:  fValue.tag,
			})
		}
	}

	if errs != nil {
		ret := objectError(errs)
		return &ret
	}
	return nil
}

func (ov *ObjectValidator) checkSelf(valMap map[string]objField) error {
	errs := []error{}
	for _, svf := range ov.selfVfuncs {
		err := svf(fieldMap2objMap(valMap))
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

func (ov *ObjectValidator) Self(vfs ...ObjectValidatorFunc) *ObjectValidator {
	// copy
	newOv := *ov
	newOv.selfVfuncs = vfs
	return &newOv
}

func Object(m M) *ObjectValidator {
	return &ObjectValidator{
		vMap: m,
	}
}
