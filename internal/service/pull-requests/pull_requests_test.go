package servpullrequests_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	servpullrequests "github.com/eragon-mdi/pr-reviewer-service/internal/service/pull-requests"
	"github.com/eragon-mdi/pr-reviewer-service/internal/service/pull-requests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPrService_Merge(t *testing.T) {
	tests := []struct {
		name      string
		prId      domain.PrId
		repoSetup func(*mocks.PullRequestsRepository)
		want      domain.PullRequest
		wantErr   error
	}{
		{
			name: "successful merge",
			prId: domain.PrId("pr-123"),
			repoSetup: func(mockRepo *mocks.PullRequestsRepository) {
				mockRepo.EXPECT().MergePullRequest(
					domain.PrId("pr-123"),
				).Return(domain.PullRequest{
					Id:        domain.PrId("pr-123"),
					Name:      domain.PrName("Test PR"),
					AuthorId:  domain.MemberId(uuid.New().String()),
					Status:    domain.PrStatusMerged,
					CreatedAt: time.Now(),
					MergedAt:  time.Now(),
				}, nil)
			},
			want: domain.PullRequest{
				Id:     domain.PrId("pr-123"),
				Status: domain.PrStatusMerged,
			},
			wantErr: nil,
		},
		{
			name: "pr not found",
			prId: domain.PrId("pr-999"),
			repoSetup: func(mockRepo *mocks.PullRequestsRepository) {
				mockRepo.EXPECT().MergePullRequest(
					domain.PrId("pr-999"),
				).Return(domain.PullRequest{}, domain.ErrNotFound)
			},
			want:    domain.PullRequest{},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "conflict error",
			prId: domain.PrId("pr-456"),
			repoSetup: func(mockRepo *mocks.PullRequestsRepository) {
				mockRepo.EXPECT().MergePullRequest(
					domain.PrId("pr-456"),
				).Return(domain.PullRequest{}, domain.ErrConflict)
			},
			want:    domain.PullRequest{},
			wantErr: domain.ErrConflict,
		},
		{
			name: "internal error",
			prId: domain.PrId("pr-789"),
			repoSetup: func(mockRepo *mocks.PullRequestsRepository) {
				mockRepo.EXPECT().MergePullRequest(
					domain.PrId("pr-789"),
				).Return(domain.PullRequest{}, errors.New("database error"))
			},
			want:    domain.PullRequest{},
			wantErr: domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPullRequestsRepository(t)
			mockMemberService := mocks.NewMemberService(t)
			tt.repoSetup(mockRepo)

			service := servpullrequests.NewPullRequestService(mockRepo, mockMemberService)
			got, err := service.Merge(tt.prId)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Status, got.Status)
			}
		})
	}
}

