// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	auth "auth/internal/services/auth"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// RefreshTokenStorage is an autogenerated mock type for the RefreshTokenStorage type
type RefreshTokenStorage struct {
	mock.Mock
}

// Create provides a mock function with given fields: token
func (_m *RefreshTokenStorage) Create(token *auth.RefreshToken) (uuid.UUID, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(*auth.RefreshToken) (uuid.UUID, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(*auth.RefreshToken) uuid.UUID); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(*auth.RefreshToken) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *RefreshTokenStorage) Delete(id uuid.UUID) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *RefreshTokenStorage) Get(id uuid.UUID) (*auth.RefreshToken, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *auth.RefreshToken
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*auth.RefreshToken, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *auth.RefreshToken); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.RefreshToken)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRefreshTokenStorage creates a new instance of RefreshTokenStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRefreshTokenStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *RefreshTokenStorage {
	mock := &RefreshTokenStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
