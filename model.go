package odm

type Map map[string]interface{}

type Model interface {
}

type TableMeta struct {
	TableName    string
	PartitionKey string
	SortingKey   string
}

type TableConfig struct {
	Name     string
	UseCache bool
	TTL      int64
}

type TableConfigGetter interface {
	TableConfig() *TableConfig
}
