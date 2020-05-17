// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	base "github.com/awfm/consensus/model/base"

	message "github.com/awfm/consensus/model/message"

	mock "github.com/stretchr/testify/mock"
)

// Signer is an autogenerated mock type for the Signer type
type Signer struct {
	mock.Mock
}

// Proposal provides a mock function with given fields: vertex
func (_m *Signer) Proposal(vertex *base.Vertex) (*message.Proposal, error) {
	ret := _m.Called(vertex)

	var r0 *message.Proposal
	if rf, ok := ret.Get(0).(func(*base.Vertex) *message.Proposal); ok {
		r0 = rf(vertex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*message.Proposal)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*base.Vertex) error); ok {
		r1 = rf(vertex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Self provides a mock function with given fields:
func (_m *Signer) Self() (base.Hash, error) {
	ret := _m.Called()

	var r0 base.Hash
	if rf, ok := ret.Get(0).(func() base.Hash); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(base.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Vote provides a mock function with given fields: vertex
func (_m *Signer) Vote(vertex *base.Vertex) (*message.Vote, error) {
	ret := _m.Called(vertex)

	var r0 *message.Vote
	if rf, ok := ret.Get(0).(func(*base.Vertex) *message.Vote); ok {
		r0 = rf(vertex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*message.Vote)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*base.Vertex) error); ok {
		r1 = rf(vertex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
