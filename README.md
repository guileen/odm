# golang ODM(Object Document Mapper) abstract layer

1. 首先支持 DynamoDB，按DynamoDB的习惯来操作。
2. 后续支持 MongoDB

## 建立连接池

TODO 连接池相关。

早期业务层自己管理连接

### 获取连接
如果需要在一条独立的连接上执行操作，则需要使用连接池，使用结束后，需要释放连接回到连接池。
pool.getConn()
conn.Close()

## Scheme 操作
TODO 根据Model定义生成表
对表的创建、建立索引由运维手动完成。暂不由代码控制。

## Table 接口

Table 用于抽象数据表的操作，对应Dynamo的Table，MongoDB的Collection，MySQL的Table。操作以Dynamo为基础进行设计。

TODO connectionPool.GetTable("User")

早期的Table按如下方式构建

```
table := &odm.DynamoTable{
    Connection: &mydynamodb,
    TableName: "User",
}
```

### Key 类型

Key是DynamoDB的概念，可以由1个或2个字段构成。使用一个map结构来表达。

`type Key map[string]interface{}`

### Model 类型

`type Model interface{}` 仅仅是一个指针，可以是任何结构体。

### 操作Options 类型

NameParams 对应 ExpressionAttributeNames
ValueParams 对应 ExpressionAttributeValues
```
type Condition struct {
	ConditionExpression       *string
    NameParams
    ValueParams
}
```
Condition 类型是一个条件表达式，仅当表达式成立时，操作才能成功。

### PutItem(key Key, cond Condition, item Model) error
PutItem 操作。替换整个item。

### UpdateItem(key Key, updateExpression string, opt UpdateOption, item Model) error
Update 部分字段，根据ReturnValues返回数据到item中。

### GetItem(key Key, consistentRead bool, item Model) error
consistentRead 代表是否是一致性读。

### DeleteItem(key Key, opt DeleteOption, item Model) error
被删除对象将填充到item。

### Query(startKey Key, QueryOption, items []Model) error
查询列表。
startKey 用来作性能优化。查询将从startKey开始。查询完成后，startKey将被更新。

```
type QueryOption struct {
	ConsistentRead            *bool                      `type:"boolean"`
	NameParams                map[string]*string         `type:"map"`
	ValueParams               map[string]*AttributeValue `type:"map"`
	FilterExpression          *string                    `type:"string"`
	IndexName                 *string                    `min:"3" type:"string"`
	KeyConditionExpression    *string                    `type:"string"`
	Limit                     *int64                     `min:"1" type:"integer"`
	ProjectionExpression      *string                    `type:"string"`
	ScanIndexForward          *bool                      `type:"boolean"`
	Select                    *string                    `type:"string" enum:"Select"`
}
```
	
### TODO Scan(startKey M, ScanOption, items []Model) error
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
