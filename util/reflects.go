package util

import (
	"reflect"
	"strings"

	"git.devops.com/go/odm/types"
)

// ExtractModelInfo 抽出 key 部分信息和其余部分信息
func ExtractModelInfo(pkField string, skField string, model types.Model) (types.Map, types.Map) {
	keys := make(types.Map)
	rest := make(types.Map)
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
func MapToExpression(m types.Map) (string, types.Map) {
	var strs []string
	attr := make(types.Map)
	for k, v := range m {
		// a=:a
		// TODO support a.b=:a_b
		strs = append(strs, k+"=:"+k)
		attr[":"+k] = v
	}
	return strings.Join(strs, " and "), attr
}
