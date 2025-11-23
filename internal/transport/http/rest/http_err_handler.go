package resttransport

import (
	"errors"
	"net/http"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func HTTPErrorHandler(err error, c echo.Context) {
	var customHttpError *domain.CustomHttpError

	if errors.As(err, &customHttpError) {
		_ = c.JSON(customHttpError.HttpCode, ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    string(customHttpError.Code),
				Message: customHttpError.Message,
			},
		})
		return
	}

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		msg := httpErr.Message
		if s, ok := msg.(string); ok {
			msg = s
		}
		_ = c.JSON(httpErr.Code, ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: msg.(string),
			},
		})
		return
	}

	_ = c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		},
	})
}
