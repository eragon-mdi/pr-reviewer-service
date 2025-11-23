package sqlrepo

import (
	"context"
	"database/sql"
	"time"

	sqlstore "github.com/eragon-mdi/go-playground/storage/sql"
	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/eragon-mdi/pr-reviewer-service/internal/repository/sql/queries"
	servpullrequests "github.com/eragon-mdi/pr-reviewer-service/internal/service/pull-requests"
	"github.com/go-faster/errors"
	"github.com/lib/pq"
)

type pullRequestsRepo struct {
	s sqlstore.Storage
}

func NewPullRequestsRepo(s sqlstore.Storage) *pullRequestsRepo {
	return &pullRequestsRepo{s: s}
}

func (r *pullRequestsRepo) CreatePullRequest(pr domain.PullRequest) (domain.PullRequest, error) {
	ctx := context.Background()

	var existingID int
	err := r.s.QueryRow("SELECT id FROM pull_requests WHERE uuid = $1", pr.Id.String()).Scan(&existingID)
	if err == nil {
		return domain.PullRequest{}, domain.ErrDuplicate
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return domain.PullRequest{}, errors.Wrap(err, ErrFailedQuery)
	}

	var authorID int
	err = r.s.QueryRow("SELECT id FROM members WHERE uuid = $1", pr.AuthorId.String()).Scan(&authorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PullRequest{}, domain.ErrNotFound
		}
		return domain.PullRequest{}, errors.Wrap(err, ErrFailedQuery)
	}

	_, err = r.s.ExecContext(ctx, queries.CreatePullRequest, pr.Id.String(), pr.Name.String(), authorID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return domain.PullRequest{}, domain.ErrDuplicate
		}
		return domain.PullRequest{}, errors.Wrap(err, ErrFailedExec)
	}

	createdPr, err := r.GetPullRequestByUUID(ctx, pr.Id)
	if err != nil {
		var existingID int
		checkErr := r.s.QueryRow("SELECT id FROM pull_requests WHERE uuid = $1", pr.Id.String()).Scan(&existingID)
		if checkErr == nil {
			return domain.PullRequest{}, domain.ErrDuplicate
		}
		return domain.PullRequest{}, err
	}

	return createdPr, nil
}

func (r *pullRequestsRepo) GetPullRequestByUUID(ctx context.Context, prId domain.PrId) (domain.PullRequest, error) {
	var id int
	var uuid string
	var title string
	var authorUUID string
	var status string
	var createdAt time.Time
	var mergedAt sql.NullTime
	var version int

	err := r.s.QueryRow(queries.GetPullRequestByUUID, prId.String()).Scan(
		&id, &uuid, &title, &authorUUID, &status, &createdAt, &mergedAt, &version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PullRequest{}, domain.ErrNotFound
		}
		return domain.PullRequest{}, errors.Wrap(err, ErrFailedQuery)
	}

	var prStatus domain.PrStatus
	if status == "MERGED" {
		prStatus = domain.PrStatusMerged
	} else {
		prStatus = domain.PrStatusOpen
	}

	pr := domain.PullRequest{
		Id:        domain.PrId(uuid),
		Name:      domain.PrName(title),
		AuthorId:  domain.MemberId(authorUUID),
		Status:    prStatus,
		CreatedAt: createdAt,
	}

	if mergedAt.Valid {
		pr.MergedAt = mergedAt.Time
	}

	reviewers, err := r.GetPullRequestReviewers(ctx, prId)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return domain.PullRequest{}, err
	}
	pr.AssignedReviews = reviewers

	return pr, nil
}

