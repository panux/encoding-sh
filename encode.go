package sh

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

var basic = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

func encode(i interface{}) string {
	switch i.(type) {
	case int:
		return strconv.Itoa(i.(int))
	case int8:
		return encode(int(i.(int8)))
	case int16:
		return encode(int(i.(int16)))
	case int32:
		return encode(int(i.(int32)))
	case int64:
		return encode(int(i.(int64)))
	case uint:
		return encode(int(i.(uint)))
	case uint8:
		return encode(int(i.(uint8)))
	case uint16:
		return encode(int(i.(uint16)))
	case uint32:
		return encode(int(i.(uint32)))
	case uint64:
		return encode(int(i.(uint64)))
	case string:
		simple := basic + "/"
		smap := make(map[rune]bool)
		for _, c := range []rune(simple) {
			smap[c] = true
		}
		isSimple := func() bool {
			for _, c := range []rune(i.(string)) {
				if !smap[c] {
					return false
				}
			}
			return true
		}()
		if isSimple {
			return i.(string)
		}
		return fmt.Sprintf("%q", i)
	case fmt.Stringer:
		return encode(i.(fmt.Stringer).String())
	default:
		rt := reflect.TypeOf(i)
		ival := reflect.ValueOf(i)
		switch rt.Kind() {
		case reflect.Slice:
			o := make([]string, ival.Len())
			for j := 0; j < ival.Len(); j++ {
				v := encode(ival.Index(j).Interface())

				f, _ := utf8.DecodeRuneInString(v)
				if f == '"' {
					l, _ := utf8.DecodeLastRuneInString(v)
					if l == '"' {
						v = v[1 : len(v)-1] //Remove quotes
					}
				}
				o[j] = v
			}
			return fmt.Sprintf(`"%s"`, strings.Join(o, " "))
		case reflect.Array:
			return encode(reflect.Indirect(ival.Addr()).Slice(0, ival.Len()).Interface())
		case reflect.Ptr:
			return encode(ival.Elem().Interface())
		}
	}
	panic("bad type")
}

//Encode encodes a map with string keys or a struct to shell format
func Encode(i interface{}) (b []byte, e error) {
	defer func() {
		err := recover()
		if err != nil {
			e = err.(error)
			b = nil
		}
	}()
	t := reflect.TypeOf(i)
	ival := reflect.ValueOf(i)
	vmap := make(map[rune]bool)
	for _, c := range []rune(basic + "_") {
		vmap[c] = true
	}
	switch t.Kind() {
	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			panic(errors.New("map key must be string"))
		}
		strs := make([]string, ival.Len())
		for j, k := range ival.MapKeys() {
			for _, c := range []rune(k.String()) {
				if !vmap[c] {
					panic(fmt.Errorf("invalid name character %s", string(c)))
				}
			}
			strs[j] = fmt.Sprintf("%s=%s", strings.ToUpper(k.String()), encode(ival.MapIndex(k)))
		}
		return []byte(strings.Join(strs, "\n")), nil
	case reflect.Struct:
		strs := make([]string, t.NumField())
		for j := 0; j < len(strs); j++ {
			strs[j] = fmt.Sprintf("%s=%s", strings.ToUpper(t.Field(j).Name), encode(ival.Field(j).Interface()))
		}
		return []byte(strings.Join(strs, "\n")), nil
	}
	panic(errors.New("unsupported type"))
}
