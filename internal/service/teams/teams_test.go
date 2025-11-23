package servteams_test

import (
	"errors"
	"testing"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	servteams "github.com/eragon-mdi/pr-reviewer-service/internal/service/teams"
	"github.com/eragon-mdi/pr-reviewer-service/internal/service/teams/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTeamsService_NewTeam(t *testing.T) {
	tests := []struct {
		name      string
		team      domain.Team
		repoSetup func(*mocks.TeamsRepository, domain.Team)
		want      domain.Team
		wantErr   error
	}{
		{
			name: "successful create",
			team: func() domain.Team {
				memberId := domain.MemberId(uuid.New().String())
				return domain.NewTeam(
					domain.TeamName("backend"),
					domain.Member{
						Id:     memberId,
						Name:   "User1",
						Status: domain.MemberStatusActive,
					},
				)
			}(),
			repoSetup: func(mockRepo *mocks.TeamsRepository, team domain.Team) {
				mockRepo.EXPECT().CreateTeamWithMembers(
					team.Name,
					team.Members,
				).Return(team, nil)
			},
			want: domain.NewTeam(
				domain.TeamName("backend"),
				domain.Member{
					Name:   "User1",
					Status: domain.MemberStatusActive,
				},
			),
			wantErr: nil,
		},
		{
			name: "duplicate team",
			team: domain.NewTeam(domain.TeamName("backend")),
			repoSetup: func(mockRepo *mocks.TeamsRepository, team domain.Team) {
				mockRepo.EXPECT().CreateTeamWithMembers(
					team.Name,
					team.Members,
				).Return(domain.Team{}, domain.ErrDuplicate)
			},
			want:    domain.Team{},
			wantErr: domain.ErrDuplicate,
		},
		{
			name: "internal error",
			team: domain.NewTeam(domain.TeamName("frontend")),
			repoSetup: func(mockRepo *mocks.TeamsRepository, team domain.Team) {
				mockRepo.EXPECT().CreateTeamWithMembers(
					team.Name,
					team.Members,
				).Return(domain.Team{}, errors.New("database error"))
			},
			want:    domain.Team{},
			wantErr: domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewTeamsRepository(t)
			tt.repoSetup(mockRepo, tt.team)

			service := servteams.NewTeamsService(mockRepo)
			got, err := service.NewTeam(tt.team)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Name, got.Name)
			}
		})
	}
}

func TestTeamsService_TeamWithMembers(t *testing.T) {
	tests := []struct {
		name      string
		teamName  domain.TeamName
		repoSetup func(*mocks.TeamsRepository, domain.TeamName)
		want      domain.Team
		wantErr   error
	}{
		{
			name:     "successful get with members",
			teamName: domain.TeamName("backend"),
			repoSetup: func(mockRepo *mocks.TeamsRepository, teamName domain.TeamName) {
				mockRepo.EXPECT().GetMembersByTeamName(
					teamName,
				).Return(domain.Members{
					{Id: domain.MemberId(uuid.New().String()), Name: "User1", Status: domain.MemberStatusActive},
					{Id: domain.MemberId(uuid.New().String()), Name: "User2", Status: domain.MemberStatusActive},
				}, nil)
			},
			want: domain.NewTeam(
				domain.TeamName("backend"),
				domain.Member{Name: "User1", Status: domain.MemberStatusActive},
				domain.Member{Name: "User2", Status: domain.MemberStatusActive},
			),
			wantErr: nil,
		},
		{
			name:     "empty members",
			teamName: domain.TeamName("empty-team"),
			repoSetup: func(mockRepo *mocks.TeamsRepository, teamName domain.TeamName) {
				mockRepo.EXPECT().GetMembersByTeamName(
					teamName,
				).Return(domain.Members{}, nil)
			},
			want:    domain.NewTeam(domain.TeamName("empty-team")),
			wantErr: domain.ErrNoContent,
		},
		{
			name:     "team not found",
			teamName: domain.TeamName("nonexistent"),
			repoSetup: func(mockRepo *mocks.TeamsRepository, teamName domain.TeamName) {
				mockRepo.EXPECT().GetMembersByTeamName(
					teamName,
				).Return(domain.Members{}, domain.ErrNotFound)
			},
			want:    domain.Team{},
			wantErr: domain.ErrNotFound,
		},
		{
			name:     "internal error",
			teamName: domain.TeamName("error-team"),
			repoSetup: func(mockRepo *mocks.TeamsRepository, teamName domain.TeamName) {
				mockRepo.EXPECT().GetMembersByTeamName(
					teamName,
				).Return(domain.Members{}, errors.New("database error"))
			},
			want:    domain.Team{},
			wantErr: domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewTeamsRepository(t)
			tt.repoSetup(mockRepo, tt.teamName)

			service := servteams.NewTeamsService(mockRepo)
			got, err := service.TeamWithMembers(tt.teamName)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Name, got.Name)
			}
		})
	}
}
