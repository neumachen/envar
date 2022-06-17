package envar

import (
	"os"
	"reflect"

	"github.com/neumachen/errorx"
)

func parse(parserCtx *ParserCtx, refValue reflect.Value) (ParserCtxAccessor, error) {
	for i := 0; i < refValue.Type().NumField(); i++ {
		refField := refValue.Field(i)
		if !refField.CanSet() {
			continue
		}

		parsedField, err := newParsedField(parserCtx, refValue.Type().Field(i))
		if err != nil {
			return nil, errorx.New(err)
		}
		if parsedField != nil {
			if parsedField.isNested() {
				if refField.Kind() == reflect.Slice {
					if err := handledNestedSlice(parserCtx, parsedField, refField); err != nil {
						return nil, errorx.New(err)
					}
					continue
				}
				if _, err := handleNested(parserCtx, parsedField, refField); err != nil {
					return nil, errorx.New(err)
				}
				continue
			}
			if err := parsedField.setEnvValue(parserCtx.GetEnvVarsMap()); err != nil {
				return nil, errorx.New(err)
			}
			if err := parsedField.validate(parserCtx); err != nil {
				return nil, errorx.New(err)
			}
			if err := parsedField.setField(parserCtx, refField); err != nil {
				return nil, errorx.New(err)
			}
			if unset := parsedField.unsetEnv(); unset {
				os.Unsetenv(parsedField.GetEnvKey())
			}
			continue
		}
	}

	return parserCtx, nil
}

// ErrNotAStructPtr is returned if you pass something that is not a pointer to a
// Struct to Parse.
var ErrNotAStructPtr = errorx.New("env: expected a pointer to a Struct")

// Parse parses a struct containing `env` tags and loads its values from
// environment variables.
func Parse(v interface{}, setterFuncs ...ParserCtxFuncSetter) (ParserCtxGetter, error) {
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return nil, ErrNotAStructPtr
	}
	refValue := ptrRef.Elem()
	if refValue.Kind() != reflect.Struct {
		return nil, ErrNotAStructPtr
	}
	parserCtx := &ParserCtx{}

	for i := range defaultParserCtxSetters {
		if err := defaultParserCtxSetters[i](parserCtx); err != nil {
			return nil, errorx.New(err)
		}
	}

	for i := range setterFuncs {
		if err := setterFuncs[i](parserCtx); err != nil {
			return nil, errorx.New(err)
		}
	}
	parserCtx.envVarsMap = parserCtx.envVarsLoaderFunc()

	return parse(parserCtx, refValue)
}
