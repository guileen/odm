package types

import "git.devops.com/go/odm"

type AsyncAction interface {
	WaitDone() error
}

type AsyncResult interface {
	Result() (interface{}, error)
}

type FindResult interface {
	One(model odm.Model) error
	List(models []odm.Model) error
}
