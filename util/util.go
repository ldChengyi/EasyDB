package util

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func SafeToString(v any) (string, error) {
	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.String:
		return v.(string), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", val.Float()), nil
	case reflect.Bool:
		return fmt.Sprintf("%t", val.Bool()), nil
	default:
		return "", errors.New("field value cannot be converted to string for substring/prefix index")
	}
}

func Compare(a, b interface{}) int {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	switch va.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai := va.Int()
		bi := vb.Int()
		if ai < bi {
			return -1
		} else if ai > bi {
			return 1
		}
		return 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		au := va.Uint()
		bu := vb.Uint()
		if au < bu {
			return -1
		} else if au > bu {
			return 1
		}
		return 0

	case reflect.Float32, reflect.Float64:
		af := va.Float()
		bf := vb.Float()
		if af < bf {
			return -1
		} else if af > bf {
			return 1
		}
		return 0

	case reflect.String:
		as := va.String()
		bs := vb.String()
		return strings.Compare(as, bs)
	}

	panic("unsupported type for compare")
}
