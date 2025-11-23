package domain

import "github.com/google/uuid"

const MemberStatusDefault = MemberStatusActive

const (
	MemberStatusActive = iota
	MemberStatusInactive
)

var MemberStatusNames = map[MemberStatus]string{
	MemberStatusActive:   "active",
	MemberStatusInactive: "inactive",
}

type MemberId string
type MemberStatus int

type Member struct {
	Id      MemberId
	Name    string
	Status  MemberStatus
	Reviews PullRequests
	Team    TeamName
}

func (mId MemberId) IsValid() bool {
	return uuid.Validate(string(mId)) == nil
}

func (mId MemberId) String() string {
	return string(mId)
}

type memberBuilder interface {
	Name(string) memberBuilder
	Status(MemberStatus) memberBuilder
	Reviews([]PullRequestShort) memberBuilder
	Build() Member
}

type memBuilder struct {
	m Member
}

func MemberBuilder(id MemberId) memberBuilder {
	return &memBuilder{
		m: Member{
			Id: id,
		},
	}
}

func (mb *memBuilder) Name(name string) memberBuilder {
	mb.m.Name = name
	return mb
}

func (mb *memBuilder) Status(st MemberStatus) memberBuilder {
	mb.m.Status = st
	return mb
}

func (mb *memBuilder) Reviews(rs []PullRequestShort) memberBuilder {
	mb.m.Reviews = PullRequests(rs)
	return mb
}

func (mb *memBuilder) Build() Member {
	return mb.m
}

type Members []Member

func (tm Members) Empty() bool {
	return len(tm) == 0
}
func (tm Members) Slice() []Member {
	return []Member(tm)
}

func (s MemberStatus) String() string {
	if name, ok := MemberStatusNames[s]; ok {
		return name
	}
	return "unknown"
}

func (s MemberStatus) IsActive() bool {
	return s == MemberStatusActive
}

func MemberStatusFromString(s string) MemberStatus {
	for k, v := range MemberStatusNames {
		if v == s {
			return k
		}
	}
	return MemberStatusDefault
}

func MemberStatusIsActiveByBool(isActive bool) MemberStatus {
	if isActive {
		return MemberStatusActive
	}
	return MemberStatusInactive
}

func MembersStatusesFromSliceOfStrings(slS []string) []MemberStatus {
	res := make([]MemberStatus, 0, len(slS))

	for _, s := range slS {
		res = append(res, MemberStatusFromString(s))
	}

	return res
}
