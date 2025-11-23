package restteams

import "github.com/eragon-mdi/pr-reviewer-service/internal/domain"

type TeamMember struct {
	UserID   string `json:"user_id" validate:"required,uuid"`
	Username string `json:"username" validate:"required"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type TeamRequest struct {
	TeamName string       `json:"team_name" validate:"required"`
	Members  []TeamMember `json:"members" validate:"required,dive,required"`
}

type TeamResponse struct {
	TeamName string       `json:"team_name" validate:"required"`
	Members  []TeamMember `json:"members" validate:"required,dive,required"`
}

func (tr *TeamRequest) domain() domain.Team {
	mems := make([]domain.Member, 0, len(tr.Members))
	for _, m := range tr.Members {
		mems = append(mems, m.domain())
	}
	return domain.NewTeam(domain.TeamName(tr.TeamName), mems...)
}

func (tr *TeamMember) domain() domain.Member {
	return domain.MemberBuilder(domain.MemberId(tr.UserID)).
		Name(tr.Username).
		Status(domain.MemberStatusIsActiveByBool(tr.IsActive)).
		Build()
}

func teamResponse(t domain.Team) TeamResponse {
	return TeamResponse{
		TeamName: t.Name.String(),
		Members:  teamMembers(t.Members.Slice()),
	}
}

func teamMember(m domain.Member) TeamMember {
	return TeamMember{
		UserID:   m.Id.String(),
		Username: m.Name,
		IsActive: m.Status.IsActive(),
	}
}

func teamMembers(ms []domain.Member) []TeamMember {
	res := make([]TeamMember, 0, len(ms))
	for _, m := range ms {
		res = append(res, teamMember(m))
	}
	return res
}
