// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	chains "github.com/zeta-chain/node/pkg/chains"

	types "github.com/cosmos/cosmos-sdk/types"
)

// CrosschainAuthorityKeeper is an autogenerated mock type for the CrosschainAuthorityKeeper type
type CrosschainAuthorityKeeper struct {
	mock.Mock
}

// CheckAuthorization provides a mock function with given fields: ctx, msg
func (_m *CrosschainAuthorityKeeper) CheckAuthorization(ctx types.Context, msg types.Msg) error {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for CheckAuthorization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, types.Msg) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAdditionalChainList provides a mock function with given fields: ctx
func (_m *CrosschainAuthorityKeeper) GetAdditionalChainList(ctx types.Context) []chains.Chain {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAdditionalChainList")
	}

	var r0 []chains.Chain
	if rf, ok := ret.Get(0).(func(types.Context) []chains.Chain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]chains.Chain)
		}
	}

	return r0
}

// NewCrosschainAuthorityKeeper creates a new instance of CrosschainAuthorityKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCrosschainAuthorityKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *CrosschainAuthorityKeeper {
	mock := &CrosschainAuthorityKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
