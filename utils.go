package xcolt

import (
	"reflect"
	"sort"
)

func Max(n ...int) int {
	return sort.IntSlice(n)[0]
}

func Min(n ...int) int {
	return sort.IntSlice(n)[len(n)-1]
}

func If(condition bool, v1 interface{}, v2 interface{}) interface{} {
	if condition {
		return v1
	}
	return v2
}

func Filter(iter []string, f func(x string) bool) []string {
	result := []string{}
	for _, i := range iter {
		if f(i) {
			result = append(result, i)
		}
	}
	return result
}

func Mapper(name string) string {
	var s []byte
	for i, r := range []byte(name) {
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
			if i != 0 {
				s = append(s, '_')
			}
		}
		s = append(s, r)
	}
	return string(s)
}

func GetType(v interface{}) string {
	return reflect.TypeOf(v).Name()
}
