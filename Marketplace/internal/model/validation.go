package model

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func checkCategoryID(catId int) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(int)
		if s < catId {
			return errors.New("must be greater or equal 1")
		}
		return nil
	}
}
