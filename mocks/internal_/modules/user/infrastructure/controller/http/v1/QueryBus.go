// Code generated by mockery v2.49.0. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/KyKyPy3/clean/internal/application/core"
	mock "github.com/stretchr/testify/mock"
)

// QueryBus is an autogenerated mock type for the QueryBus type
type QueryBus struct {
	mock.Mock
}

type QueryBus_Expecter struct {
	mock *mock.Mock
}

func (_m *QueryBus) EXPECT() *QueryBus_Expecter {
	return &QueryBus_Expecter{mock: &_m.Mock}
}

// Ask provides a mock function with given fields: _a0, _a1
func (_m *QueryBus) Ask(_a0 context.Context, _a1 core.Query) (any, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Ask")
	}

	var r0 any
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, core.Query) (any, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.Query) any); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(any)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.Query) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryBus_Ask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ask'
type QueryBus_Ask_Call struct {
	*mock.Call
}

// Ask is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 core.Query
func (_e *QueryBus_Expecter) Ask(_a0 interface{}, _a1 interface{}) *QueryBus_Ask_Call {
	return &QueryBus_Ask_Call{Call: _e.mock.On("Ask", _a0, _a1)}
}

func (_c *QueryBus_Ask_Call) Run(run func(_a0 context.Context, _a1 core.Query)) *QueryBus_Ask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.Query))
	})
	return _c
}

func (_c *QueryBus_Ask_Call) Return(_a0 any, _a1 error) *QueryBus_Ask_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *QueryBus_Ask_Call) RunAndReturn(run func(context.Context, core.Query) (any, error)) *QueryBus_Ask_Call {
	_c.Call.Return(run)
	return _c
}

// NewQueryBus creates a new instance of QueryBus. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQueryBus(t interface {
	mock.TestingT
	Cleanup(func())
}) *QueryBus {
	mock := &QueryBus{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
