package singleton

import (
	"sync"
)

type Singleton[T, U any] interface {
	Client() *T
	Init() error
	InitFromConfig(conf U) error
}

type BaseInstance[T, U any] struct {
	once       sync.Once
	initFn     func() error
	initConfFn func(conf U) error
}

func NewBaseInstance[T, U any](initFn func() error, initConfFn func(conf U) error) *BaseInstance[T, U] {
	i := new(BaseInstance[T, U])
	i.initFn = initFn
	i.initConfFn = initConfFn

	return i
}

func (i *BaseInstance[T, U]) Client() *T {
	return nil
}

func (i *BaseInstance[T, U]) Init() error {
	var initErr error

	i.once.Do(func() {
		initErr = i.initFn()
	})

	return initErr
}

func (i *BaseInstance[T, U]) InitFromConfig(conf U) error {
	var initErr error

	i.once.Do(func() {
		initErr = i.initConfFn(conf)
	})

	return initErr
}
