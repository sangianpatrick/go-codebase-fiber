// Code generated by mockery v2.39.2. DO NOT EDIT.

package mock_customer

import (
	context "context"

	request "github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/request"
	response "github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/response"
	mock "github.com/stretchr/testify/mock"
)

// CustomerUseCase is an autogenerated mock type for the CustomerUseCase type
type CustomerUseCase struct {
	mock.Mock
}

// SignUp provides a mock function with given fields: ctx, req
func (_m *CustomerUseCase) SignUp(ctx context.Context, req request.SignUpRequest) (response.SignUpResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for SignUp")
	}

	var r0 response.SignUpResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.SignUpRequest) (response.SignUpResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.SignUpRequest) response.SignUpResponse); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(response.SignUpResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.SignUpRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCustomerUseCase creates a new instance of CustomerUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCustomerUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *CustomerUseCase {
	mock := &CustomerUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
