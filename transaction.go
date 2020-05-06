package odm

type transaction struct {
	db *ODMDB
}

func (t *transaction) makeGet() *getTransaction {
	return &getTransaction{
		db:         t.db,
		operations: []*TransactGet{},
	}
}

func (t *transaction) makeWrite() *writeTransaction {
	return &writeTransaction{
		db:         t.db,
		operations: []*TransactWrite{},
	}
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
func (t *transaction) Update(tableName string, partitionKey interface{}, sortingKey interface{}, updateExpr string, opt *WriteOption, failValue Model) *writeTransaction {
	return t.makeWrite().Update(tableName, partitionKey, sortingKey, updateExpr, opt, failValue)
}
func (t *transaction) Delete(tableName string, partitionKey interface{}, sortingKey interface{}, opt *WriteOption, failValue Model) *writeTransaction {
	return t.makeWrite().Delete(tableName, partitionKey, sortingKey, opt, failValue)
}

type writeTransaction struct {
	db         *ODMDB
	operations []*TransactWrite
}

func (t *writeTransaction) Check(cond *WriteOption) *writeTransaction {
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
	return t
}
func (t *writeTransaction) Update(tableName string, partitionKey interface{}, sortingKey interface{}, updateExpr string, opt *WriteOption, failValue Model) *writeTransaction {
	return t
}
func (t *writeTransaction) Delete(args ...interface{}) *writeTransaction {
	t.operations = append(t.operations)
	return t
}
func (t *writeTransaction) Commit() error {
	return nil
}

type getTransaction struct {
	db         *ODMDB
	operations []*TransactGet
}

func (t *getTransaction) Get(args ...interface{}) *getTransaction {
	return t
}

func (t *getTransaction) Commit() error {
	return nil
}
