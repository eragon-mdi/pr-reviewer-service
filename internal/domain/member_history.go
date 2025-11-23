package domain

type MembersHistories []MemberHistory

type MemberHistory struct {
	Id                MemberId
	Status            MemberStatus
	Role              MemberRole
	wasAssignedBefore bool
}

func (mh MembersHistories) Slice() []MemberHistory {
	return []MemberHistory(mh)
}

func (mh MembersHistories) Empty() bool {
	return len(mh) == 0
}

func NewMemberHistory(id MemberId, status MemberStatus, role MemberRole, wasAssignedBefore bool) MemberHistory {
	return MemberHistory{
		Id:                id,
		Status:            status,
		Role:              role,
		wasAssignedBefore: wasAssignedBefore,
	}
}
