// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	core "github.com/KyKyPy3/clean/internal/application/core"
	mock "github.com/stretchr/testify/mock"
)

// Command is an autogenerated mock type for the Command type
type Command struct {
	mock.Mock
}

type Command_Expecter struct {
	mock *mock.Mock
}

func (_m *Command) EXPECT() *Command_Expecter {
	return &Command_Expecter{mock: &_m.Mock}
}

// Type provides a mock function with given fields:
func (_m *Command) Type() core.CommandType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Type")
	}

	var r0 core.CommandType
	if rf, ok := ret.Get(0).(func() core.CommandType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(core.CommandType)
	}

	return r0
}

// Command_Type_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Type'
type Command_Type_Call struct {
	*mock.Call
}

// Type is a helper method to define mock.On call
func (_e *Command_Expecter) Type() *Command_Type_Call {
	return &Command_Type_Call{Call: _e.mock.On("Type")}
}

func (_c *Command_Type_Call) Run(run func()) *Command_Type_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Command_Type_Call) Return(_a0 core.CommandType) *Command_Type_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Command_Type_Call) RunAndReturn(run func() core.CommandType) *Command_Type_Call {
	_c.Call.Return(run)
	return _c
}

// NewCommand creates a new instance of Command. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommand(t interface {
	mock.TestingT
	Cleanup(func())
}) *Command {
	mock := &Command{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
