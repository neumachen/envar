package envar

import (
	"reflect"

	"github.com/ParaServices/paratils"
)

type ParserCtx struct {
	tagName            string
	envVarsLoaderFunc  EnvVarsLoaderFunc
	envPrefix          string
	envPrefixDelim     string
	envSliceDelim      string
	parserFuncMap      ParserFuncMap
	validatorFuncsMap  ValidatorFuncsMap
	envVarsMap         EnvVarsMap
	validationErrorMap validationErrorMap
}

func (p *ParserCtx) GetTagName() string {
	return p.tagName
}

func (p *ParserCtx) GetEnvPrefix() string {
	return p.envPrefix
}

func (p *ParserCtx) GetEnvPrefixDelim() string {
	return p.envPrefixDelim
}

func (p *ParserCtx) GetEnvSliceDelim() string {
	return p.envSliceDelim
}

func (p *ParserCtx) GetParserFuncMap() ParserFuncMap {
	return p.parserFuncMap
}

func (p *ParserCtx) GetEnvVarsMap() EnvVarsMap {
	return p.envVarsMap
}

func (p *ParserCtx) GetValidationErrors() validationErrorMap {
	return p.validationErrorMap
}

func (p *ParserCtx) SetTagName(tagName string) error {
	p.tagName = tagName
	return nil
}

func (p *ParserCtx) SetEnvVarsLoaderFunc(envVarsLoadFunc EnvVarsLoaderFunc) error {
	p.envVarsLoaderFunc = envVarsLoadFunc
	return nil
}

func (p *ParserCtx) SetEnvPrefix(envPrefix string) error {
	p.envPrefix = envPrefix
	return nil
}

func (p *ParserCtx) SetEnvPrefixDelim(delim string) error {
	p.envPrefixDelim = delim
	return nil
}

func (p *ParserCtx) SetEnvSliceDelim(delim string) error {
	p.envSliceDelim = delim
	return nil
}

func (p *ParserCtx) SetParserFuncMap(fnMap ParserFuncMap) error {
	if fnMap.GetLength() < 1 {
		return nil
	}
	p.parserFuncMap = fnMap
	return nil
}

func (p *ParserCtx) AddParserFunc(rType reflect.Type, fn ParserFunc) {
	if paratils.IsNil(rType) {
		return
	}
	if p.parserFuncMap == nil {
		p.parserFuncMap = make(ParserFuncMap)
	}
	p.parserFuncMap.Add(rType, fn)
}

func (p *ParserCtx) SetValidatorFuncsMap(fnMap ValidatorFuncsMap) error {
	if fnMap.GetLength() < 1 {
		return nil
	}
	p.validatorFuncsMap = fnMap
	return nil
}

func (p *ParserCtx) AddValidatorFunc(validatorKey string, fn ValidatorFunc) {
	if paratils.IsNil(fn) {
		return
	}
	if p.validatorFuncsMap == nil {
		p.validatorFuncsMap = make(ValidatorFuncsMap)
	}
	p.validatorFuncsMap.Add(validatorKey, fn)
}

func (p *ParserCtx) AddValidationError(key, value string) {
	if p.validationErrorMap == nil {
		p.validationErrorMap = make(validationErrorMap)
	}
	p.validationErrorMap.add(key, value)
}

var _ ParserCtxAccessor = (*ParserCtx)(nil)

type ParserCtxGetter interface {
	// GetTagName returns the tag name used to check the tag values for a
	// given struct field
	GetTagName() string
	// GetEnvPrefix returns the ENV prefix to be appended to the ENV name
	// defined in the struct tag
	GetEnvPrefix() string
	// GetEnvPrefixDelima reeturns the ENV prefix delimiter that is used
	// to join the ENV prefix and the ENV name
	GetEnvPrefixDelim() string
	// GetEnvSliceDelim returns the delimiter used to split the values of
	// a collection for a given ENV var
	GetEnvSliceDelim() string
	// GetParserFuncMap returns ParserFuncMap that will be used to parse
	// the different types for a given struct field
	GetParserFuncMap() ParserFuncMap
	// GetEnvVarsMaps returns the environemnt varialbes that was mapped to
	// EnvVarsMap
	GetEnvVarsMap() EnvVarsMap
	// GetValidationErrors returns the validation errors that were added
	// when validation failed for a given struct field.
	GetValidationErrors() validationErrorMap
}

