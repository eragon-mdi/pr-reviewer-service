package restmembers

import (
	"context"
	"errors"
	"net/http"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/eragon-mdi/pr-reviewer-service/pkg/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrBadReqParam = echo.NewHTTPError(http.StatusBadRequest, "bad req param")
	ErrBadReqBody  = echo.NewHTTPError(http.StatusBadRequest, "bad req body")
)

type MembersService interface {
	MemberReviews(id domain.MemberId) (domain.Member, error)
	SetMemberIsActive(member domain.Member) (domain.Member, error)
}

func (mt *RestMembers) UserSetIsActive(c echo.Context) error {
	var req = &SetIsActiveRequest{}

	l := mt.l.With("req", req)
	l.Infof("UserSetIsActive called")

	if err := c.Bind(req); err != nil {
		l.Errorf("failed to bind request: %v", err)
		return ErrBadReqBody
	}

	if err := validate(c, req); err != nil {
		l.Errorf("failed validate: %v", err)
		return ErrBadReqBody
	}

	member := req.domain()
	updMember, err := mt.s.SetMemberIsActive(member)
	if err != nil {
		l.Errorf("failed to set member is active: %v", err)

		if errors.Is(err, domain.ErrNotFound) {
			return domain.HttpErrNotFound()
		}
		return domain.ErrInternal
	}

	l = l.With("user_id", updMember.Id.String())
	l.Infof("member status updated successfully")

	return c.JSON(http.StatusOK, echo.Map{
		"user": userResponse(updMember),
	})
}

func (mt *RestMembers) GetUserPeviewsById(c echo.Context) error {
	userID := c.Param("id")

	l := mt.l.With("user_id", userID)
	l.Infof("GetUserPeviewsById called")

	if err := idValidate(userID); err != nil {
		l.Errorf("invalid user_id format: %v", err)
		return ErrBadReqParam
	}

	member, err := mt.s.MemberReviews(domain.MemberId(userID))
	if err != nil && !errors.Is(err, domain.ErrNoContent) {
		l.Errorf("failed to get member reviews: %v", err)

		if errors.Is(err, domain.ErrNotFound) {
			return domain.HttpErrNotFound()
		}
		return domain.ErrInternal
	}

	l = l.With("user_id", member.Id.String())
	l.Infof("member reviews fetched successfully")

	return c.JSON(http.StatusOK, userReviewsResponse(member))
}

func validate(c echo.Context, structure any) error {
	return validator.Validate(ctx(c), structure)
}

func idValidate(id string) error {
	_, err := uuid.Parse(id)
	return err
}

func ctx(c echo.Context) context.Context {
	return c.Request().Context()
}
