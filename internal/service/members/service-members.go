package servmembers

import (
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/configs"
	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
)

type MembersService struct {
	repo         Repository
	allowedRoles domain.AllowedRules
}

func NewMembersService(cfg *configs.BussinesLogic, r Repository) *MembersService {
	return &MembersService{
		repo: r,
		allowedRoles: domain.NewAllowedRules(
			cfg.AllowedReuseToReasign,
			domain.MembersStatusesFromSliceOfStrings(cfg.AlloweStatusesToReasign),
			domain.MembersRolesFromSliceOfStrings(cfg.AllowedRolesToReasign),
		),
	}
}

type Repository interface {
	MembersRepository
}
