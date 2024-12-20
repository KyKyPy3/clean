// Code generated by mockery v2.49.0. DO NOT EDIT.

package mocks

import (
	core "github.com/KyKyPy3/clean/internal/application/core"
	mock "github.com/stretchr/testify/mock"
)

// Query is an autogenerated mock type for the Query type
type Query struct {
	mock.Mock
}

type Query_Expecter struct {
	mock *mock.Mock
}

func (_m *Query) EXPECT() *Query_Expecter {
	return &Query_Expecter{mock: &_m.Mock}
}

// Type provides a mock function with given fields:
func (_m *Query) Type() core.QueryType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Type")
	}

	var r0 core.QueryType
	if rf, ok := ret.Get(0).(func() core.QueryType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(core.QueryType)
	}

	return r0
}

// Query_Type_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Type'
type Query_Type_Call struct {
	*mock.Call
}

// Type is a helper method to define mock.On call
func (_e *Query_Expecter) Type() *Query_Type_Call {
	return &Query_Type_Call{Call: _e.mock.On("Type")}
}

func (_c *Query_Type_Call) Run(run func()) *Query_Type_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Query_Type_Call) Return(_a0 core.QueryType) *Query_Type_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Query_Type_Call) RunAndReturn(run func() core.QueryType) *Query_Type_Call {
	_c.Call.Return(run)
	return _c
}

// NewQuery creates a new instance of Query. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQuery(t interface {
	mock.TestingT
	Cleanup(func())
}) *Query {
	mock := &Query{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
