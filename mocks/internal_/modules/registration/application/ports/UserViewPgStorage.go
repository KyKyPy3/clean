// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/KyKyPy3/clean/internal/domain/common"

	entity "github.com/KyKyPy3/clean/internal/modules/user/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// UserViewPgStorage is an autogenerated mock type for the UserViewPgStorage type
type UserViewPgStorage struct {
	mock.Mock
}

type UserViewPgStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *UserViewPgStorage) EXPECT() *UserViewPgStorage_Expecter {
	return &UserViewPgStorage_Expecter{mock: &_m.Mock}
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *UserViewPgStorage) GetByEmail(ctx context.Context, email common.Email) (entity.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetByEmail")
	}

	var r0 entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Email) (entity.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Email) entity.User); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Email) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserViewPgStorage_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type UserViewPgStorage_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email common.Email
func (_e *UserViewPgStorage_Expecter) GetByEmail(ctx interface{}, email interface{}) *UserViewPgStorage_GetByEmail_Call {
	return &UserViewPgStorage_GetByEmail_Call{Call: _e.mock.On("GetByEmail", ctx, email)}
}

func (_c *UserViewPgStorage_GetByEmail_Call) Run(run func(ctx context.Context, email common.Email)) *UserViewPgStorage_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Email))
	})
	return _c
}

func (_c *UserViewPgStorage_GetByEmail_Call) Return(_a0 entity.User, _a1 error) *UserViewPgStorage_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserViewPgStorage_GetByEmail_Call) RunAndReturn(run func(context.Context, common.Email) (entity.User, error)) *UserViewPgStorage_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserViewPgStorage creates a new instance of UserViewPgStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserViewPgStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserViewPgStorage {
	mock := &UserViewPgStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
