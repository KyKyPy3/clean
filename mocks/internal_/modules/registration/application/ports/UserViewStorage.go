// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/KyKyPy3/clean/internal/domain/common"

	entity "github.com/KyKyPy3/clean/internal/modules/user/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// UserViewStorage is an autogenerated mock type for the UserViewStorage type
type UserViewStorage struct {
	mock.Mock
}

type UserViewStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *UserViewStorage) EXPECT() *UserViewStorage_Expecter {
	return &UserViewStorage_Expecter{mock: &_m.Mock}
}

// Fetch provides a mock function with given fields: ctx, limit, offset
func (_m *UserViewStorage) Fetch(ctx context.Context, limit int64, offset int64) ([]entity.User, error) {
	ret := _m.Called(ctx, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for Fetch")
	}

	var r0 []entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) ([]entity.User, error)); ok {
		return rf(ctx, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) []entity.User); ok {
		r0 = rf(ctx, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserViewStorage_Fetch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fetch'
type UserViewStorage_Fetch_Call struct {
	*mock.Call
}

// Fetch is a helper method to define mock.On call
//   - ctx context.Context
//   - limit int64
//   - offset int64
func (_e *UserViewStorage_Expecter) Fetch(ctx interface{}, limit interface{}, offset interface{}) *UserViewStorage_Fetch_Call {
	return &UserViewStorage_Fetch_Call{Call: _e.mock.On("Fetch", ctx, limit, offset)}
}

func (_c *UserViewStorage_Fetch_Call) Run(run func(ctx context.Context, limit int64, offset int64)) *UserViewStorage_Fetch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(int64))
	})
	return _c
}

func (_c *UserViewStorage_Fetch_Call) Return(_a0 []entity.User, _a1 error) *UserViewStorage_Fetch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserViewStorage_Fetch_Call) RunAndReturn(run func(context.Context, int64, int64) ([]entity.User, error)) *UserViewStorage_Fetch_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *UserViewStorage) GetByEmail(ctx context.Context, email common.Email) (entity.User, error) {
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

// UserViewStorage_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type UserViewStorage_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email common.Email
func (_e *UserViewStorage_Expecter) GetByEmail(ctx interface{}, email interface{}) *UserViewStorage_GetByEmail_Call {
	return &UserViewStorage_GetByEmail_Call{Call: _e.mock.On("GetByEmail", ctx, email)}
}

func (_c *UserViewStorage_GetByEmail_Call) Run(run func(ctx context.Context, email common.Email)) *UserViewStorage_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Email))
	})
	return _c
}

func (_c *UserViewStorage_GetByEmail_Call) Return(_a0 entity.User, _a1 error) *UserViewStorage_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserViewStorage_GetByEmail_Call) RunAndReturn(run func(context.Context, common.Email) (entity.User, error)) *UserViewStorage_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *UserViewStorage) GetByID(ctx context.Context, id common.UID) (entity.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.UID) (entity.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.UID) entity.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.UID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserViewStorage_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type UserViewStorage_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id common.UID
func (_e *UserViewStorage_Expecter) GetByID(ctx interface{}, id interface{}) *UserViewStorage_GetByID_Call {
	return &UserViewStorage_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *UserViewStorage_GetByID_Call) Run(run func(ctx context.Context, id common.UID)) *UserViewStorage_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.UID))
	})
	return _c
}

func (_c *UserViewStorage_GetByID_Call) Return(_a0 entity.User, _a1 error) *UserViewStorage_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserViewStorage_GetByID_Call) RunAndReturn(run func(context.Context, common.UID) (entity.User, error)) *UserViewStorage_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserViewStorage creates a new instance of UserViewStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserViewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserViewStorage {
	mock := &UserViewStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}