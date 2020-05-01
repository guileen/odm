package util

import (
	"reflect"
	"testing"

	"git.devops.com/go/odm"
)

type Book struct {
	Author string `odm:"partitionKey"`
	Title  string `odm:"sortingKey"`
	Age    int
}

func TestExtractModelInfo(t *testing.T) {
	type args struct {
		pkField string
		skField string
		model   odm.Model
	}
	tests := []struct {
		name  string
		args  args
		want  odm.Map
		want1 odm.Map
	}{
		{
			"Extract PK only", args{"Author", "", &Book{"Tome", "Hello", 15}},
			odm.Map{"Author": "Tome"}, odm.Map{"Title": "Hello", "Age": 15},
		},
		{
			"Extract PK only", args{"Author", "Title", &Book{"Tome", "Hello", 15}},
			odm.Map{"Author": "Tome", "Title": "Hello"}, odm.Map{"Age": 15},
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
		m odm.Map
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 odm.Map
	}{
		{"PlainObject", args{odm.Map{"Author": "Tom", "Title": "Hello", "Age": 13}}, "Author=:Author and Title=:Title and Age=:Age", odm.Map{":Author": "Tom", ":Title": "Hello", ":Age": 13}},
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
