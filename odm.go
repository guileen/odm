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

// Update 更新一个对象，确保对象已经有主键
func Update(model types.Model) types.AsyncAction {
	return util.RunAsync(func() error {
		return GetAccessor(model).Update(model)
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
