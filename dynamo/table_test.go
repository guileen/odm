package dynamo

import (
	"fmt"
	"strconv"
	"testing"

	"git.devops.com/go/odm"

	"github.com/stretchr/testify/assert"
)

type Book struct {
	Author string `odm:"PK"`
	Title  string `odm:"SK"`
	Age    int64
	// 自定义数据库字段
	JSONInfo  string `json:"json_info"`
	DyTagInfo string `json:"dyInfo" dynamodbav:"dy_info"`
}

// Localhost
var dbpath = "AccessKey=123;SecretKey=456;Token=789;Region=localhost;Endpoint=http://127.0.0.1:8000"

// Development environment
// var dbpath = "AccessKey=AKIAX24KZ5UPZSJY4FGV;SecretKey=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;Region=cn-northwest-1"

func GetTestTable(t *testing.T) odm.Table {
	db, err := odm.Open("dynamo", dbpath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	table := db.Table("book")
	return table
}

func TestTable_PutItem(t *testing.T) {
	t.Run("PutItem", func(t *testing.T) {
		resetDB(t)
		book := &Book{
			Author: "Tom",
			Title:  "Hello",
			Age:    10,
		}
		table := GetTestTable(t)
		err := table.PutItem(book, nil, nil)
		assert.NoError(t, err)
	})
}

func TestTable_UpdateItem(t *testing.T) {
	t.Run("UpdateItem", func(t *testing.T) {
		book := &Book{
			Author: "Tom",
			Title:  "2",
			Age:    10,
		}
		table := GetTestTable(t)
		err := table.PutItem(book, nil, nil)
		assert.NoError(t, err)
		book1 := &Book{}
		err = table.UpdateItem("Tom", "2", "SET json_info=:Info", &odm.WriteOption{
			ValueParams: odm.Map{
				":Info": "World",
			},
		}, book)
		assert.NoError(t, err)
		assert.Equal(t, &Book{
			Author: "Tom",
			Title:  "2",
			Age:    10,
			// JSONInfo is mapped with json_info
			JSONInfo: "World",
		}, book)
		table.GetItem("Tom", "2", nil, book1)
		assert.Equal(t, &Book{
			Author: "Tom",
			Title:  "2",
			Age:    10,
			//
			JSONInfo: "World",
		}, book1)
	})
}

func TestTable_GetItem(t *testing.T) {
	t.Run("GetItem", func(t *testing.T) {
		book := &Book{
			Author:    "Tom",
			Title:     "Hello",
			Age:       10,
			JSONInfo:  "JSON",
			DyTagInfo: "DyTag",
		}
		table := GetTestTable(t)
		err := table.PutItem(book, nil, nil)
		assert.NoError(t, err)
		book1 := &Book{}
		err = table.GetItem("Tom", "Hello", nil, book1)
		assert.NoError(t, err)
		assert.Equal(t, &Book{
			Author:    "Tom",
			Title:     "Hello",
			Age:       10,
			JSONInfo:  "JSON",
			DyTagInfo: "DyTag",
		}, book1)
	})
}

func TestTable_DeleteItem(t *testing.T) {
	t.Run("GetItem", func(t *testing.T) {
		book := &Book{
			Author: "Tom",
			Title:  "3",
			Age:    10,
		}
		table := GetTestTable(t)
		err := table.PutItem(book, nil, nil)
		assert.NoError(t, err)
		err = table.DeleteItem("Tom", "3", nil, nil)
		assert.NoError(t, err)
		book1 := &Book{}
		table.GetItem("Tom", "3", nil, book1)
		assert.Equal(t, &Book{
			Author: "",
			Title:  "",
			Age:    0,
		}, book1)
	})
}

func TestTable_Query(t *testing.T) {
	table := GetTestTable(t)
	allBooks := []Book{}
	for i := 0; i < 10; i++ {
		allBooks = append(allBooks, Book{
			Author: "Jack",
			Title:  "Book" + strconv.Itoa(i),
			Age:    int64(i),
		})
		table.PutItem(&allBooks[i], nil, nil)
	}
	t.Run("ASC page", func(t *testing.T) {
		books := []Book{}
		offsetKey := make(odm.Map)
		err := table.Query(&odm.QueryOption{
			KeyFilter: "Author = :Author and Title > :Title",
			ValueParams: odm.Map{
				":Author": "Jack",
				":Title":  "Book",
			},
			Limit: 3,
		}, offsetKey, &books)
		assert.NoError(t, err)
		assert.Equal(t, allBooks[:3], books)
		err = table.Query(&odm.QueryOption{
			KeyFilter: "Author = :Author and Title > :Title",
			ValueParams: odm.Map{
				":Author": "Jack",
				":Title":  "Book",
			},
			Limit: 5,
		}, offsetKey, &books)
		assert.NoError(t, err)
		assert.Equal(t, allBooks[3:8], books)
	})
	t.Run("DESC page", func(t *testing.T) {
		books := []Book{}
		offsetKey := make(odm.Map)
		err := table.Query(&odm.QueryOption{
			KeyFilter: "Author = :Author and Title > :Title",
			ValueParams: odm.Map{
				":Author": "Jack",
				":Title":  "Book",
			},
			Desc: true,
		}, offsetKey, &books)
		assert.NoError(t, err)
		assert.Equal(t, len(allBooks), len(books))
		assert.NotEqual(t, allBooks, books)
	})
	t.Run("Filter and Projection", func(t *testing.T) {
		books := []Book{}
		err := table.Query(&odm.QueryOption{
			KeyFilter: "Author = :Author and Title > :Title",
			ValueParams: odm.Map{
				":Author": "Jack",
				":Title":  "Book",
				":Age":    5,
			},
			Filter: "Age=:Age",
			Select: "Title, Age",
		}, nil, &books)
		assert.NoError(t, err)
		assert.Equal(t, &Book{
			Author: "",
			Title:  "Book5",
			Age:    5,
		}, &books[0])
	})
}

func ExampleTable_Query() {
	db, err := odm.Open("dynamo", dbpath)
	if err != nil {
		fmt.Errorf("Can't connect to dynamo db. %s\n", err.Error())
	}
	table := db.Table(&Book{})
	allBooks := []Book{}
	for i := 0; i < 10; i++ {
		allBooks = append(allBooks, Book{
			Author: "Alice",
			Title:  "Book" + strconv.Itoa(i),
			Age:    int64(i),
		})
		table.PutItem(&allBooks[i], nil, nil)
	}
	offsetKey := make(odm.Map)
	books := []Book{}
	err = table.Query(&odm.QueryOption{
		KeyFilter: "Author = :Author and Title > :Title",
		ValueParams: odm.Map{
			":Author": "Jack",
			":Title":  "Book2",
		},
		Limit: 1,
	}, offsetKey, &books)
	fmt.Println(books[0].Title)
	err = table.Query(&odm.QueryOption{
		KeyFilter: "Author = :Author and Title > :Title",
		ValueParams: odm.Map{
			":Author": "Jack",
			":Title":  "Book2",
		},
		Limit: 1,
	}, offsetKey, &books)
	fmt.Println(books[0].Title)
	// Output:
	// Book3
	// Book4
}
