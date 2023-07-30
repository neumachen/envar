package envar

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/neumachen/errorx"
	"github.com/neumachen/gobag"
)

func newParsedField(parserCtx *ParserCtx, structField reflect.StructField) (*parsedField, error) {
	pField := &parsedField{
		structField:    structField,
		tagName:        parserCtx.GetTagName(),
		envPrefix:      parserCtx.GetEnvPrefix(),
		envPrefixDelim: parserCtx.GetEnvPrefixDelim(),
	}

	if err := parseStructField(pField); err != nil {
		return nil, errorx.New(err)
	}

	if !pField.tagFound {
		return nil, nil
	}
	return pField, nil
}

type parsedField struct {
	tagOpts        tagOpts
	structField    reflect.StructField
	envPrefix      string
	envPrefixDelim string
	tagName        string
	envName        string
	envValue       string
	tagFound       bool
	keyFound       bool
}

func (p *parsedField) GetStructField() reflect.StructField {
	return p.structField
}

// GetEnvValue is the value when the env vars was parsed using the GetEnvKey
func (p *parsedField) GetEnvValue() string {
	return p.envValue
}

// GetDefaultValue is the value defined in the struct tag key default
func (p *parsedField) GetDefaultValue() string {
	return p.tagOpts.getDefaultValue()
}

// getFieldValue returns the default value if GetEnvValue returns an empty
// string ("even if the value is blank")
func (p *parsedField) getFieldValue() string {
	if v := p.GetEnvValue(); !gobag.StringIsEmpty(v) {
		return v
	}
	return p.GetDefaultValue()
}

// GetEnvFound represents whether the environment variable was found
func (p *parsedField) GetEnvFound() bool {
	return p.keyFound
}

// GetEnvPrefix that will be appending to the env key when looking up the
// environment variables
func (p *parsedField) GetEnvPrefix() string {
	return p.envPrefix
}

// GetEnvKey joins the EnvPrefix and EnvName
func (p *parsedField) GetEnvKey() string {
	if gobag.StringIsEmpty(p.GetEnvPrefix()) {
		return p.GetEnvName()
	}
	return strings.Join([]string{p.GetEnvPrefix(), p.GetEnvName()}, p.envPrefixDelim)
}

// GetEnvName returns the value of the tag, e.g env:"FOO", the value will be
// FOO
func (p *parsedField) GetEnvName() string {
	return p.envName
}

// GetTagName returns the tag key used when parsing the struct field. The
// default is defaultTagName (env)
func (p *parsedField) GetTagName() string {
	return p.tagName
}

func (p *parsedField) isNested() bool {
	return p.tagOpts.getNested()
}

func (p *parsedField) unsetEnv() bool {
	return p.tagOpts.getUnsetKey()
}

func (p *parsedField) validate(parserCtx *ParserCtx) error {
	if len(p.tagOpts.getValidate()) < 1 {
		return nil
	}

	for i := range p.tagOpts.getValidate() {
		vFunc, ok := parserCtx.validatorFuncsMap[p.tagOpts.getValidate()[i]]
		if !ok {
			return errorx.New(fmt.Sprintf("validator func: %s not found", p.tagOpts.getValidate()[i]))
		}
		if err := vFunc(parserCtx, p); err != nil {
			return errorx.New(err)
		}
	}

	return nil
}

func (p *parsedField) setTagOpts(key, value string) {
	if p.tagOpts == nil {
		p.tagOpts = make(tagOpts)
	}
	p.tagOpts[key] = value
}

func (p *parsedField) setEnvName(tagValue string) error {
	p.envName = tagValue
	return nil
}

func (p *parsedField) setEnvValue(eMap EnvVarsMap) error {
	if eMap.GetLength() < 1 {
		return nil
	}
	p.envValue, p.keyFound = eMap.Get(p.GetEnvKey())
	return nil
}

func (p *parsedField) setField(parserCtx *ParserCtx, fieldValue reflect.Value) error {
	if v := p.getFieldValue(); gobag.StringIsEmpty(v) {
		return nil
	}

	if fieldValue.Kind() == reflect.Slice {
		return handleSlice(parserCtx, p, fieldValue)
	}

	if unmarshaler := asTextUnmarshaler(fieldValue); unmarshaler != nil {
		if err := unmarshaler.UnmarshalText([]byte(p.getFieldValue())); err != nil {
			return newParseError(p.GetStructField(), err)
		}
		return nil
	}

	sType := p.GetStructField().Type

	if p.GetStructField().Type.Kind() == reflect.Ptr {
		sType = sType.Elem()
		fieldValue = fieldValue.Elem()
	}

	parserFunc := parserCtx.GetParserFuncMap().Get(p.GetStructField().Type)
	if !gobag.IsNil(parserFunc) {
		val, err := parserFunc(p.getFieldValue())
		if err != nil {
			return newParseError(p.GetStructField(), err)
		}

		fieldValue.Set(reflect.ValueOf(val).Convert(sType))
		return nil
	}

	return newNoParserError(p.GetStructField())
}

