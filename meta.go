package odm

import (
	"fmt"
	"reflect"
	"strings"

	"git.devops.com/go/odm/util"
)

// GetModelMeta 根据指针获取表的元信息
func GetModelMeta(model Model) *TableMeta {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	meta := new(TableMeta)
	if getter, ok := model.(TableConfigGetter); ok {
		meta.TableName = getter.TableConfig().Name
	}
	if meta.TableName == "" {
		// meta.Name = inflection.Plural(util.ToSnakeCase(t.Name()))
		meta.TableName = util.ToSnakeCase(t.Name())
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("odm")
		if tag == "" {
			continue
		}
		args := strings.Split(tag, ",")
		for j := 0; j < len(args); j++ {
			if args[j] == "" {
				continue
			}
			switch strings.TrimSpace(args[j]) {
			case "partitionKey":
				meta.PartitionKey = util.ToSnakeCase(f.Name)
			case "sortingKey":
				meta.SortingKey = util.ToSnakeCase(f.Name)
			default:
				fmt.Printf("Model attribute not supported: %s\n", args[j])
			}
		}
	}
	return meta
}
