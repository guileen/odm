package types

// Table 表的基本底层操作, 按照DynamoDB的操作进行对应抽象
type Table interface {
	// GetDB returns instanceof DB
	GetDB() DB
	// put a item, will replace entire item.
	PutItem(item Model, cond *Condition) error
	// Update attributes. item will fill base on ReturnValues.
	UpdateItem(key Key, updateExpr string, opt *Condition, result Model) error
	// get a item
	GetItem(key Key, opt *GetOption, result Model) error
	// returns deleted item
	DeleteItem(key Key, opt *Condition, result Model) error
	// Query and fill in items, StartKey will be replaced after query
	// result is slice of Model
	// Example:
	// 		offsetKey := make(types.Key)
	// 		items := []Item{}
	// 		table.Query(query, offsetKey, &items)
	Query(query *QueryOption, offsetKey Key, results interface{}) error
}

type Key map[string]interface{}

type Condition struct {
	ConditionExpression string            `type:"string"`
	NameParams          map[string]string `type:"map"`
	ValueParams         Map               `type:"map"`
}

type GetOption struct {
	ConsistentRead       bool `type:"boolean"`
	ProjectionExpression string
	NameParams           map[string]string
}

type UpdateOption struct {
	// Condition
	ConditionExpression string            `type:"string"`
	NameParams          map[string]string `type:"map"`
	ValueParams         Map               `type:"map"`
	// UpdateExpression          string
	// ReturnValues     *string `type:"string" enum:"ReturnValue"`
}

type DeleteOption struct {
	// Condition
	ConditionExpression string            `type:"string"`
	NameParams          map[string]string `type:"map"`
	ValueParams         Map               `type:"map"`
	// ReturnValues              *string
}

type QueryOption struct {
	// 查询表达式
	FilterExpression       string `type:"string"`
	KeyConditionExpression string `type:"string"`
	ProjectionExpression   string `type:"string"`

	// 查询参数
	NameParams  map[string]string `type:"map"`
	ValueParams Map               `type:"map"`

	// 查询限制
	ConsistentRead bool   `type:"boolean"`
	Limit          int64  `min:"1" type:"integer"`
	IndexName      string `min:"3" type:"string"`
	Desc           bool   // 默认升序，默认false。向其他数据库迁移的时候，这里需要注意，可能不兼容，需要提供额外的排序信息。
	// 这个决定了Key的扫描方向
	// ScanIndexForward       *bool              `type:"boolean"`   default is true. ascending

	// ExclusiveStartKey 由 StartKey 参数提供
	// ExclusiveStartKey      Map               `type:"map"`

	// Select 并没有什么用，只是关于信息范围的。无需指定
	// Select string `type:"string" enum:"Select"`
}

type ScanOption struct {
	QueryOption
	// ScanOption 没有KeyConditionExpression
	// 下面两个字段是关于多进程并行扫描的
	Segment       *int64 `type:"integer"`
	TotalSegments *int64 `min:"1" type:"integer"`
}
