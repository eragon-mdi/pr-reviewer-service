package domain

import "time"

type PrId string
type PrName string
type PrStatus int

const (
	PrStatusDefault = PrStatusOpen
	PrStatusOpen    = iota
	PrStatusMerged
)

type PullRequest struct {
	Id              PrId
	Name            PrName
	AuthorId        MemberId
	Status          PrStatus
	CreatedAt       time.Time
	MergedAt        time.Time
	version         int
	AssignedReviews Members
}

type PullRequests []PullRequestShort

type PullRequestShort struct {
	Id       PrId
	Name     PrName
	AuthorId MemberId
	Status   PrStatus
}

func (prs *PullRequestShort) Create() PullRequest {
	return PullRequest{
		Id:              prs.Id,
		Name:            prs.Name,
		AuthorId:        prs.AuthorId,
		Status:          PrStatusDefault,
		CreatedAt:       time.Now(),
		MergedAt:        time.Time{},
		version:         0,
		AssignedReviews: nil,
	}
}

func (pId PrId) String() string {
	return string(pId)
}

func (pName PrName) String() string {
	return string(pName)
}

func (mr PullRequests) Empty() bool {
	return len(mr) == 0
}

type PrReasignMember struct {
	PrId     PrId
	MemberId MemberId
}

type PrWithReasignMember struct {
	PullRequest
	MemberId MemberId
}
