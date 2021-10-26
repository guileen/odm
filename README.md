# 该项目目前仅供演示使用，后续代码已闭源。

# golang ODM(Object Document Mapper) abstract layer

1. 首先支持 DynamoDB，按DynamoDB的习惯来操作。
2. 后续支持 MongoDB

## 建立连接 odm.Open(dbtype, connect_string)
```
import (
	"git.devops.com/go/odm"
	_ "git.devops.com/go/odm/dynamo"
)
db,err := odm.Open("dynamodb", "AccessKey=123;SecretKey=456;Token=789;Region=localhost;Endpoint=http://127.0.0.1:8000")
db.Close()
```

`db,err := odm.Open("mysql", "db_user:password@tcp(localhost:3306)/my_db")`

NOTE: 业务层在使用时不需要关心连接池

也可以使用AWS的底层配置对象来实现

```
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

...
creds := credentials.NewStaticCredentials("123", "123", "")

db, err := dynamo.OpenDB(&aws.Config{
	Credentials: creds,
	Endpoint:    aws.String(END_POINT),
	Region:      aws.String("localhost"),
})
```

## Scheme 操作
TODO 根据Model定义生成表
对表的创建、建立索引由运维手动完成。暂不由代码控制。

## Table 接口

Table 用于抽象数据表的操作，对应Dynamo的Table，MongoDB的Collection，MySQL的Table。操作以Dynamo为基础进行设计。

有两种方式获得Table。 `db.Table("table_name")` 这种方式只能获取已经存在于数据库中的表。`db.Table(&Book{})` 这种方式可以获得Book对应的表，如果表不存在则会在*localhost测试环境*自动创建（*生产环境仍要运维创建*）。为了能够获得创建表的必要信息，我们需要对Table增加注解。

```
type Book struct {
	Author string `odm:"PK"`  // odm:"PK" 注解标记主键
	Title  string `odm:"SK"`  // odm:"SK" 注解标记排序
	Age    int64
	// 自定义数据库字段
	JSONInfo  string `json:"json_info"`
	DyTagInfo string `json:"dyInfo" dynamodbav:"dy_info"`
}
```

### Map 类型

`type Map map[string]interface{}`

可以使用来方便的创建一个Map
```
odm.Map{
	"Author": "Tom",
	"Title": "Hello world",
}
```

### Model 类型

`type Model interface{}` 仅仅是一个指针，可以是任何结构体。

### 操作Options 类型

为了方便操作，简化了一些AWS SDK的字段名。

- NameParams 对应 ExpressionAttributeNames
- ValueParams 对应 ExpressionAttributeValues
- Filter 对应 ConditionExpression
- Select 对应 ProjectionExpression
- KeyFilter 对应 KeyConditionExpression

```
type WriteOption struct {
	ConditionExpression       *string
    NameParams
    ValueParams
}
```
Condition 类型是一个条件表达式，仅当表达式成立时，操作才能成功。

### PutItem(item Model, opt WriteOption, ) error
PutItem 操作。替换整个item。

### UpdateItem(pk interface{}, sk interface{}, updateExpression string, opt WriteOption, item Model) error
Update 部分字段，根据ReturnValues返回数据到item中。

### GetItem(pk interface{}, sk interface{}, opt GetOption, item Model) error
Consistent 代表是否是一致性读。

### DeleteItem(pk interface{}, sk interface{}, opt WriteOption, item Model) error
被删除对象将填充到item。

### Query(QueryOption, offsetKey Map, items []Model) error
查询列表。
startKey 用来作性能优化。查询将从startKey开始。查询完成后，startKey将被更新。

```
type QueryOption struct {
	// 查询表达式
	Filter    string
	KeyFilter string
	Select string

	// 查询参数
	NameParams  map[string]string
	ValueParams Map

	// 查询限制
	Consistent bool
	Limit      int64
	IndexName  string
	Desc       bool // 默认升序，默认false。向其他数据库迁移的时候，这里需要注意，可能不兼容，需要提供额外的排序信息。
}
```
	
### TODO Scan(ScanOption, offsetKey Map, items []Model) error
TODO 暂不支持Scan, 与Query使用同一方法。


## RedisTable
TODO 使用Redis实现类似Table的功能。只能支持一些简单的查询。接口形式为Table

## CachedTable
TODO 组合Cache（RedisCache、MemoryCache、MixCache）和Table（DynamoTable、MongoTable）的一个实现，接口形式为Table。对PutItem、UpdateItem、DeleteItem、GetItem

## 缓存相关设计

### Cache 接口

Cache接口封装缓存常规操作。包含 GetItem、PutItem、DeleteItem

#### GetItem(key string) ([]byte, error)
#### PutItem(key string, []byte, ttl int64) error
#### DeleteItem(key string) ([]byte, error)

### Cache 实现
RedisCache、MemoryCache、MixCache（级联 MemoryCache 和 RedisCache）


## Test

```
go test ./...
```
