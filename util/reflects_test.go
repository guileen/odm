package util

import (
	"fmt"
	"reflect"
	"testing"
)

type Book struct {
	Author string `odm:"hashKey"`
	Title  string `odm:"rangeKey"`
	Age    int
}

func TestExtractModelInfo(t *testing.T) {
	type args struct {
		pkField string
		skField string
		model   interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]interface{}
		want1 map[string]interface{}
	}{
		{
			"Extract PK only", args{"Author", "", &Book{"Tome", "Hello", 15}},
			map[string]interface{}{"Author": "Tome"}, map[string]interface{}{"Title": "Hello", "Age": 15},
		},
		{
			"Extract PK only", args{"Author", "Title", &Book{"Tome", "Hello", 15}},
			map[string]interface{}{"Author": "Tome", "Title": "Hello"}, map[string]interface{}{"Age": 15},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ExtractModelInfo(tt.args.pkField, tt.args.skField, tt.args.model)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractModelInfo() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ExtractModelInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMapToExpression(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 map[string]interface{}
	}{
		{"PlainObject", args{map[string]interface{}{"Author": "Tom", "Title": "Hello", "Age": 13}}, "Author=:Author and Title=:Title and Age=:Age", map[string]interface{}{":Author": "Tom", ":Title": "Hello", ":Age": 13}},
		// TODO: Add test cases for a.b=:a_b
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := MapToExpression(tt.args.m)
			if got != tt.want {
				t.Errorf("MapToExpression() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("MapToExpression() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func ExampleClearSlice() {
	book := []Book{Book{"Tome", "Hi", 1}}
	fmt.Println("begin", book)
	ClearSlice(&book)
	fmt.Println("final", book)
	// Output:
	// begin [{Tome Hi 1}]
	// final []
}
