package types

// Config is Connection Configuration.
type Config interface {
}

// Conn 连接
type DB interface {
	Table(model Model) Table
	GetTable(name string) Table
	Close()
}

// Pool 连接池
type Pool interface {
	DB()
}
