// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	common "github.com/KyKyPy3/clean/internal/domain/common"
	mock "github.com/stretchr/testify/mock"
)

// UniquenessPolicer is an autogenerated mock type for the UniquenessPolicer type
type UniquenessPolicer struct {
	mock.Mock
}

type UniquenessPolicer_Expecter struct {
	mock *mock.Mock
}

func (_m *UniquenessPolicer) EXPECT() *UniquenessPolicer_Expecter {
	return &UniquenessPolicer_Expecter{mock: &_m.Mock}
}

// IsUnique provides a mock function with given fields: email
func (_m *UniquenessPolicer) IsUnique(email common.Email) (bool, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for IsUnique")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Email) (bool, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(common.Email) bool); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(common.Email) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UniquenessPolicer_IsUnique_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsUnique'
type UniquenessPolicer_IsUnique_Call struct {
	*mock.Call
}

// IsUnique is a helper method to define mock.On call
//   - email common.Email
func (_e *UniquenessPolicer_Expecter) IsUnique(email interface{}) *UniquenessPolicer_IsUnique_Call {
	return &UniquenessPolicer_IsUnique_Call{Call: _e.mock.On("IsUnique", email)}
}

func (_c *UniquenessPolicer_IsUnique_Call) Run(run func(email common.Email)) *UniquenessPolicer_IsUnique_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.Email))
	})
	return _c
}

func (_c *UniquenessPolicer_IsUnique_Call) Return(_a0 bool, _a1 error) *UniquenessPolicer_IsUnique_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UniquenessPolicer_IsUnique_Call) RunAndReturn(run func(common.Email) (bool, error)) *UniquenessPolicer_IsUnique_Call {
	_c.Call.Return(run)
	return _c
}

// NewUniquenessPolicer creates a new instance of UniquenessPolicer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUniquenessPolicer(t interface {
	mock.TestingT
	Cleanup(func())
}) *UniquenessPolicer {
	mock := &UniquenessPolicer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
