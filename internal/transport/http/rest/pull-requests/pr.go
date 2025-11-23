package restpullrequests

import (
	"context"
	"errors"
	"net/http"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/eragon-mdi/pr-reviewer-service/pkg/validator"
	"github.com/labstack/echo/v4"
)

var (
	ErrBadReqParam = echo.NewHTTPError(http.StatusBadRequest, "bad req param")
	ErrBadReqBody  = echo.NewHTTPError(http.StatusBadRequest, "bad req body")
)

type PullRequestService interface {
	Merge(id domain.PrId) (domain.PullRequest, error)
	NewPullRequest(basePR domain.PullRequestShort) (domain.PullRequest, error)
	Reasign(ctx context.Context, prReasMem domain.PrReasignMember) (domain.PrWithReasignMember, error)
}

func (prt *RestPullRequests) CreatePullRequest(c echo.Context) error {
	var req = &CreatePRRequest{}

	l := prt.l.With("req", req)
	l.Infof("CreatePullRequest called")

	if err := c.Bind(req); err != nil {
		l.Errorf("failed to bind request: %v", err)
		return ErrBadReqBody
	}

	if err := validate(c, req); err != nil {
		l.Errorf("failed validate: %v", err)
		return ErrBadReqBody
	}

	pr, err := prt.s.NewPullRequest(req.domain())
	if err != nil {
		l.Errorf("failed to create pull request: %v", err)

		if errors.Is(err, domain.ErrDuplicate) {
			return domain.HttpErrPRExists()
		}
		if errors.Is(err, domain.ErrNotFound) {
			return domain.HttpErrNotFound()
		}
		return domain.ErrInternal
	}

	l = l.With("pr_id", pr.Id.String())
	l.Infof("pull request created successfully")

	return c.JSON(http.StatusCreated, echo.Map{
		"pr": pullRequestResponse(pr),
	})
}

func (prt *RestPullRequests) MergePullRequest(c echo.Context) error {
	var req = &MergePRRequest{}

	l := prt.l.With("req", req)
	l.Infof("MergePullRequest called")

	if err := c.Bind(req); err != nil {
		l.Errorf("failed to bind request: %v", err)
		return ErrBadReqBody
	}

	if err := validate(c, req); err != nil {
		l.Errorf("failed validate: %v", err)
		return ErrBadReqBody
	}

	pr, err := prt.s.Merge(domain.PrId(req.PullRequestID))
	if err != nil {
		l.Errorf("failed to merge pull request: %v", err)

		if errors.Is(err, domain.ErrNotFound) {
			return domain.HttpErrNotFound()
		}
		if errors.Is(err, domain.ErrConflict) {
			return domain.HttpErrPRMerged()
		}
		return domain.ErrInternal
	}

	l = l.With("pr_id", pr.Id.String())
	l.Infof("pull request merged successfully")

	return c.JSON(http.StatusOK, echo.Map{
		"pr": pullRequestResponse(pr),
	})
}

func (prt *RestPullRequests) ReassignUserForPullRequest(c echo.Context) error {
	var req = &ReassignPRRequest{}

	l := prt.l.With("req", req)
	l.Infof("ReassignUserForPullRequest called")

	if err := c.Bind(req); err != nil {
		l.Errorf("failed to bind request: %v", err)
		return ErrBadReqBody
	}

	if err := validate(c, req); err != nil {
		l.Errorf("failed validate: %v", err)
		return ErrBadReqBody
	}

	prReasMem := domain.PrReasignMember{
		PrId:     domain.PrId(req.PullRequestID),
		MemberId: domain.MemberId(req.OldUserID),
	}

	prWithNewMember, err := prt.s.Reasign(c.Request().Context(), prReasMem)
	if err != nil {
		l.Errorf("failed to reassign pull request: %v", err)

		if errors.Is(err, domain.ErrNotFound) {
			return domain.HttpErrNotFound()
		}
		if errors.Is(err, domain.ErrConflict) {
			return domain.HttpErrPRMerged()
		}
		if errors.Is(err, domain.ErrForbidden) {
			if errors.Is(err, domain.ErrNoContent) {
				return domain.HttpErrNoCandidate()
			}
			return domain.HttpErrNotAssigned()
		}
		return domain.ErrInternal
	}

	l = l.With("pr_id", prWithNewMember.Id.String(),
		"replaced_by", prWithNewMember.MemberId.String())
	l.Infof("pull request reassigned successfully")

	return c.JSON(http.StatusOK, reassignPRResponse(prWithNewMember))
}

func validate(c echo.Context, structure any) error {
	return validator.Validate(c.Request().Context(), structure)
}
