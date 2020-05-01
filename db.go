package odm

// Config is Connection Configuration.
type Config interface {
}

// DB 是对数据库的抽象
type DB interface {
	Table(model Model) Table
	GetTable(name string) Table
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

// Pool 连接池
type Pool interface {
	DB()
}

type TransGet struct {
	TableName            string
	ProjectionExpression string
	NameParams           map[string]string
	Key                  Key
}

type TransWrite struct {
	TableName string
	// PUT UPDATE DELETE
	Operation string
	// Condition check
	ConditionExpression string                 `type:"string"`
	NameParams          map[string]string      `type:"map"`
	ValueParams         map[string]interface{} `type:"map"`
	// required for PUT
	Item interface{}
	// required for DELETE and UPDATE
	Key Key `type:"map" required:"true"`
	// required for UPDATE
	UpdateExpression string
	// enum: NONE and ALL_OLD  (for PUT, DELETE)
	// enum: NONE, ALL_OLD, UPDATED_OLD, ALL_NEW, UPDATED_NEW (for UPDATE)
	ReturnValuesOnConditionCheckFailure string
}

type BatchGet struct {
	TableName            string
	ConsistentRead       bool `type:"boolean"`
	ProjectionExpression string
	NameParams           map[string]string
	Keys                 []Key
}

type BatchWrite struct {
	TableName  string
	PutItems   interface{}
	DeleteKeys []Key
}
