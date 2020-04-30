package types

type AsyncAction interface {
	WaitDone() error
}

type AsyncResult interface {
	Result() (interface{}, error)
}

type FindResult interface {
	One(model Model) error
	List(models []Model) error
}
