package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/weldonkipchirchir/simple_bank/util"
)

// custom validator
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		//check currency
		return util.IsSupportedCurrency(currency)
	}
	return false
}