func (r *pullRequestsRepo) GetPullRequestReviewers(ctx context.Context, prId domain.PrId) (domain.Members, error) {
	rows, err := r.s.Query(queries.GetPullRequestReviewers, prId.String())
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

func (r *pullRequestsRepo) MergePullRequest(prId domain.PrId) (domain.PullRequest, error) {
	ctx := context.Background()

	pr, err := r.GetPullRequestByUUID(ctx, prId)
	if err != nil {
		return domain.PullRequest{}, err
	}

	if pr.Status == domain.PrStatusMerged {
		return pr, nil
	}

	var id int
	var uuid string
	var title string
	var authorID int
	var statusID int
	var createdAt time.Time
	var mergedAt sql.NullTime
	var version int

	err = r.s.QueryRow(queries.MergePullRequest, prId.String()).Scan(
		&id, &uuid, &title, &authorID, &statusID, &createdAt, &mergedAt, &version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r.GetPullRequestByUUID(ctx, prId)
		}
		return domain.PullRequest{}, errors.Wrap(err, ErrFailedQuery)
	}

	return r.GetPullRequestByUUID(ctx, prId)
}

func (r *pullRequestsRepo) BeginReasignTx(ctx context.Context) (servpullrequests.ReassignTx, error) {
	tx, err := r.s.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrFailedStartTX)
	}
	return &reassignTx{tx: tx, s: r.s}, nil
}

type reassignTx struct {
	tx interface{}
	s  sqlstore.Storage
}

func (rtx *reassignTx) GetPullRequestMembersHistories(ctx context.Context, prId domain.PrId) (domain.MembersHistories, error) {
	rows, err := rtx.s.Query(queries.GetPullRequestMembersHistories, prId.String(), "")
	if err != nil {
		return nil, errors.Wrap(err, ErrFailedQuery)
	}
	defer rows.Close()

	histories := make([]domain.MemberHistory, 0)
	for rows.Next() {
		var id int
		var uuid string
		var name string
		var isActive bool
		var role string
		var assignedAt sql.NullTime
		var wasAssignedBefore bool

		if err := rows.Scan(&id, &uuid, &name, &isActive, &role, &assignedAt, &wasAssignedBefore); err != nil {
			return nil, errors.Wrap(err, ErrFailedScan)
		}

		var memberRole domain.MemberRole
		if role == "author" {
			memberRole = domain.MemberRolePrAuthor
		} else if role == "reassigned" {
			memberRole = domain.MemberRoleHadReasigned
		} else {
			memberRole = domain.MemberRoleDefault
		}

		status := domain.MemberStatusIsActiveByBool(isActive)
		history := domain.NewMemberHistory(
			domain.MemberId(uuid),
			status,
			memberRole,
			wasAssignedBefore,
		)
		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, ErrRowsIterations)
	}

	if len(histories) == 0 {
		return nil, domain.ErrNoContent
	}

	return domain.MembersHistories(histories), nil
}

func (rtx *reassignTx) AssignMember(ctx context.Context, prId domain.PrId, newMemberId domain.MemberId) (domain.PullRequest, error) {
	var id int
	var uuid string
	var title string
	var authorUUID string
	var status string
	var createdAt time.Time
	var mergedAt sql.NullTime
	var version int

	err := rtx.s.QueryRow(queries.AssignMemberToPR, prId.String(), "", newMemberId.String()).Scan(
		&id, &uuid, &title, &authorUUID, &status, &createdAt, &mergedAt, &version,
	)
	if err != nil {
		return domain.PullRequest{}, errors.Wrap(err, ErrFailedQuery)
	}

	var prStatus domain.PrStatus
	if status == "MERGED" {
		prStatus = domain.PrStatusMerged
	} else {
		prStatus = domain.PrStatusOpen
	}

	pr := domain.PullRequest{
		Id:        domain.PrId(uuid),
		Name:      domain.PrName(title),
		AuthorId:  domain.MemberId(authorUUID),
		Status:    prStatus,
		CreatedAt: createdAt,
	}

	if mergedAt.Valid {
		pr.MergedAt = mergedAt.Time
	}

	reviewers, err := rtx.getPullRequestReviewers(ctx, prId)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return domain.PullRequest{}, err
	}
	pr.AssignedReviews = reviewers

	return pr, nil
}

func (rtx *reassignTx) getPullRequestReviewers(ctx context.Context, prId domain.PrId) (domain.Members, error) {
	rows, err := rtx.s.Query(queries.GetPullRequestReviewers, prId.String())
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

func (rtx *reassignTx) Commit() error {
	return nil
}

func (rtx *reassignTx) Rollback() error {
	return nil
}
