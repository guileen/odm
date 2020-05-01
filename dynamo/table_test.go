package dynamo

import (
	"strconv"
	"testing"

	"git.devops.com/go/odm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/stretchr/testify/assert"
)

const END_POINT = "http://127.0.0.1:8000"

type Book struct {
	Author string
	Title  string
	Age    int64
	Info   string
}

func GetTestTable(t *testing.T) odm.Table {
	creds := credentials.NewStaticCredentials("123", "123", "")

	db, err := OpenDB(&aws.Config{
		Credentials: creds,
		Endpoint:    aws.String(END_POINT),
		Region:      aws.String("localhost"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, db)
	table := db.GetTable("book")
	return table
}

func TestTable_PutItem(t *testing.T) {
	t.Run("PutItem", func(t *testing.T) {
		book := &Book{
			Author: "Tom",
			Title:  "Hello",
			Age:    10,
		}
		table := GetTestTable(t)
		err := table.PutItem(book, nil)
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
		err := table.PutItem(book, nil)
		assert.NoError(t, err)
		book1 := &Book{}
		err = table.UpdateItem(odm.Key{"Author": "Tom", "Title": "2"}, "SET Info=:Info", &odm.WriteOption{
			ValueParams: odm.Map{
				":Info": "World",
			},
		}, book)
		assert.NoError(t, err)
		assert.Equal(t, &Book{
			Author: "Tom",
			Title:  "2",
			Age:    10,
			Info:   "World",
		}, book)
		table.GetItem(odm.Key{"Author": "Tom", "Title": "2"}, nil, book1)
		assert.Equal(t, &Book{
			Author: "Tom",
			Title:  "2",
			Age:    10,
			Info:   "World",
		}, book1)
	})
}

func TestTable_GetItem(t *testing.T) {
	t.Run("GetItem", func(t *testing.T) {
		book := &Book{
			Author: "Tom",
			Title:  "Hello",
			Age:    10,
		}
		table := GetTestTable(t)
		err := table.PutItem(book, nil)
		assert.NoError(t, err)
		book1 := &Book{}
		err = table.GetItem(odm.Key{"Author": "Tom", "Title": "Hello"}, nil, book1)
		assert.NoError(t, err)
		assert.Equal(t, &Book{
			Author: "Tom",
			Title:  "Hello",
			Age:    10,
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
		err := table.PutItem(book, nil)
		assert.NoError(t, err)
		err = table.DeleteItem(odm.Key{"Author": "Tom", "Title": "3"}, nil, nil)
		assert.NoError(t, err)
		book1 := &Book{}
		table.GetItem(odm.Key{"Author": "Tom", "Title": "3"}, nil, book1)
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
		table.PutItem(&allBooks[i], nil)
	}
	t.Run("ASC page", func(t *testing.T) {
		books := []Book{}
		offsetKey := make(odm.Key)
		err := table.Query(&odm.QueryOption{
			KeyFilter: "Author = :Author and Title > :Title",
			ValueParams: odm.Map{
				":Author": "Jack",
				":Title":  "Book",
			},
			Limit: 5,
		}, offsetKey, &books)
		assert.NoError(t, err)
		assert.Equal(t, allBooks[:5], books)
		err = table.Query(&odm.QueryOption{
			KeyFilter: "Author = :Author and Title > :Title",
			ValueParams: odm.Map{
				":Author": "Jack",
				":Title":  "Book",
			},
			Limit: 5,
		}, offsetKey, &books)
		assert.NoError(t, err)
		assert.Equal(t, allBooks[5:], books)
	})
	t.Run("DESC page", func(t *testing.T) {
		books := []Book{}
		offsetKey := make(odm.Key)
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
