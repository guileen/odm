package types

import "git.devops.com/go/odm"

// PK 分区键
type PK interface{}

// SK 排序键
type SK interface{}

// PKSK 分区键、排序键（联合主键）
type PKSK struct {
	pk PK
	sk SK
}

// Accessor 是一个高级存储访问对象,针对某个表
type Accessor interface {
	// Create 插入一个对象?...
	Insert(model odm.Model) error

	// Update 更新一个对象
	Update(model odm.Model) error

	// UpdateByPK  PK是唯一主键。对于dynamodb，调用该方法的表只存在分区主键，不应该存在排序主键，参数是PartitionKey。
	// 对于mongodb，pk应该是_id。在数据迁移时，对于没有排序主键的表，应将PartionKey迁移为 _id，以兼容该函数。
	UpdateOne(pk PK, params odm.Map) error

	// UpdateMany ，cond是搜索条件， 对于 dynamodb， cond 可能是 分区主键+排序主键，对于mongodb 则是一个查询条件（需要建索引）
	// 但他们在形式上都是统一的  cond : {"UserId": "1234", "Time": "20190304"}
	UpdateMany(cond odm.Map, params odm.Map) error

	// GetByPK 针对 只有分区键 的访问
	FindOneByPK(pk PK, model odm.Model) error

	// GetByPKSK 针对 分区键、排序键 的访问
	FindOne(cond odm.Map, model odm.Model) error

	// GetByPKSK 针对 分区键、排序键 的访问
	FindMany(cond odm.Map, model []odm.Model) error

	// 根据主键删除一条数据
	DeleteOne(pk PK) error

	// 删除多个 条件应该传主键信息
	DeleteMany(cond odm.Map) error

	// BatchGetByPKs 针对 只有分区键 的批量访问
	BatchGetByPKs(pks []PK, models []odm.Model) error

	// BatchGetByPKSKs 针对 分区键、排序键 的批量访问
	BatchGetByPKSKs(pksks []PKSK, model []odm.Model) error

	// Find 针对 Index、Scan、Filter 类的操作
	Find(expression string, params odm.Map, models []odm.Model) error
}
