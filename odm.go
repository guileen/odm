package odm

import (
	"git.devops.com/go/odm/types"
	"git.devops.com/go/odm/util"
)

// GetAccessor 负责建立获取 Accessor，内部封装了建立连接相关的逻辑
func GetAccessor(model types.Model) types.Accessor {
	// 区分model和models
	return nil
}

// Create 新增一个对象
func Create(model types.Model) types.AsyncAction {
	return util.RunAsync(func() error {
		return GetAccessor(model).Create(model)
	})
}

// Update 更新一个对象，确保对象已经有主键
func Update(model types.Model) types.AsyncAction {
	return util.RunAsync(func() error {
		return GetAccessor(model).Update(model)
	})
}

// Save is create or update
func Save(model types.Model) types.AsyncAction {
	return util.RunAsync(func() error {
		return GetAccessor(model).Save(model)
	})
}

// GetByPK 针对 只有分区键 的访问,
func GetByPK(pk types.PK, model types.Model) types.AsyncAction {
	return util.RunAsync(func() error {
		return GetAccessor(model).GetByPK(pk, model)
	})
}

// GetByPKSK 针对 分区键、排序键 的访问
func GetByPKSK(pk types.PK, sk types.SK, model types.Model) types.AsyncAction {
	return util.RunAsync(func() error {
		return GetAccessor(model).GetByPKSK(pk, sk, model)
	})
}

// BatchGetByPKs 用于批量获取一组对象
func BatchGetByPKs(pks []types.PK, models []types.Model) types.AsyncAction {
	return util.RunAsync(func() error {
		return GetAccessor(models).BatchGetByPKs(pks, models)
	})
}

// BatchGetByPKSKs 用于批量获取一组对象
func BatchGetByPKSKs(pksks []types.PKSK, models []types.Model) types.AsyncAction {
	return util.RunAsync(func() error {
		return GetAccessor(models).BatchGetByPKSKs(pksks, models)
	})
}

// Find 进行除主键外的查询操作
func Find(expression string, params types.M, models []types.Model) types.AsyncAction {
	return util.RunAsync(func() error {
		return GetAccessor(models).Find(expression, params, models)
	})
}
