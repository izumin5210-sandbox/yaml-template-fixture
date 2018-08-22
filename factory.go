package tfx

import "testing"

type Factory interface {
	Load(tb testing.TB, name string, v interface{}, opts ...LoadOption)
}

var DefaultFactory Factory = New()

func Load(tb testing.TB, name string, v interface{}, opts ...LoadOption) {
	tb.Helper()
	DefaultFactory.Load(tb, name, v, opts...)
}

func New() Factory {
	return &factoryImpl{}
}

type factoryImpl struct {
}

func (f *factoryImpl) Load(tb testing.TB, name string, v interface{}, opts ...LoadOption) {
	tb.Helper()
	// TODO
}
