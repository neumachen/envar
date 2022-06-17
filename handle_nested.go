package envar

import (
	"fmt"
	"reflect"

	"github.com/neumachen/errorx"
)

var ErrNestedNotStruct = func(rValue reflect.Value) error {
	return errorx.New(fmt.Sprintf("field: %v is not a struct but has nested tag option", rValue.Elem().Type().String()))
}

func handleNested(
	parserCtx *ParserCtx,
	pField *parsedField,
	rValue reflect.Value,
) (ParserCtxAccessor, error) {
	if rValue.Kind() == reflect.Ptr {
		if rValue.IsNil() {
			rValue.Set(reflect.New(rValue.Type().Elem()))
		}
		rValue = rValue.Elem()
	}
	if rValue.Kind() != reflect.Struct {
		return nil, ErrNestedNotStruct(rValue)
	}
	return parse(parserCtx, rValue)
}

func handledNestedSlice(
	parserCtx *ParserCtx,
	pField *parsedField,
	rValues reflect.Value,
) error {
	if rValues.Kind() != reflect.Slice {
		return errorx.New(fmt.Sprintf("nested field: %v is not a slice", rValues.Type().Elem().Name()))
	}

	result := reflect.MakeSlice(pField.GetStructField().Type, 0, rValues.Len())
	for i := 0; i < rValues.Len(); i++ {
		if rValues.Index(i).Kind() != reflect.Struct {
			return ErrNestedNotStruct(rValues.Index(i))
		}
		nestedValue := reflect.New(rValues.Type().Elem())
		if _, err := parse(parserCtx, nestedValue); err != nil {
			return errorx.New(err)
		}
		result = reflect.Append(result, nestedValue)
	}
	rValues.Set(result)
	return nil
}
