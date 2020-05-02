package odm

// Table 表的基本底层操作, 按照DynamoDB的操作进行对应抽象
type Table interface {
	// GetDB returns instanceof DB
	GetDB() DialectDB
	// put a item, will replace entire item.
	PutItem(item Model, cond *WriteOption) error
	// Update attributes. item will fill base on ReturnValues.
	UpdateItem(partitionKey interface{}, sortingKey interface{}, updateExpr string, opt *WriteOption, result Model) error
	// get a item
	GetItem(partitionKey interface{}, sortingKey interface{}, opt *GetOption, result Model) error
	// returns deleted item
	DeleteItem(partitionKey interface{}, sortingKey interface{}, opt *WriteOption, result Model) error
	// Query and fill in items, StartKey will be replaced after query
	// result is slice of Model
	// Example:
	// 		offsetKey := make(odm.Key)
	// 		items := []Item{}
	// 		table.Query(query, offsetKey, &items)
	Query(query *QueryOption, offsetKey Map, results interface{}) error
}

type WriteOption struct {
	Filter      string
	NameParams  map[string]string
	ValueParams Map
}

type GetOption struct {
	Consistent bool
	Select     string
	NameParams map[string]string
}

type QueryOption struct {
	// 查询表达式
	Filter    string
	KeyFilter string
	// Select 对应着 ProjectionExpression，更简短更易理解
	Select string

	// 查询参数
	NameParams  map[string]string
	ValueParams Map

	// 查询限制
	Consistent bool
	Limit      int64
	IndexName  string
	Desc       bool // 默认升序，默认false。向其他数据库迁移的时候，这里需要注意，可能不兼容，需要提供额外的排序信息。
	// 这个决定了Key的扫描方向
	// ScanIndexForward       *bool              `type:"boolean"`   default is true. ascending

	// ExclusiveStartKey 由 StartKey 参数提供
	// ExclusiveStartKey      Map               `type:"map"`

	// QueryInput.Select 并没有什么用，只是关于信息范围的。无需指定
	// Select string `type:"string" enum:"Select"`
}

type ScanOption struct {
	QueryOption
	// ScanOption 没有KeyConditionExpression
	// 下面两个字段是关于多进程并行扫描的
	Segment       int64
	TotalSegments int64
}
