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
	CreateTable(*TableMeta) error
	CreateTableIfNotExists(*TableMeta) error
	DropTable(tableName string) error
	// 对多表读取，不保证一致性
	BatchGetItem(options []*BatchGet, unprocessedItems *[]*BatchGet, results ...interface{}) error
	// 对多表增、删，不保证一致性
	BatchWriteItem(options []*BatchWrite, unprocessedItems *[]*BatchWrite) error
	// 一致性读，一起成功
	// TransactGetItem([]{get1, get2}, &result1, &result2)
	TransactGetItems(gets []*TransactGet, results ...Model) error
	// 一致性写，一起成功、一起失败
	TransactWriteItems(writes []*TransactWrite) error
	Close()
}

type TransactGet struct {
	TableName  string
	Select     string
	NameParams map[string]string
	Key        Map
}

type TransactWrite struct {
	ConditionCheck *ConditionCheck
	// First Order
	Put *Put
	// Second Order
	Update *Update
	// Third Order
	Delete *Delete
}

type ConditionCheck struct {
	TableName   string
	NameParams  map[string]string
	ValueParams Map
	PK          string
	SK          string
	// enum: NONE and ALL_OLD
	ReturnValuesOnConditionCheckFailure string
}

type Put struct {
	TableName   string
	Item        interface{}
	Condition   string
	NameParams  map[string]string
	ValueParams Map
	// enum: NONE and ALL_OLD  (for PUT, DELETE)
	ReturnValuesOnConditionCheckFailure string
}

type Update struct {
	TableName    string
	Expression   string
	PartitionKey interface{}
	SortingKey   interface{}
	Condition    string
	NameParams   map[string]string
	ValueParams  Map
	// enum: NONE, ALL_OLD, UPDATED_OLD, ALL_NEW, UPDATED_NEW (for UPDATE)
	ReturnValuesOnConditionCheckFailure string
}

type Delete struct {
	TableName    string
	PartitionKey interface{}
	SortingKey   interface{}
	Condition    string
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

func (db *ODMDB) ResetTable(model Model) (Table, error) {
	metaInfo := GetModelMeta(model)
	if IsDropTableEnabled() {
		db.DropTable(metaInfo.TableName)
	}
	err := db.CreateTableIfNotExists(metaInfo)
	if err != nil {
		return nil, err
	}
	return db.Table(model), nil
}

func (db *ODMDB) Table(model Model) Table {
	metaInfo := GetModelMeta(model)
	return db.GetDialectTable(metaInfo)
}
