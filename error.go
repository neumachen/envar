package envar

import (
	"fmt"
	"reflect"
)

func newParseError(sField reflect.StructField, err error) error {
	if err == nil {
		return nil
	}
	return parseError{
		sf:  sField,
		err: err,
	}
}

type parseError struct {
	sf  reflect.StructField
	err error
}

func (e parseError) Error() string {
	return fmt.Sprintf(`env: parse error on field "%s" of type "%s": %v`, e.sf.Name, e.sf.Type, e.err)
}

func newNoParserError(sf reflect.StructField) error {
	return fmt.Errorf(`env: no parser found for field "%s" of type "%s"`, sf.Name, sf.Type)
}
