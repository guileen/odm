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
	GetDialectTable(*TableMeta) Table
	DeleteTable(tableName string) error
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
	Key        Map
}

type TransWrite struct {
	TableName string
	// First Order
	PutItem interface{}
	// Second Order
	Update UpdateWrite
	// Third Order
	Delete DeleteWrite
}

type PutWrite struct {
	Item        interface{}
	Filter      string
	NameParams  map[string]string
	ValueParams Map
	// enum: NONE and ALL_OLD  (for PUT, DELETE)
	ReturnValuesOnConditionCheckFailure string
}

type UpdateWrite struct {
	Expression   string
	PartitionKey interface{}
	SortingKey   interface{}
	Filter       string
	NameParams   map[string]string
	ValueParams  Map
	// enum: NONE, ALL_OLD, UPDATED_OLD, ALL_NEW, UPDATED_NEW (for UPDATE)
	ReturnValuesOnConditionCheckFailure string
}

type DeleteWrite struct {
	PartitionKey interface{}
	SortingKey   interface{}
	Filter       string
	NameParams   map[string]string
	ValueParams  Map
	// enum: NONE and ALL_OLD  (for PUT, DELETE)
	ReturnValuesOnConditionCheckFailure string
}

type BatchGet struct {
	TableName  string
	Consistent bool
	Select     string
	NameParams map[string]string
	Keys       []Map
}

type BatchWrite struct {
	TableName  string
	PutItems   interface{}
	DeleteKeys []Map
}

func (db *ODMDB) Table(model Model) Table {
	metaInfo := GetModelMeta(model)
	return db.GetDialectTable(metaInfo)
}
