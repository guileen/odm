package odm

import "errors"

type DBOpener interface {
	Open(connectString string) (DB, error)
}

var openers map[string]DBOpener

func init() {
	openers = make(map[string]DBOpener)
}

func Register(dbtype string, opener DBOpener) {
	openers[dbtype] = opener
}

func Open(dbtype string, connectString string) (DB, error) {
	opener := openers[dbtype]
	if opener == nil {
		return nil, errors.New("No DB dialect <" + dbtype + "> register. Try `import \"git.devops.com/go/odm/dynamodb\"`")
	}
	return opener.Open(connectString)
}
