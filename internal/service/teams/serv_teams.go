package servteams

type TeamsService struct {
	repo Repository
}

func NewTeamsService(r Repository) *TeamsService {
	return &TeamsService{
		repo: r,
	}
}

type Repository interface {
	TeamsRepository
}
