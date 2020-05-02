package odm

import "errors"

var dialectMap map[string]Dialect

func init() {
	dialectMap = make(map[string]Dialect)
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
