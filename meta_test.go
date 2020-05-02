package odm

import (
	"reflect"
	"testing"
)

type Book struct {
	Author string `odm:"partitionKey"`
	Title  string `odm:"sortingKey"`
	Age    int
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
	type args struct {
		model Model
	}
	tests := []struct {
		name string
		args args
		want *TableMeta
	}{
		{"Normal", args{new(Book)}, &TableMeta{TableName: "book", PartitionKey: "author", SortingKey: "title"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetModelMeta(tt.args.model); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetModelMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}
