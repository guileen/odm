package util

import (
	"errors"
	"reflect"
	"strings"
)

// ExtractModelInfo 抽出 key 部分信息和其余部分信息
func ExtractModelInfo(pkField string, skField string, model interface{}) (map[string]interface{}, map[string]interface{}) {
	keys := make(map[string]interface{})
	rest := make(map[string]interface{})
	val := reflect.ValueOf(model).Elem()
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		f := val.Field(i)
		tf := t.Field(i)
		name := tf.Name
		if name == pkField || name == skField {
			keys[name] = f.Interface()
		} else {
			rest[name] = f.Interface()
		}
	}
	return keys, rest
}

// MapToExpression convert {"a":"123"} to "a=:a" and {":a", "123"}.
// returns expression and attribute Map
func MapToExpression(m map[string]interface{}) (string, map[string]interface{}) {
	var strs []string
	attr := make(map[string]interface{})
	for k, v := range m {
		// a=:a
		// TODO support a.b=:a_b
		strs = append(strs, k+"=:"+k)
		attr[":"+k] = v
	}
	return strings.Join(strs, " and "), attr
}

// ClearSlice ClearSlice(&items)  items=[]Something{...}
func ClearSlice(aryptr interface{}) error {
	t := reflect.TypeOf(aryptr)
	val := reflect.ValueOf(aryptr)
	if t.Kind() != reflect.Ptr {
		return errors.New("Input is not pointer to slice.")
	}
	t = t.Elem()
	val = val.Elem()
	if t.Kind() != reflect.Slice {
		return errors.New("Input is not slice.")
	}
	val.SetLen(0)
	return nil
}
