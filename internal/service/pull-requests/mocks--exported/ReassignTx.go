package mocks

import (
	context "context"

	domain "github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

type ReassignTx struct {
	mock.Mock
}

type ReassignTx_Expecter struct {
	mock *mock.Mock
}

func (_m *ReassignTx) EXPECT() *ReassignTx_Expecter {
	return &ReassignTx_Expecter{mock: &_m.Mock}
}

func (_m *ReassignTx) AssignMember(_a0 context.Context, _a1 domain.PrId, _a2 domain.MemberId) (domain.PullRequest, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for AssignMember")
	}

	var r0 domain.PullRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.PrId, domain.MemberId) (domain.PullRequest, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.PrId, domain.MemberId) domain.PullRequest); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(domain.PullRequest)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.PrId, domain.MemberId) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type ReassignTx_AssignMember_Call struct {
	*mock.Call
}

func (_e *ReassignTx_Expecter) AssignMember(_a0 interface{}, _a1 interface{}, _a2 interface{}) *ReassignTx_AssignMember_Call {
	return &ReassignTx_AssignMember_Call{Call: _e.mock.On("AssignMember", _a0, _a1, _a2)}
}

func (_c *ReassignTx_AssignMember_Call) Run(run func(_a0 context.Context, _a1 domain.PrId, _a2 domain.MemberId)) *ReassignTx_AssignMember_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.PrId), args[2].(domain.MemberId))
	})
	return _c
}

func (_c *ReassignTx_AssignMember_Call) Return(_a0 domain.PullRequest, _a1 error) *ReassignTx_AssignMember_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ReassignTx_AssignMember_Call) RunAndReturn(run func(context.Context, domain.PrId, domain.MemberId) (domain.PullRequest, error)) *ReassignTx_AssignMember_Call {
	_c.Call.Return(run)
	return _c
}

func (_m *ReassignTx) Commit() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type ReassignTx_Commit_Call struct {
	*mock.Call
}

func (_e *ReassignTx_Expecter) Commit() *ReassignTx_Commit_Call {
	return &ReassignTx_Commit_Call{Call: _e.mock.On("Commit")}
}

func (_c *ReassignTx_Commit_Call) Run(run func()) *ReassignTx_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReassignTx_Commit_Call) Return(_a0 error) *ReassignTx_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ReassignTx_Commit_Call) RunAndReturn(run func() error) *ReassignTx_Commit_Call {
	_c.Call.Return(run)
	return _c
}

func (_m *ReassignTx) GetPullRequestMembersHistories(_a0 context.Context, _a1 domain.PrId) (domain.MembersHistories, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetPullRequestMembersHistories")
	}

	var r0 domain.MembersHistories
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.PrId) (domain.MembersHistories, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.PrId) domain.MembersHistories); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.MembersHistories)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.PrId) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type ReassignTx_GetPullRequestMembersHistories_Call struct {
	*mock.Call
}

func (_e *ReassignTx_Expecter) GetPullRequestMembersHistories(_a0 interface{}, _a1 interface{}) *ReassignTx_GetPullRequestMembersHistories_Call {
	return &ReassignTx_GetPullRequestMembersHistories_Call{Call: _e.mock.On("GetPullRequestMembersHistories", _a0, _a1)}
}

func (_c *ReassignTx_GetPullRequestMembersHistories_Call) Run(run func(_a0 context.Context, _a1 domain.PrId)) *ReassignTx_GetPullRequestMembersHistories_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.PrId))
	})
	return _c
}

func (_c *ReassignTx_GetPullRequestMembersHistories_Call) Return(_a0 domain.MembersHistories, _a1 error) *ReassignTx_GetPullRequestMembersHistories_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ReassignTx_GetPullRequestMembersHistories_Call) RunAndReturn(run func(context.Context, domain.PrId) (domain.MembersHistories, error)) *ReassignTx_GetPullRequestMembersHistories_Call {
	_c.Call.Return(run)
	return _c
}

func (_m *ReassignTx) Rollback() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Rollback")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type ReassignTx_Rollback_Call struct {
	*mock.Call
}

func (_e *ReassignTx_Expecter) Rollback() *ReassignTx_Rollback_Call {
	return &ReassignTx_Rollback_Call{Call: _e.mock.On("Rollback")}
}

func (_c *ReassignTx_Rollback_Call) Run(run func()) *ReassignTx_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReassignTx_Rollback_Call) Return(_a0 error) *ReassignTx_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ReassignTx_Rollback_Call) RunAndReturn(run func() error) *ReassignTx_Rollback_Call {
	_c.Call.Return(run)
	return _c
}

func NewReassignTx(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReassignTx {
	mock := &ReassignTx{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
