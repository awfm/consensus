// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	model "github.com/alvalor/consensus/model"
	mock "github.com/stretchr/testify/mock"
)

// Database is an autogenerated mock type for the Database type
type Database struct {
	mock.Mock
}

// Block provides a mock function with given fields: blockID
func (_m *Database) Block(blockID model.Hash) (*model.Block, error) {
	ret := _m.Called(blockID)

	var r0 *model.Block
	if rf, ok := ret.Get(0).(func(model.Hash) *model.Block); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Block)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.Hash) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}