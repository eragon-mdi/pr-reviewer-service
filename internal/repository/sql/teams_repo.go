package sqlrepo

import (
	"context"
	"database/sql"

	sqlstore "github.com/eragon-mdi/go-playground/storage/sql"
	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/eragon-mdi/pr-reviewer-service/internal/repository/sql/queries"
	"github.com/go-faster/errors"
	"github.com/lib/pq"
)

const (
	ErrFailedQuery        = "repo: failed query"
	ErrFailedExec         = "repo: failed exec"
	ErrFailedScan         = "repo: failed to scan row"
	ErrFailedAffectedRows = "repo: failed to get number of affected rows"
	ErrFailedStartTX      = "repo: failed to start tx"
	ErrFailedCommitTX     = "repo: failed to commit tx"
	ErrFailedRollbackTX   = "repo: failed rollback tx"
	ErrRowsIterations     = "repo: rows iteration error"
)

type teamsRepo struct {
	s sqlstore.Storage
}

func NewTeamsRepo(s sqlstore.Storage) *teamsRepo {
	return &teamsRepo{s: s}
}

func (r *teamsRepo) CreateTeamWithMembers(teamName domain.TeamName, members domain.Members) (domain.Team, error) {
	ctx := context.Background()

	tx, err := r.s.BeginTx(ctx, nil)
	if err != nil {
		return domain.Team{}, errors.Wrap(err, ErrFailedStartTX)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var teamID int
	err = tx.QueryRow("SELECT id FROM teams WHERE name = $1", teamName.String()).Scan(&teamID)
	if err == nil {
		_ = tx.Rollback()
		return domain.Team{}, domain.ErrDuplicate
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return domain.Team{}, errors.Wrap(err, ErrFailedQuery)
	}

	err = tx.QueryRow("INSERT INTO teams (name) VALUES ($1) RETURNING id", teamName.String()).Scan(&teamID)
	if err != nil {
		return domain.Team{}, errors.Wrap(err, ErrFailedExec)
	}

	if len(members) > 0 {
		uuids := make([]string, len(members))
		names := make([]string, len(members))
		isActives := make([]bool, len(members))

		for i, m := range members {
			uuids[i] = m.Id.String()
			names[i] = m.Name
			isActives[i] = m.Status.IsActive()
		}

		rows, err := tx.Query(queries.CreateTeamWithMembers, teamName.String(), pq.Array(uuids), pq.Array(names), pq.Array(isActives))
		if err != nil {
			return domain.Team{}, errors.Wrap(err, ErrFailedQuery)
		}
		defer rows.Close()

		memberUUIDs := make([]string, 0, len(members))
		for rows.Next() {
			var id int
			var uuid string
			var name string
			var isActive bool
			if err := rows.Scan(&id, &uuid, &name, &isActive); err != nil {
				return domain.Team{}, errors.Wrap(err, ErrFailedScan)
			}
			memberUUIDs = append(memberUUIDs, uuid)
		}

		_, err = tx.Exec(queries.LinkMembersToTeam, teamID, pq.Array(memberUUIDs))
		if err != nil {
			return domain.Team{}, errors.Wrap(err, ErrFailedExec)
		}
	}

	if err := tx.Commit(); err != nil {
		return domain.Team{}, errors.Wrap(err, ErrFailedCommitTX)
	}

	return r.GetTeamWithMembers(context.Background(), teamName)
}

func (r *teamsRepo) GetTeamWithMembers(ctx context.Context, teamName domain.TeamName) (domain.Team, error) {
	rows, err := r.s.Query(queries.GetMembersByTeamName, teamName.String())
	if err != nil {
		return domain.Team{}, errors.Wrap(err, ErrFailedQuery)
	}
	defer rows.Close()

	members := make([]domain.Member, 0)
	for rows.Next() {
		var id int
		var uuid string
		var name string
		var isActive bool

		if err := rows.Scan(&id, &uuid, &name, &isActive); err != nil {
			return domain.Team{}, errors.Wrap(err, ErrFailedScan)
		}

		status := domain.MemberStatusIsActiveByBool(isActive)
		member := domain.MemberBuilder(domain.MemberId(uuid)).
			Name(name).
			Status(status).
			Build()
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return domain.Team{}, errors.Wrap(err, ErrRowsIterations)
	}

	if len(members) == 0 {
		return domain.Team{}, domain.ErrNotFound
	}

	return domain.NewTeam(teamName, members...), nil
}

func (r *teamsRepo) GetMembersByTeamName(teamName domain.TeamName) (domain.Members, error) {
	rows, err := r.s.Query(queries.GetMembersByTeamName, teamName.String())
	if err != nil {
		return nil, errors.Wrap(err, ErrFailedQuery)
	}
	defer rows.Close()

	members := make([]domain.Member, 0)
	for rows.Next() {
		var id int
		var uuid string
		var name string
		var isActive bool

		if err := rows.Scan(&id, &uuid, &name, &isActive); err != nil {
			return nil, errors.Wrap(err, ErrFailedScan)
		}

		status := domain.MemberStatusIsActiveByBool(isActive)
		member := domain.MemberBuilder(domain.MemberId(uuid)).
			Name(name).
			Status(status).
			Build()
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, ErrRowsIterations)
	}

	return domain.Members(members), nil
}

func (r *teamsRepo) GetTeamNameByMemberId(ctx context.Context, memberId domain.MemberId) (domain.TeamName, error) {
	var teamName string
	err := r.s.QueryRow(queries.GetTeamNameByMemberId, memberId.String()).Scan(&teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TeamName(""), domain.ErrNotFound
		}
		return domain.TeamName(""), errors.Wrap(err, ErrFailedQuery)
	}
	return domain.TeamName(teamName), nil
}
