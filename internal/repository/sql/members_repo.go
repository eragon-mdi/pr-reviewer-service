package sqlrepo

import (
	"context"
	"database/sql"

	sqlstore "github.com/eragon-mdi/go-playground/storage/sql"
	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/eragon-mdi/pr-reviewer-service/internal/repository/sql/queries"
	"github.com/go-faster/errors"
)

type membersRepo struct {
	s sqlstore.Storage
}

func NewMembersRepo(s sqlstore.Storage) *membersRepo {
	return &membersRepo{s: s}
}

func (r *membersRepo) UpdateMemberStatus(memberId domain.MemberId, status domain.MemberStatus) (domain.Member, error) {
	ctx := context.Background()

	var id int
	var uuid string
	var name string
	var isActive bool

	err := r.s.QueryRow(queries.UpdateMemberStatus, memberId.String(), status.IsActive()).Scan(&id, &uuid, &name, &isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Member{}, domain.ErrNotFound
		}
		return domain.Member{}, errors.Wrap(err, ErrFailedQuery)
	}

	teamName, err := r.getTeamNameByMemberId(ctx, memberId)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return domain.Member{}, err
	}

	member := domain.MemberBuilder(domain.MemberId(uuid)).
		Name(name).
		Status(domain.MemberStatusIsActiveByBool(isActive)).
		Build()
	member.Team = teamName

	return member, nil
}

func (r *membersRepo) GetPrReviewsByMember(memberId domain.MemberId) (domain.PullRequests, error) {
	rows, err := r.s.Query(queries.GetPrReviewsByMember, memberId.String())
	if err != nil {
		return nil, errors.Wrap(err, ErrFailedQuery)
	}
	defer rows.Close()

	prs := make([]domain.PullRequestShort, 0)
	for rows.Next() {
		var prID string
		var prName string
		var authorID string
		var status string

		if err := rows.Scan(&prID, &prName, &authorID, &status); err != nil {
			return nil, errors.Wrap(err, ErrFailedScan)
		}

		var prStatus domain.PrStatus
		if status == "MERGED" {
			prStatus = domain.PrStatusMerged
		} else {
			prStatus = domain.PrStatusOpen
		}

		pr := domain.PullRequestShort{
			Id:       domain.PrId(prID),
			Name:     domain.PrName(prName),
			AuthorId: domain.MemberId(authorID),
			Status:   prStatus,
		}
		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, ErrRowsIterations)
	}

	return domain.PullRequests(prs), nil
}

func (r *membersRepo) getTeamNameByMemberId(ctx context.Context, memberId domain.MemberId) (domain.TeamName, error) {
	var teamName string
	err := r.s.QueryRowContext(ctx, queries.GetTeamNameByMemberId, memberId.String()).Scan(&teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TeamName(""), domain.ErrNotFound
		}
		return domain.TeamName(""), errors.Wrap(err, ErrFailedQuery)
	}
	return domain.TeamName(teamName), nil
}
