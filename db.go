package odm

// Config is Connection Configuration.
type Config interface {
}

// ODMDB 是对数据库的抽象
type ODMDB struct {
	DialectDB
}

type Dialect interface {
	Open(connectString string) (DialectDB, error)
	GetName() string
}

type DialectDB interface {
	GetTable(tableName string) Table
	GetDialectTable(*TableMeta) Table
	// 对多表读取，不保证一致性
	BatchGetItem(options []BatchGet, unprocessedItems *[]BatchGet, results ...interface{}) error
	// 对多表增、删，不保证一致性
	BatchWriteItem(options []BatchWrite, unprocessedItems *[]BatchWrite) error
	// 一致性读，一起成功
	// TransactGetItem([]{get1, get2}, &result1, &result2)
	TransactGetItems(gets []TransGet, results ...Model) error
	// 一致性写，一起成功、一起失败
	TransactWriteItems(writes []TransWrite) error
	Close()
}

type TransGet struct {
	TableName  string
	Select     string
	NameParams map[string]string
	Key        Key
}

type TransWrite struct {
	TableName string
	// PUT UPDATE DELETE
	Operation string
	// Filter check
	Filter      string
	NameParams  map[string]string
	ValueParams map[string]interface{}
	// required for PUT
	Item interface{}
	// required for DELETE and UPDATE
	Key Key
	// required for UPDATE
	Update string
	// enum: NONE and ALL_OLD  (for PUT, DELETE)
	// enum: NONE, ALL_OLD, UPDATED_OLD, ALL_NEW, UPDATED_NEW (for UPDATE)
	ReturnValuesOnConditionCheckFailure string
}

type BatchGet struct {
	TableName  string
	Consistent bool
	Select     string
	NameParams map[string]string
	Keys       []Key
}

type BatchWrite struct {
	TableName  string
	PutItems   interface{}
	DeleteKeys []Key
}

func (db *ODMDB) Table(model Model) Table {
	metaInfo := GetModelMeta(model)
	return db.GetDialectTable(metaInfo)
}

func (db *ODMDB) GetTable(tableName string) Table {
	return db.GetTable(tableName)
}
