// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TrManager is an autogenerated mock type for the TrManager type
type TrManager struct {
	mock.Mock
}

type TrManager_Expecter struct {
	mock *mock.Mock
}

func (_m *TrManager) EXPECT() *TrManager_Expecter {
	return &TrManager_Expecter{mock: &_m.Mock}
}

// Do provides a mock function with given fields: ctx, fn
func (_m *TrManager) Do(ctx context.Context, fn func(context.Context) error) error {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for Do")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TrManager_Do_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Do'
type TrManager_Do_Call struct {
	*mock.Call
}

// Do is a helper method to define mock.On call
//   - ctx context.Context
//   - fn func(context.Context) error
func (_e *TrManager_Expecter) Do(ctx interface{}, fn interface{}) *TrManager_Do_Call {
	return &TrManager_Do_Call{Call: _e.mock.On("Do", ctx, fn)}
}

func (_c *TrManager_Do_Call) Run(run func(ctx context.Context, fn func(context.Context) error)) *TrManager_Do_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error))
	})
	return _c
}

func (_c *TrManager_Do_Call) Return(err error) *TrManager_Do_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *TrManager_Do_Call) RunAndReturn(run func(context.Context, func(context.Context) error) error) *TrManager_Do_Call {
	_c.Call.Return(run)
	return _c
}

// NewTrManager creates a new instance of TrManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTrManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *TrManager {
	mock := &TrManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
