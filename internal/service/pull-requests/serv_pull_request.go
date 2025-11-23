package servpullrequests

type PrService struct {
	repo    Repository
	memServ MemberService
}

func NewPullRequestService(r Repository, ms MemberService) *PrService {
	return &PrService{
		repo:    r,
		memServ: ms,
	}
}

type Repository interface {
	PullRequestsRepository
}
