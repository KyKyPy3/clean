// Code generated by mockery v2.49.0. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/KyKyPy3/clean/internal/domain/common"

	entity "github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// RegistrationPgStorage is an autogenerated mock type for the RegistrationPgStorage type
type RegistrationPgStorage struct {
	mock.Mock
}

type RegistrationPgStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *RegistrationPgStorage) EXPECT() *RegistrationPgStorage_Expecter {
	return &RegistrationPgStorage_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, registration
func (_m *RegistrationPgStorage) Create(ctx context.Context, registration entity.Registration) error {
	ret := _m.Called(ctx, registration)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Registration) error); ok {
		r0 = rf(ctx, registration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegistrationPgStorage_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type RegistrationPgStorage_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - registration entity.Registration
func (_e *RegistrationPgStorage_Expecter) Create(ctx interface{}, registration interface{}) *RegistrationPgStorage_Create_Call {
	return &RegistrationPgStorage_Create_Call{Call: _e.mock.On("Create", ctx, registration)}
}

func (_c *RegistrationPgStorage_Create_Call) Run(run func(ctx context.Context, registration entity.Registration)) *RegistrationPgStorage_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.Registration))
	})
	return _c
}

func (_c *RegistrationPgStorage_Create_Call) Return(_a0 error) *RegistrationPgStorage_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RegistrationPgStorage_Create_Call) RunAndReturn(run func(context.Context, entity.Registration) error) *RegistrationPgStorage_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *RegistrationPgStorage) GetByID(ctx context.Context, id common.UID) (entity.Registration, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 entity.Registration
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.UID) (entity.Registration, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.UID) entity.Registration); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(entity.Registration)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.UID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegistrationPgStorage_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type RegistrationPgStorage_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id common.UID
func (_e *RegistrationPgStorage_Expecter) GetByID(ctx interface{}, id interface{}) *RegistrationPgStorage_GetByID_Call {
	return &RegistrationPgStorage_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *RegistrationPgStorage_GetByID_Call) Run(run func(ctx context.Context, id common.UID)) *RegistrationPgStorage_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.UID))
	})
	return _c
}

func (_c *RegistrationPgStorage_GetByID_Call) Return(_a0 entity.Registration, _a1 error) *RegistrationPgStorage_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RegistrationPgStorage_GetByID_Call) RunAndReturn(run func(context.Context, common.UID) (entity.Registration, error)) *RegistrationPgStorage_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, registration
func (_m *RegistrationPgStorage) Update(ctx context.Context, registration entity.Registration) error {
	ret := _m.Called(ctx, registration)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Registration) error); ok {
		r0 = rf(ctx, registration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegistrationPgStorage_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type RegistrationPgStorage_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - registration entity.Registration
func (_e *RegistrationPgStorage_Expecter) Update(ctx interface{}, registration interface{}) *RegistrationPgStorage_Update_Call {
	return &RegistrationPgStorage_Update_Call{Call: _e.mock.On("Update", ctx, registration)}
}

func (_c *RegistrationPgStorage_Update_Call) Run(run func(ctx context.Context, registration entity.Registration)) *RegistrationPgStorage_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.Registration))
	})
	return _c
}

func (_c *RegistrationPgStorage_Update_Call) Return(_a0 error) *RegistrationPgStorage_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RegistrationPgStorage_Update_Call) RunAndReturn(run func(context.Context, entity.Registration) error) *RegistrationPgStorage_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewRegistrationPgStorage creates a new instance of RegistrationPgStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRegistrationPgStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *RegistrationPgStorage {
	mock := &RegistrationPgStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
