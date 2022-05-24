package envar

import (
	"fmt"

	"github.com/ParaServices/errgo"
	"github.com/ParaServices/paratils"
)

type validationErrorMap map[string][]string

func (v validationErrorMap) HasErrors(key string) bool {
	value, ok := v[key]
	if !ok {
		return false
	}
	return len(value) > 0
}

func (v validationErrorMap) add(key, value string) {
	if !v.HasErrors(key) {
		v[key] = make([]string, 0)
	}
	v[key] = append(v[key], value)
}

func (v validationErrorMap) GetLength() int {
	return len(v)
}

// ValidatorFunc defines the signature of the function that will validate the
// v. It is expected to return an error if it fails the validation
type ValidatorFunc func(parserCtx ParserCtxAccessor, parsedField ParsedFieldGetter) error

func validateRequired(parserCtx ParserCtxAccessor, parsedField ParsedFieldGetter) error {
	if paratils.IsNil(parsedField) {
		return errgo.NewF("parsed field is nil")
	}
	if !parsedField.GetEnvFound() {
		parserCtx.AddValidationError(
			parsedField.GetStructField().Name,
			fmt.Sprintf("env key: %s not found", parsedField.GetEnvKey()),
		)
	}
	return nil
}

func validateNotEmpty(parserCtx ParserCtxAccessor, parsedField ParsedFieldGetter) error {
	if paratils.IsNil(parsedField) {
		return errgo.NewF("parsed field is nil")
	}
	if paratils.StringIsEmpty(parsedField.GetEnvValue()) {
		parserCtx.AddValidationError(
			parsedField.GetStructField().Name,
			fmt.Sprintf("env key: %s value is empty", parsedField.GetEnvKey()),
		)
	}
	return nil
}

type ValidatorFuncs []ValidatorFunc

func (v ValidatorFuncs) GetLength() int {
	return len(v)
}

type ValidatorFuncsMap map[string]ValidatorFunc

func (v ValidatorFuncsMap) GetLength() int {
	return len(v)
}

func (v ValidatorFuncsMap) Add(key string, fn ValidatorFunc) {
	v[key] = fn
}

var defaultValidatorsFunc = ValidatorFuncsMap{
	"required":  validateRequired,
	"not_empty": validateNotEmpty,
}
