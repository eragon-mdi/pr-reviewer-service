package servmembers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eragon-mdi/pr-reviewer-service/internal/common/configs"
	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	servmembers "github.com/eragon-mdi/pr-reviewer-service/internal/service/members"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMembersService_ReasignMember(t *testing.T) {
	tests := []struct {
		name        string
		memId       domain.MemberId
		mems        domain.MembersHistories
		cfg         *configs.BussinesLogic
		wantErr     bool
		checkResult func(t *testing.T, result domain.MemberId)
	}{
		{
			name:  "successful reassign - allowed member found",
			memId: domain.MemberId(uuid.New().String()),
			mems: domain.MembersHistories{
				domain.NewMemberHistory(
					domain.MemberId(uuid.New().String()),
					domain.MemberStatusActive,
					domain.MemberRoleDefault,
					false,
				),
			},
			cfg: &configs.BussinesLogic{
				AllowedReuseToReasign:   true,
				AlloweStatusesToReasign: []string{"active"},
				AllowedRolesToReasign:   []string{"default"},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result domain.MemberId) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:  "no allowed members - forbidden",
			memId: domain.MemberId(uuid.New().String()),
			mems: domain.MembersHistories{
				domain.NewMemberHistory(
					domain.MemberId(uuid.New().String()),
					domain.MemberStatusInactive,
					domain.MemberRoleDefault,
					false,
				),
			},
			cfg: &configs.BussinesLogic{
				AllowedReuseToReasign:   true,
				AlloweStatusesToReasign: []string{"active"},
				AllowedRolesToReasign:   []string{"default"},
			},
			wantErr: true,
			checkResult: func(t *testing.T, result domain.MemberId) {
				assert.Empty(t, result)
			},
		},
		{
			name:  "was assigned before, history reuse false",
			memId: domain.MemberId(uuid.New().String()),
			mems: domain.MembersHistories{
				domain.NewMemberHistory(
					domain.MemberId(uuid.New().String()),
					domain.MemberStatusActive,
					domain.MemberRoleDefault,
					true,
				),
			},
			cfg: &configs.BussinesLogic{
				AllowedReuseToReasign:   false,
				AlloweStatusesToReasign: []string{"active"},
				AllowedRolesToReasign:   []string{"default"},
			},
			wantErr: true,
			checkResult: func(t *testing.T, result domain.MemberId) {
				assert.Empty(t, result)
			},
		},
		{
			name:  "multiple members, one allowed",
			memId: domain.MemberId(uuid.New().String()),
			mems: domain.MembersHistories{
				domain.NewMemberHistory(
					domain.MemberId(uuid.New().String()),
					domain.MemberStatusInactive,
					domain.MemberRoleDefault,
					false,
				),
				domain.NewMemberHistory(
					domain.MemberId(uuid.New().String()),
					domain.MemberStatusActive,
					domain.MemberRoleDefault,
					false,
				),
			},
			cfg: &configs.BussinesLogic{
				AllowedReuseToReasign:   true,
				AlloweStatusesToReasign: []string{"active"},
				AllowedRolesToReasign:   []string{"default"},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result domain.MemberId) {
				assert.NotEmpty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := servmembers.NewMembersService(tt.cfg, nil)

			result, err := service.ReasignMember(context.Background(), tt.memId, tt.mems)

			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrForbidden))
			} else {
				assert.NoError(t, err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}