var _ ParsedFieldGetter = (*parsedField)(nil)

type ParsedFieldGetter interface {
	isNested() bool
	unsetEnv() bool
	GetStructField() reflect.StructField
	GetEnvValue() string
	GetEnvFound() bool
	GetEnvKey() string
}

const DefaultTagName = "env"
const DefaultEnvPrefixDelim = "_"
const DefaultEnvSliceDelim = ","

const tagOptsDelim = ","
const tagOptsNested = "nested"
const tagOptsValidateKey = "validate"

// TODO: implement this later
// const tagOptsExpandKey = "expand"
const tagOptsUnsetKey = "unset"
const tagOptsDefaultKey = "default"

const validateDelim = "|"
const defaultDelim = "|"

var trueStrs = []string{
	"true",
	"t",
	"1",
}

type tagOpts map[string]string

func (t tagOpts) getLength() int {
	return len(t)
}

func (t tagOpts) getNested() bool {
	_, ok := t[tagOptsNested]
	// we don't care what the value is as long as the key is there, it's
	// expected to be a nested struct
	return ok
}

func (t tagOpts) getValidate() []string {
	v, ok := t[tagOptsValidateKey]
	if !ok {
		return nil
	}
	return strings.Split(v, validateDelim)
}

func (t tagOpts) getUnsetKey() bool {
	v, ok := t[tagOptsUnsetKey]
	if !ok {
		return false
	}
	return gobag.ArrayContainsStr(trueStrs, v)
}

// func (t tagOpts) getExpandKey() bool {
// 	v, ok := t[tagOptsExpandKey]
// 	if !ok {
// 		return false
// 	}
// 	return gobag.ArrayContainsStr(trueStrs, v)
// }

func (t tagOpts) getDefaultValue() string {
	v, ok := t[tagOptsDefaultKey]
	if !ok {
		return ""
	}
	return v
}

func removeEmptyString(slice *[]string) {
	i := 0
	p := *slice
	for _, entry := range p {
		if !gobag.StringIsEmpty(entry) {
			p[i] = entry
			i++
		}
	}
	*slice = p[0:i]
}

type tagValues []string

func (t tagValues) getLength() int {
	return len(t)
}

func (t tagValues) isNested() bool {
	for i := range t {
		if t[i] == tagOptsNested {
			return true
		}
	}
	return false
}

// getTagValues returns the values for the given tag, e.g, `env:"FOO" wil
// lreturn tagValues{"FOO"}. If the tag is present but no env name is given,
// e.g, `env:",nested" the tagValues is not removed of empty strings. The
// reason behind this is that there are times when we want to process a field
// that is a struct.
func getTagValues(p *parsedField) tagValues {
	tagVals := strings.Split(p.GetStructField().Tag.Get(p.GetTagName()), tagOptsDelim)
	return tagValues(tagVals)
}

func parseStructField(p *parsedField) error {
	tagVals := getTagValues(p)
	if p.tagFound = tagVals.getLength() > 0; !p.tagFound {
		return nil
	}
	p.setEnvName(tagVals[0])

	// if it's less than two, there are no options
	if tagVals.getLength() < 2 {
		return nil
	}

	for i := range tagVals[1:] {
		tagOptsKeyVals := strings.Split(tagVals[1:][i], "=")

		switch tagOptsKeyVals[0] {
		case tagOptsNested:
			p.setTagOpts(tagOptsNested, "true")
		// case tagOptsExpandKey:
		// 	p.setTagOpts(tagOptsExpandKey, "true")
		case tagOptsUnsetKey:
			p.setTagOpts(tagOptsUnsetKey, "true")
		case tagOptsValidateKey:
			p.setTagOpts(tagOptsValidateKey, tagOptsKeyVals[1])
		case tagOptsDefaultKey:
			p.setTagOpts(tagOptsDefaultKey, tagOptsKeyVals[1])
		default:
			return errorx.New(fmt.Sprintf("unrecognized field option key: %s", tagOptsKeyVals[0]))
		}
	}

	return nil
}