func TestPrService_Reasign(t *testing.T) {
	tests := []struct {
		name        string
		prReasMem   domain.PrReasignMember
		repoSetup   func(*mocks.PullRequestsRepository, *mocks.ReassignTx, domain.PrReasignMember)
		memberSetup func(*mocks.MemberService, domain.PrReasignMember)
		want        domain.PrWithReasignMember
		wantErr     error
	}{
		{
			name: "successful reassign",
			prReasMem: func() domain.PrReasignMember {
				oldMemberID := uuid.New().String()
				return domain.PrReasignMember{
					PrId:     domain.PrId("pr-123"),
					MemberId: domain.MemberId(oldMemberID),
				}
			}(),
			repoSetup: func(mockRepo *mocks.PullRequestsRepository, mockTx *mocks.ReassignTx, prReasMem domain.PrReasignMember) {
				candidateID := "candidate-123"
				newMemberID := "new-member-123"

				mockRepo.EXPECT().BeginReasignTx(context.Background()).Return(mockTx, nil)
				mockTx.EXPECT().GetPullRequestMembersHistories(
					context.Background(),
					prReasMem.PrId,
				).Return(domain.MembersHistories{
					domain.NewMemberHistory(
						domain.MemberId(candidateID),
						domain.MemberStatusActive,
						domain.MemberRoleDefault,
						false,
					),
				}, nil)
				mockTx.EXPECT().AssignMember(
					context.Background(),
					prReasMem.PrId,
					domain.MemberId(newMemberID),
				).Return(domain.PullRequest{
					Id:     prReasMem.PrId,
					Status: domain.PrStatusOpen,
				}, nil)
				mockTx.EXPECT().Commit().Return(nil)
			},
			memberSetup: func(mockMemberService *mocks.MemberService, prReasMem domain.PrReasignMember) {
				candidateID := "candidate-123"
				newMemberID := "new-member-123"

				mockMemberService.EXPECT().ReasignMember(
					context.Background(),
					prReasMem.MemberId,
					domain.MembersHistories{
						domain.NewMemberHistory(
							domain.MemberId(candidateID),
							domain.MemberStatusActive,
							domain.MemberRoleDefault,
							false,
						),
					},
				).Return(domain.MemberId(newMemberID), nil)
			},
			want: domain.PrWithReasignMember{
				PullRequest: domain.PullRequest{
					Id:     domain.PrId("pr-123"),
					Status: domain.PrStatusOpen,
				},
			},
			wantErr: nil,
		},
		{
			name: "pr not found",
			prReasMem: domain.PrReasignMember{
				PrId:     domain.PrId("pr-999"),
				MemberId: domain.MemberId(uuid.New().String()),
			},
			repoSetup: func(mockRepo *mocks.PullRequestsRepository, mockTx *mocks.ReassignTx, prReasMem domain.PrReasignMember) {
				mockRepo.EXPECT().BeginReasignTx(context.Background()).Return(mockTx, nil)
				mockTx.EXPECT().GetPullRequestMembersHistories(
					context.Background(),
					prReasMem.PrId,
				).Return(domain.MembersHistories{}, domain.ErrNotFound)
				mockTx.EXPECT().Rollback().Return(nil)
			},
			memberSetup: func(mockMemberService *mocks.MemberService, prReasMem domain.PrReasignMember) {},
			want:        domain.PrWithReasignMember{},
			wantErr:     domain.ErrNotFound,
		},
		{
			name: "no content",
			prReasMem: domain.PrReasignMember{
				PrId:     domain.PrId("pr-123"),
				MemberId: domain.MemberId(uuid.New().String()),
			},
			repoSetup: func(mockRepo *mocks.PullRequestsRepository, mockTx *mocks.ReassignTx, prReasMem domain.PrReasignMember) {
				mockRepo.EXPECT().BeginReasignTx(context.Background()).Return(mockTx, nil)
				mockTx.EXPECT().GetPullRequestMembersHistories(
					context.Background(),
					prReasMem.PrId,
				).Return(domain.MembersHistories{}, domain.ErrNoContent)
				mockTx.EXPECT().Rollback().Return(nil)
			},
			memberSetup: func(mockMemberService *mocks.MemberService, prReasMem domain.PrReasignMember) {},
			want:        domain.PrWithReasignMember{},
			wantErr:     domain.ErrNoContent,
		},
		{
			name: "member service error",
			prReasMem: func() domain.PrReasignMember {
				oldMemberID := uuid.New().String()
				return domain.PrReasignMember{
					PrId:     domain.PrId("pr-123"),
					MemberId: domain.MemberId(oldMemberID),
				}
			}(),
			repoSetup: func(mockRepo *mocks.PullRequestsRepository, mockTx *mocks.ReassignTx, prReasMem domain.PrReasignMember) {
				candidateID := "candidate-456"

				mockRepo.EXPECT().BeginReasignTx(context.Background()).Return(mockTx, nil)
				mockTx.EXPECT().GetPullRequestMembersHistories(
					context.Background(),
					prReasMem.PrId,
				).Return(domain.MembersHistories{
					domain.NewMemberHistory(
						domain.MemberId(candidateID),
						domain.MemberStatusActive,
						domain.MemberRoleDefault,
						false,
					),
				}, nil)
				mockTx.EXPECT().Rollback().Return(nil)
			},
			memberSetup: func(mockMemberService *mocks.MemberService, prReasMem domain.PrReasignMember) {
				candidateID := "candidate-456"

				mockMemberService.EXPECT().ReasignMember(
					context.Background(),
					prReasMem.MemberId,
					domain.MembersHistories{
						domain.NewMemberHistory(
							domain.MemberId(candidateID),
							domain.MemberStatusActive,
							domain.MemberRoleDefault,
							false,
						),
					},
				).Return(domain.MemberId(""), domain.ErrForbidden)
			},
			want:    domain.PrWithReasignMember{},
			wantErr: domain.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPullRequestsRepository(t)
			mockMemberService := mocks.NewMemberService(t)
			mockTx := mocks.NewReassignTx(t)

			tt.repoSetup(mockRepo, mockTx, tt.prReasMem)
			tt.memberSetup(mockMemberService, tt.prReasMem)

			service := servpullrequests.NewPullRequestService(mockRepo, mockMemberService)
			got, err := service.Reasign(context.Background(), tt.prReasMem)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.PullRequest.Id, got.PullRequest.Id)
			}
		})
	}
}
