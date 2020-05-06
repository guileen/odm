package odm

import (
	"errors"
)

func (db *ODMDB) Transact() *transaction {
	return &transaction{
		db: db,
	}
}

type transaction struct {
	db               *ODMDB
	getTransaction   *getTransaction
	writeTransaction *writeTransaction
}

func (t *transaction) makeGet() *getTransaction {
	if t.writeTransaction != nil {
		panic("Transaction is initialized for write only")
	}
	if t.getTransaction != nil {
		return t.getTransaction
	}
	t.getTransaction = &getTransaction{
		db:         t.db,
		operations: []*TransactGet{},
	}
	return t.getTransaction
}

func (t *transaction) makeWrite() *writeTransaction {
	if t.getTransaction != nil {
		panic("Transaction is initialized for read only")
	}
	if t.writeTransaction != nil {
		return t.writeTransaction
	}
	t.writeTransaction = &writeTransaction{
		db:         t.db,
		operations: []*TransactWrite{},
	}
	return t.writeTransaction
}

func (t *transaction) Commit() error {
	if t.getTransaction != nil {
		return t.getTransaction.Commit()
	}
	if t.writeTransaction != nil {
		return t.writeTransaction.Commit()
	}
	return errors.New("Nothing to commit")
}

func (t *transaction) Get(args ...interface{}) *getTransaction {
	return t.makeGet().Get(args...)
}
func (t *transaction) Check(cond *WriteOption) *writeTransaction {
	return t.makeWrite().Check(cond)
}
func (t *transaction) Put(item Model, cond *WriteOption, failResult Model) *writeTransaction {
	return t.makeWrite().Put(item, cond, failResult)
}
func (t *transaction) Update(tableName string, hashKey interface{}, rangeKey interface{}, updateExpr string, opt *WriteOption, failValue Model) *writeTransaction {
	return t.makeWrite().Update(tableName, hashKey, rangeKey, updateExpr, opt, failValue)
}
func (t *transaction) Delete(tableName string, hashKey interface{}, rangeKey interface{}, opt *WriteOption, failValue Model) *writeTransaction {
	return t.makeWrite().Delete(tableName, hashKey, rangeKey, opt, failValue)
}

type writeTransaction struct {
	db         *ODMDB
	operations []*TransactWrite
}

func (t *writeTransaction) Check(cond *WriteOption) *writeTransaction {
	panic("not implemented")
	return nil
}

func (t *writeTransaction) Put(item interface{}, cond *WriteOption, failValue interface{}) *writeTransaction {
	meta := GetModelMeta(item)
	t.operations = append(t.operations, &TransactWrite{
		Put: &Put{
			TableName: meta.TableName,
			Item:      item,
		},
	})
	panic("not implemented")
	return t
}
func (t *writeTransaction) Update(tableName string, hashKey interface{}, rangeKey interface{}, updateExpr string, opt *WriteOption, failValue Model) *writeTransaction {
	op := &TransactWrite{
		Update: &Update{
			TableName:   tableName,
			HashKey:     hashKey,
			RangeKey:    rangeKey,
			Expression:  updateExpr,
			WriteOption: opt,
		},
	}
	if failValue != nil {
		op.Update.ReturnValuesOnConditionCheckFailure = "ALL_NEW"
	}
	t.operations = append(t.operations, op)
	return t
}
func (t *writeTransaction) Delete(args ...interface{}) *writeTransaction {
	t.operations = append(t.operations)
	panic("not implemented")
	return t
}
func (t *writeTransaction) Commit() error {
	return t.db.TransactWriteItems(t.operations)
}

type getTransaction struct {
	db         *ODMDB
	operations []*TransactGet
}

func (t *getTransaction) Get(args ...interface{}) *getTransaction {
	panic("not implemented")
	return t
}

func (t *getTransaction) Commit() error {
	panic("not implemented")
	return nil
}
