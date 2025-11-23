package servmembers

import (
	"errors"
	"testing"

	"github.com/eragon-mdi/pr-reviewer-service/internal/common/configs"
	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/eragon-mdi/pr-reviewer-service/internal/service/members/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMembersService_SetMemberIsActive(t *testing.T) {
	tests := []struct {
		name      string
		member    domain.Member
		repoSetup func(*mocks.MembersRepository, domain.Member)
		want      domain.Member
		wantErr   error
	}{
		{
			name: "successful update",
			member: domain.Member{
				Id:     domain.MemberId(uuid.New().String()),
				Status: domain.MemberStatusActive,
			},
			repoSetup: func(mockRepo *mocks.MembersRepository, member domain.Member) {
				mockRepo.EXPECT().UpdateMemberStatus(
					member.Id,
					member.Status,
				).Return(domain.Member{
					Id:     member.Id,
					Name:   "Test User",
					Status: domain.MemberStatusActive,
				}, nil)
			},
			want: domain.Member{
				Name:   "Test User",
				Status: domain.MemberStatusActive,
			},
			wantErr: nil,
		},
		{
			name: "member not found",
			member: domain.Member{
				Id:     domain.MemberId(uuid.New().String()),
				Status: domain.MemberStatusActive,
			},
			repoSetup: func(mockRepo *mocks.MembersRepository, member domain.Member) {
				mockRepo.EXPECT().UpdateMemberStatus(
					member.Id,
					member.Status,
				).Return(domain.Member{}, domain.ErrNotFound)
			},
			want:    domain.Member{},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "internal error",
			member: domain.Member{
				Id:     domain.MemberId(uuid.New().String()),
				Status: domain.MemberStatusInactive,
			},
			repoSetup: func(mockRepo *mocks.MembersRepository, member domain.Member) {
				mockRepo.EXPECT().UpdateMemberStatus(
					member.Id,
					member.Status,
				).Return(domain.Member{}, errors.New("database error"))
			},
			want:    domain.Member{},
			wantErr: domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMembersRepository(t)
			tt.repoSetup(mockRepo, tt.member)

			cfg := &configs.BussinesLogic{
				AllowedReuseToReasign:   true,
				AlloweStatusesToReasign: []string{"active"},
				AllowedRolesToReasign:   []string{"default"},
			}

			service := NewMembersService(cfg, mockRepo)
			got, err := service.SetMemberIsActive(tt.member)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.Status, got.Status)
			}
		})
	}
}

func TestMembersService_MemberReviews(t *testing.T) {
	tests := []struct {
		name      string
		memberId  domain.MemberId
		repoSetup func(*mocks.MembersRepository, domain.MemberId)
		want      domain.Member
		wantErr   error
	}{
		{
			name:     "successful get with reviews",
			memberId: domain.MemberId(uuid.New().String()),
			repoSetup: func(mockRepo *mocks.MembersRepository, memberId domain.MemberId) {
				mockRepo.EXPECT().GetPrReviewsByMember(
					memberId,
				).Return(domain.PullRequests{
					{Id: domain.PrId("pr-1"), Name: domain.PrName("PR1"), AuthorId: domain.MemberId(uuid.New().String()), Status: domain.PrStatusOpen},
				}, nil)
			},
			want: domain.Member{
				Reviews: domain.PullRequests{
					{Id: domain.PrId("pr-1"), Name: domain.PrName("PR1"), AuthorId: domain.MemberId(uuid.New().String()), Status: domain.PrStatusOpen},
				},
			},
			wantErr: nil,
		},
		{
			name:     "empty reviews",
			memberId: domain.MemberId(uuid.New().String()),
			repoSetup: func(mockRepo *mocks.MembersRepository, memberId domain.MemberId) {
				mockRepo.EXPECT().GetPrReviewsByMember(
					memberId,
				).Return(domain.PullRequests{}, nil)
			},
			want:    domain.Member{},
			wantErr: domain.ErrNoContent,
		},
		{
			name:     "member not found",
			memberId: domain.MemberId(uuid.New().String()),
			repoSetup: func(mockRepo *mocks.MembersRepository, memberId domain.MemberId) {
				mockRepo.EXPECT().GetPrReviewsByMember(
					memberId,
				).Return(domain.PullRequests{}, domain.ErrNotFound)
			},
			want:    domain.Member{},
			wantErr: domain.ErrNotFound,
		},
		{
			name:     "internal error",
			memberId: domain.MemberId(uuid.New().String()),
			repoSetup: func(mockRepo *mocks.MembersRepository, memberId domain.MemberId) {
				mockRepo.EXPECT().GetPrReviewsByMember(
					memberId,
				).Return(domain.PullRequests{}, errors.New("database error"))
			},
			want:    domain.Member{},
			wantErr: domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMembersRepository(t)
			tt.repoSetup(mockRepo, tt.memberId)

			cfg := &configs.BussinesLogic{
				AllowedReuseToReasign:   true,
				AlloweStatusesToReasign: []string{"active"},
				AllowedRolesToReasign:   []string{"default"},
			}

			service := NewMembersService(cfg, mockRepo)
			got, err := service.MemberReviews(tt.memberId)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
				if !tt.want.Reviews.Empty() {
					assert.Equal(t, len(tt.want.Reviews), len(got.Reviews))
				}
			}
		})
	}
}
