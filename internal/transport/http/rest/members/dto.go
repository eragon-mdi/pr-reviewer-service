package restmembers

import "github.com/eragon-mdi/pr-reviewer-service/internal/domain"

type SetIsActiveRequest struct {
	UserID   string `json:"user_id" validate:"required,uuid"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserReviewsResponse struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}

type PullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func (req *SetIsActiveRequest) domain() domain.Member {
	return domain.MemberBuilder(domain.MemberId(req.UserID)).
		Status(domain.MemberStatusIsActiveByBool(req.IsActive)).
		Build()
}

func userResponse(m domain.Member) UserResponse {
	return UserResponse{
		UserID:   m.Id.String(),
		Username: m.Name,
		TeamName: m.Team.String(),
		IsActive: m.Status.IsActive(),
	}
}

func userReviewsResponse(m domain.Member) UserReviewsResponse {
	return UserReviewsResponse{
		UserID:       m.Id.String(),
		PullRequests: pullRequestShorts(m.Reviews),
	}
}

func pullRequestShort(pr domain.PullRequestShort) PullRequestShortDTO {
	status := "OPEN"
	if pr.Status == domain.PrStatusMerged {
		status = "MERGED"
	}
	return PullRequestShortDTO{
		PullRequestID:   pr.Id.String(),
		PullRequestName: pr.Name.String(),
		AuthorID:        pr.AuthorId.String(),
		Status:          status,
	}
}

func pullRequestShorts(prs domain.PullRequests) []PullRequestShortDTO {
	if prs.Empty() {
		return []PullRequestShortDTO{}
	}

	res := make([]PullRequestShortDTO, 0, len(prs))
	for _, pr := range prs {
		res = append(res, pullRequestShort(pr))
	}
	return res
}
