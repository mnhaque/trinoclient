package trinodriver

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type UnsupportedArgError struct {
	t string
}

func (e UnsupportedArgError) Error() string {
	return fmt.Sprintf("trino: unsupported arg type: %s", e.t)
}

// Numeric is a string representation of a number, such as "10", "5.5" or in scientific form
// If another string format is used it will error to serialise
type Numeric string

// Serial converts any supported value to its equivalent string for as a Trino parameter
// See https://trino.io/docs/current/language/types.html
func Serial(v interface{}) (string, error) {
	switch x := v.(type) {
	case nil:
		return "", UnsupportedArgError{"<nil>"}

	// numbers convertible to int
	case int8:
		return strconv.Itoa(int(x)), nil
	case int16:
		return strconv.Itoa(int(x)), nil
	case int32:
		return strconv.Itoa(int(x)), nil
	case int:
		return strconv.Itoa(x), nil
	case uint16:
		return strconv.Itoa(int(x)), nil

	case int64:
		return strconv.FormatInt(x, 10), nil

	case uint32:
		return strconv.FormatUint(uint64(x), 10), nil
	case uint:
		return strconv.FormatUint(uint64(x), 10), nil
	case uint64:
		return strconv.FormatUint(x, 10), nil

		// float32, float64 not supported because digit precision will easily cause large problems
	case float32:
		return "", UnsupportedArgError{"float32"}
	case float64:
		return "", UnsupportedArgError{"float64"}

	case Numeric:
		if _, err := strconv.ParseFloat(string(x), 64); err != nil {
			return "", err
		}
		return string(x), nil

		// note byte and uint are not supported, this is because byte is an alias for uint8
		// if you were to use uint8 (as a number) it could be interpreted as a byte, so it is unsupported
		// use string instead of byte and any other uint/int type for uint8
	case byte:
		return "", UnsupportedArgError{"byte/uint8"}

	case bool:
		return strconv.FormatBool(x), nil

	case string:
		return "'" + strings.Replace(x, "'", "''", -1) + "'", nil

		// TODO - []byte should probably be matched to 'VARBINARY' in trino
	case []byte:
		return "", UnsupportedArgError{"[]byte"}

		// time.Time and time.Duration not supported as time and date take several different formats in Trino
	case time.Time:
		return "", UnsupportedArgError{"time.Time"}
	case time.Duration:
		return "", UnsupportedArgError{"time.Duration"}

		// TODO - json.RawMesssage should probably be matched to 'JSON' in Trino
	case json.RawMessage:
		return "", UnsupportedArgError{"json.RawMessage"}
	}

	if reflect.TypeOf(v).Kind() == reflect.Slice {
		x := reflect.ValueOf(v)
		if x.IsNil() {
			return "", UnsupportedArgError{"[]<nil>"}
		}

		slice := make([]interface{}, x.Len())

		for i := 0; i < x.Len(); i++ {
			slice[i] = x.Index(i).Interface()
		}

		return serialSlice(slice)
	}

	if reflect.TypeOf(v).Kind() == reflect.Map {
		// are Trino MAPs indifferent to order? Golang maps are, if Trino aren't then the two types can't be compatible
		return "", UnsupportedArgError{"map"}
	}

	// TODO - consider the remaining types in https://trino.io/docs/current/language/types.html (Row, IP, ...)

	return "", UnsupportedArgError{fmt.Sprintf("%T", v)}
}

func serialSlice(v []interface{}) (string, error) {
	ss := make([]string, len(v))

	for i, x := range v {
		s, err := Serial(x)
		if err != nil {
			return "", err
		}
		ss[i] = s
	}

	return "ARRAY[" + strings.Join(ss, ", ") + "]", nil
}
