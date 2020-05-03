package odm

import (
	"errors"
	"flag"
	"os"
	"strings"
)

var dialectMap map[string]Dialect
var dropTableEnabled bool

func init() {
	dialectMap = make(map[string]Dialect)
	if strings.HasSuffix(os.Args[0], ".test") || flag.Lookup("test.v") != nil {
		dropTableEnabled = true
	}
}

func IsDropTableEnabled() bool {
	return dropTableEnabled
}

func RegisterDialect(dbtype string, opener Dialect) {
	dialectMap[dbtype] = opener
}

func GetDialect(dbtype string) Dialect {
	return dialectMap[dbtype]
}
func Open(dbtype string, connectString string) (*ODMDB, error) {
	dialect := GetDialect(dbtype)
	if dialect == nil {
		return nil, errors.New("No DB dialect <" + dbtype + "> register. Try `import \"git.devops.com/go/odm/dynamodb\"`")
	}
	dialectDB, err := dialect.Open(connectString)
	if err != nil {
		return nil, err
	}
	return &ODMDB{
		DialectDB: dialectDB,
	}, nil
}
