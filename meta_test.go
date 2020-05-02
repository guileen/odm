package odm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Book struct {
	Author     string `odm:"PK" dynamodbav:"author"`
	Title      string `odm:"SK" json:"title" dynamodbav:"subject"`
	Age        byte
	FooBar     float32 `json:"foo_bar"`
	Img        []byte
	NotSupport map[string]string
	Ignored    string `json:"-"`
}

func (b *Book) GetConfig() *TableConfig {
	return &TableConfig{
		Name:     "books",
		UseCache: true,
		TTL:      60,
	}
}

func (b *Book) TableName() string {
	return "books"
}

func TestGetModelMeta(t *testing.T) {
	meta := GetModelMeta(&Book{})
	assert.Equal(t, "book", meta.TableName)
	assert.Equal(t, &FieldDefine{
		ModelFieldName: "Author",
		SchemaFieldName: map[string]string{
			"json":     "Author",
			"dynamodb": "author",
		},
		PK:   true,
		Type: "S",
	}, meta.PK)
	assert.Equal(t, &FieldDefine{
		ModelFieldName: "Title",
		SchemaFieldName: map[string]string{
			"json":     "title",
			"dynamodb": "subject",
		},
		SK:   true,
		Type: "S",
	}, meta.SK)
	assert.Equal(t, &FieldDefine{
		ModelFieldName: "Author",
		SchemaFieldName: map[string]string{
			"json":     "Author",
			"dynamodb": "author",
		},
		PK:   true,
		Type: "S",
	}, meta.Fields[0])
	assert.Equal(t, &FieldDefine{
		ModelFieldName: "Title",
		SchemaFieldName: map[string]string{
			"json":     "title",
			"dynamodb": "subject",
		},
		SK:   true,
		Type: "S",
	}, meta.Fields[1])
	assert.Equal(t, &FieldDefine{
		ModelFieldName: "Age",
		SchemaFieldName: map[string]string{
			"json":     "Age",
			"dynamodb": "Age",
		},
		Type: "N",
	}, meta.Fields[2])
	assert.Equal(t, &FieldDefine{
		ModelFieldName: "FooBar",
		SchemaFieldName: map[string]string{
			"json":     "foo_bar",
			"dynamodb": "foo_bar",
		},
		Type: "N",
	}, meta.Fields[3])
	assert.Equal(t, &FieldDefine{
		ModelFieldName: "Img",
		SchemaFieldName: map[string]string{
			"json":     "Img",
			"dynamodb": "Img",
		},
		Type: "B",
	}, meta.Fields[4])
}
