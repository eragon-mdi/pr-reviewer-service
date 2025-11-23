package restpullrequests

import (
	"time"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
)

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" validate:"required,uuid"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorID        string `json:"author_id" validate:"required,uuid"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required,uuid"`
}

type ReassignPRRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required,uuid"`
	OldUserID     string `json:"old_reviewer_id" validate:"required,uuid"`
}

type PRResponse struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type ReassignPRResponse struct {
	PR         PRResponse `json:"pr"`
	ReplacedBy string     `json:"replaced_by"`
}

func (req *CreatePRRequest) domain() domain.PullRequestShort {
	return domain.PullRequestShort{
		Id:       domain.PrId(req.PullRequestID),
		Name:     domain.PrName(req.PullRequestName),
		AuthorId: domain.MemberId(req.AuthorID),
		Status:   domain.PrStatusDefault,
	}
}

func pullRequestResponse(pr domain.PullRequest) PRResponse {
	status := "OPEN"
	if pr.Status == domain.PrStatusMerged {
		status = "MERGED"
	}

	assignedReviewers := make([]string, 0, len(pr.AssignedReviews))
	for _, member := range pr.AssignedReviews.Slice() {
		assignedReviewers = append(assignedReviewers, member.Id.String())
	}

	var createdAt *string
	if !pr.CreatedAt.IsZero() {
		createdAtStr := pr.CreatedAt.Format(time.RFC3339)
		createdAt = &createdAtStr
	}

	var mergedAt *string
	if !pr.MergedAt.IsZero() {
		mergedAtStr := pr.MergedAt.Format(time.RFC3339)
		mergedAt = &mergedAtStr
	}

	return PRResponse{
		PullRequestID:     pr.Id.String(),
		PullRequestName:   pr.Name.String(),
		AuthorID:          pr.AuthorId.String(),
		Status:            status,
		AssignedReviewers: assignedReviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
}

func reassignPRResponse(pr domain.PrWithReasignMember) ReassignPRResponse {
	return ReassignPRResponse{
		PR:         pullRequestResponse(pr.PullRequest),
		ReplacedBy: string(pr.MemberId),
	}
}
