package envar

import (
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ParaServices/errgo"
	"github.com/ParaServices/paratils"
)

// ParserFunc defines the signature of a function that can be used within `CustomParsers`.
type ParserFunc func(v string) (interface{}, error)

type ParserFuncMap map[string]ParserFunc

func (p ParserFuncMap) GetLength() int {
	return len(p)
}

func (p ParserFuncMap) Get(rType reflect.Type) ParserFunc {
	typeStr := rType.String()
	if rType.Kind() == reflect.Ptr {
		typeStr = rType.Elem().String()
	}
	f, ok := p[strings.ToLower(typeStr)]
	if !ok {
		return nil
	}
	return f
}

func (p ParserFuncMap) Add(rType reflect.Type, pFunc ParserFunc) {
	p[strings.ToLower(rType.String())] = pFunc
}

func defaultParserFuncs() ParserFuncMap {
	parserFuncs := map[reflect.Type]ParserFunc{
		reflect.TypeOf(false): func(v string) (interface{}, error) {
			return strconv.ParseBool(v)
		},
		reflect.TypeOf(string("")): func(v string) (interface{}, error) {
			return v, nil
		},
		reflect.TypeOf(int(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 32)
			return int(i), err
		},
		reflect.TypeOf(int8(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 8)
			return int(i), err
		},
		reflect.TypeOf(int16(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 16)
			return int(i), err
		},
		reflect.TypeOf(int32(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 32)
			return int(i), err
		},
		reflect.TypeOf(int64(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 64)
			return int(i), err
		},
		reflect.TypeOf(uint(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 32)
			return uint(i), err
		},
		reflect.TypeOf(uint8(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 8)
			return uint8(i), err
		},
		reflect.TypeOf(uint16(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 16)
			return uint16(i), err
		},
		reflect.TypeOf(uint32(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 32)
			return uint32(i), err
		},
		reflect.TypeOf(uint64(0)): func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 64)
			return uint64(i), err
		},
		reflect.TypeOf(float32(0.0)): func(v string) (interface{}, error) {
			return strconv.ParseFloat(v, 32)
		},
		reflect.TypeOf(float64(0.0)): func(v string) (interface{}, error) {
			return strconv.ParseFloat(v, 64)
		},
		reflect.TypeOf(url.URL{}): func(v string) (interface{}, error) {
			u, err := url.Parse(v)
			if err != nil {
				return nil, errgo.NewF("unable to parse URL: %v", err)
			}
			return *u, nil
		},
		reflect.TypeOf(os.File{}): func(v string) (interface{}, error) {
			if paratils.StringIsEmpty(v) {
				return nil, errgo.NewF("The file %v can not be empty", v)
			}

			fileInfo, err := os.Stat(v)
			if err != nil {
				if os.IsNotExist(err) {
					return nil, errgo.NewF("The file %v does not exist", v)
				}

				return nil, err
			}

			if fileInfo.IsDir() {
				return nil, errgo.NewF("The file %v is a directory", v)
			}

			f, err := os.Open(v)
			if err != nil {
				return nil, err
			}
			return *f, nil
		},
		reflect.TypeOf(time.Nanosecond): func(v string) (interface{}, error) {
			d, err := time.ParseDuration(strings.TrimSpace(v))
			if err != nil {
				return nil, errgo.NewF("unable to parse duration: %v", err)
			}
			return d, err
		},
	}

	defFuncs := make(ParserFuncMap)

	for k, v := range parserFuncs {
		defFuncs.Add(k, v)
	}
	return defFuncs
}