type ParserCtxSetter interface {
	// SetTagName sets the tagname that will be used to read the tag values
	// for a given struct field.
	SetTagName(tagName string) error
	// SetEnvVarsLoaderFunc sets the env var loader func that will be used
	// to set the EnvVarsMap taht will hold the environmental variables
	SetEnvVarsLoaderFunc(fn EnvVarsLoaderFunc) error
	// SetEnvPrefix sets the env prefix
	SetEnvPrefix(envPrefix string) error
	// SetEnvPrefixDelim sets the env prefix delimiter
	SetEnvPrefixDelim(delim string) error
	// SetEnvSliceDelima sets the delimiter for the ENV var that contains
	// a collection
	SetEnvSliceDelim(delim string) error
	// SetParserFuncMap sets the ParserFuncMap using the fnMap. This overrides the
	// ParserFuncMap for the given ParserCtx. This should only if it's not desired
	// to use the parsers already defined in this package.
	SetParserFuncMap(fnMap ParserFuncMap) error
	// AddParserFunc adds a ParserFunc. This does not reset teh ParserFuncMap for
	// the given ParserCtxSetter but will override a ParserFunc that's already
	// defined for the given rType.
	AddParserFunc(rType reflect.Type, fn ParserFunc)
	// SetValidatorFuncsMap sets the ValidatorFuncsMap for the given
	// ParserCtxSetter. This should only be used if it's not desired to use the
	// defined validators in the package and which to override ALL of them.
	SetValidatorFuncsMap(fnMap ValidatorFuncsMap) error
	// AddValidatorFunc adds a validator to the givne ParserCtxSetter's
	// ValidatorFuncsMap. This does not reset the ValidatorFuncsMap but will
	// override any existing validator that matches the validatorKey.
	AddValidatorFunc(validatorKey string, fn ValidatorFunc)
	// AddValidationError adds the validation error that was generated
	// when a validation failed.
	AddValidationError(key, value string)
}

type ParserCtxAccessor interface {
	ParserCtxGetter
	ParserCtxSetter
}

type ParserCtxFuncSetter func(setter ParserCtxSetter) error

// SetTagName set the tag name that will be used to for the struct field
func SetTagName(tagName string) ParserCtxFuncSetter {
	return func(setter ParserCtxSetter) error {
		return setter.SetTagName(tagName)
	}
}

// SetEnvVarsLoaderFunc sets the func that is used ot load the env vars in to
// the EnvVarsMap
func SetEnvVarsLoaderFunc(fn EnvVarsLoaderFunc) ParserCtxFuncSetter {
	return func(setter ParserCtxSetter) error {
		return setter.SetEnvVarsLoaderFunc(fn)
	}
}

// SetEnvPrefix sets the prefix that will be appending to the env var name in
// the struct field. This is useful when an existing struct is reused but the
// env var defined in the struct field ha a prefix.
func SetEnvPrefix(envPrefix string) ParserCtxFuncSetter {
	return func(setter ParserCtxSetter) error {
		return setter.SetEnvPrefix(envPrefix)
	}
}

// SetEnvPrefixDelim sets the delimiter used when appending a prefix to the
// existing env vars defined in a struct field. Defaault is DefaultEnvPrefixDelim
func SetEnvPrefixDelim(delim string) ParserCtxFuncSetter {
	return func(setter ParserCtxSetter) error {
		return setter.SetEnvPrefixDelim(delim)
	}
}

// SetEnvSliceDelim defines the delimiter used to split a collection value in
// a given ENV var value.
func SetEnvSliceDelim(delim string) ParserCtxFuncSetter {
	return func(setter ParserCtxSetter) error {
		return setter.SetEnvSliceDelim(delim)
	}
}

// SetParserFuncMap sets the ParserFuncMap using the fnMap. This overrides the
// ParserFuncMap for the given ParserCtx. This should only if it's not desired
// to use the parsers already defined in this package.
func SetParserFuncMap(fnMap ParserFuncMap) ParserCtxFuncSetter {
	return func(setter ParserCtxSetter) error {
		return setter.SetParserFuncMap(fnMap)
	}
}

// AddParserFunc adds a ParserFunc. This does not reset teh ParserFuncMap for
// the given ParserCtxSetter but will override a ParserFunc that's already
// defined for the given rType.
func AddParserFunc(rType reflect.Type, fn ParserFunc) ParserCtxFuncSetter {
	return func(setter ParserCtxSetter) error {
		setter.AddParserFunc(rType, fn)
		return nil
	}
}

// SetValidatorFuncsMap sets the ValidatorFuncsMap for the given
// ParserCtxSetter. This should only be used if it's not desired to use the
// defined validators in the package and which to override ALL of them.
func SetValidatorFuncsMap(fnMap ValidatorFuncsMap) ParserCtxFuncSetter {
	return func(setter ParserCtxSetter) error {
		return setter.SetValidatorFuncsMap(fnMap)
	}
}

// AddValidatorFunc adds a validator to the givne ParserCtxSetter's
// ValidatorFuncsMap. This does not reset the ValidatorFuncsMap but will
// override any existing validator that matches the validatorKey.
func AddValidatorFunc(validatorKey string, fn ValidatorFunc) ParserCtxFuncSetter {
	return func(setter ParserCtxSetter) error {
		setter.AddValidatorFunc(validatorKey, fn)
		return nil
	}
}

var defaultParserCtxSetters = []ParserCtxFuncSetter{
	SetTagName(DefaultTagName),
	SetEnvPrefixDelim(DefaultEnvPrefixDelim),
	SetEnvSliceDelim(DefaultEnvSliceDelim),
	SetEnvVarsLoaderFunc(loadEnvVarsToMap),
	SetParserFuncMap(defaultParserFuncs()),
	SetValidatorFuncsMap(defaultValidatorsFunc),
}
