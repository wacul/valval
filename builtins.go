package valval

import (
	"errors"
	"fmt"
	"regexp"
)

func errf(str string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(str, args...))
}

func GreaterThan(min float64) ValidatorFunc {
	return NewFloatValidator(func(val float64) error {
		if val <= min {
			return errf("must be greater than %g", val)
		}
		return nil
	})
}

func Min(min float64) ValidatorFunc {
	return NewFloatValidator(func(val float64) error {
		if val < min {
			return errf("must be %g and grater", val)
		}
		return nil
	})
}

func LessThan(max float64) ValidatorFunc {
	return NewFloatValidator(func(val float64) error {
		if val >= max {
			return errf("must be less than %g", val)
		}
		return nil
	})
}

func Max(max float64) ValidatorFunc {
	return NewFloatValidator(func(val float64) error {
		if val > max {
			return errf("must be %g or less", val)
		}
		return nil
	})
}

func Regexp(r *regexp.Regexp) ValidatorFunc {
	return NewStringValidator(func(val string) error {
		if !r.MatchString(val) {
			return errf("must be match to the regexp %s", r.String())
		}
		return nil
	})
}

func In(targetValues ...interface{}) ValidatorFunc {
	return func(val interface{}) error {
		for _, tv := range targetValues {
			if val == tv {
				return nil
			}
		}
		return errf("invalid value. allowed are %v", targetValues)
	}
}

func And(funcs ...ValidatorFunc) ValidatorFunc {
	return func(val interface{}) error {
		for _, f := range funcs {
			if err := f(val); err != nil {
				return err
			}
		}
		return nil
	}
}

func Or(funcs ...ValidatorFunc) ValidatorFunc {
	return func(val interface{}) error {
		var firstError error
		for _, f := range funcs {
			if err := f(val); err == nil {
				return nil
			} else {
				if firstError == nil {
					firstError = err
				}
			}
		}
		return firstError
	}
}
