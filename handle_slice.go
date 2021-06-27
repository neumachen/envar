package envar

import (
	"encoding"
	"reflect"
	"strings"

	"github.com/ParaServices/paratils"
)

func handleSlice(
	parserCtx *ParserCtx,
	pField *parsedField,
	rValue reflect.Value,
) error {
	delim := parserCtx.GetEnvSliceDelim()
	if v := pField.GetEnvValue(); paratils.StringIsEmpty(v) {
		delim = defaultDelim
	}
	parts := strings.Split(pField.getFieldValue(), delim)

	fieldElem := pField.GetStructField().Type.Elem()
	if fieldElem.Kind() == reflect.Ptr {
		fieldElem = fieldElem.Elem()
	}

	if _, ok := reflect.New(fieldElem).Interface().(encoding.TextUnmarshaler); ok {
		return parseTextUnmarshalers(rValue, parts, pField.GetStructField())
	}

	parserFunc := parserCtx.GetParserFuncMap().Get(fieldElem)
	if paratils.IsNil(parserFunc) {
		return newNoParserError(pField.GetStructField())
	}

	result := reflect.MakeSlice(pField.GetStructField().Type, 0, len(parts))
	for _, part := range parts {
		r, err := parserFunc(part)
		if err != nil {
			return newParseError(pField.GetStructField(), err)
		}
		v := reflect.ValueOf(r).Convert(fieldElem)
		if pField.GetStructField().Type.Elem().Kind() == reflect.Ptr {
			v = reflect.New(fieldElem)
			v.Elem().Set(reflect.ValueOf(r).Convert(fieldElem))
		}
		result = reflect.Append(result, v)
	}
	rValue.Set(result)
	return nil
}
