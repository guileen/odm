package odm

import (
	"reflect"
	"sort"
	"strings"

	"git.devops.com/go/odm/util"
)

type Map map[string]interface{}

type Model interface {
}

type TableMeta struct {
	TableName string
	PK        *FieldDefine
	SK        *FieldDefine
	Fields    []*FieldDefine
}

type FieldDefine struct {
	ModelFieldName  string
	SchemaFieldName map[string]string
	// The data type for the attribute, where:
	//
	//    * S - the attribute is of type String
	//
	//    * N - the attribute is of type Number
	//
	//    * B - the attribute is of type Binary
	//
	// AttributeType is a required field
	Type      string
	PK        bool
	SK        bool
	OmitEmpty bool
}

func (f *FieldDefine) GetDBFieldName(dbname string) string {
	name := f.SchemaFieldName[dbname]
	if name == "" {
		return f.SchemaFieldName["json"]
	}
	return name
}

type TableConfig struct {
	Name     string
	UseCache bool
	TTL      int64
}

type TableConfigGetter interface {
	TableConfig() *TableConfig
}

// GetModelMeta 根据指针获取表的元信息
func GetModelMeta(model Model) *TableMeta {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.String {
		return &TableMeta{
			TableName: model.(string),
		}
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	meta := &TableMeta{
		Fields: []*FieldDefine{},
	}
	if getter, ok := model.(TableConfigGetter); ok {
		meta.TableName = getter.TableConfig().Name
	}
	if meta.TableName == "" {
		// meta.Name = inflection.Plural(util.ToSnakeCase(t.Name()))
		meta.TableName = util.ToSnakeCase(t.Name())
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fd := getFieldDefine(&f)
		if fd == nil {
			continue
		}
		meta.Fields = append(meta.Fields, fd)
		if fd.PK {
			meta.PK = fd
		} else if fd.SK {
			meta.SK = fd
		}
	}
	sort.Slice(meta.Fields, func(i, j int) bool {
		f1 := meta.Fields[i]
		f2 := meta.Fields[j]
		if f1.PK {
			return true
		}
		if f2.PK {
			return false
		}
		if f1.SK {
			return true
		}
		if f2.SK {
			return false
		}
		return strings.Compare(f1.ModelFieldName, f2.ModelFieldName) < 0
	})

	return meta
}

var typeOfBytes = reflect.TypeOf([]byte(nil))

func getFieldDefine(f *reflect.StructField) *FieldDefine {
	t := ""
	switch f.Type.Kind() {
	case reflect.String:
		t = "S"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		t = "N"
	case reflect.Slice:
		if f.Type == typeOfBytes {
			t = "B"
		}
	default:
		// Not support field type. won't create field.
		return nil
	}
	odmTags := strings.Split(f.Tag.Get("odm"), ",")
	jsonTags := strings.Split(f.Tag.Get("json"), ",")
	dyTags := strings.Split(f.Tag.Get("dynamodbav"), ",")
	if odmTags[0] == "" && jsonTags[0] == "-" {
		return nil
	}
	// snakeName := util.ToSnakeCase(f.Name)
	d := &FieldDefine{
		ModelFieldName: f.Name,
		Type:           t,
		PK:             odmTags[0] == "PK" || odmTags[0] == "hashkey",
		SK:             odmTags[0] == "SK" || odmTags[0] == "rangekey",
		OmitEmpty:      len(jsonTags) > 1 && jsonTags[1] == "omitempty",
		SchemaFieldName: map[string]string{
			"json":     util.StringsOr(jsonTags[0], f.Name),
			"dynamodb": util.StringsOr(dyTags[0], jsonTags[0], f.Name),
		},
	}
	return d
}
