package domain

import "net/http"

type ErrorCode string

const (
	CodeTeamExists  ErrorCode = "TEAM_EXISTS"
	CodePRExists    ErrorCode = "PR_EXISTS"
	CodePRMerged    ErrorCode = "PR_MERGED"
	CodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	CodeNoCandidate ErrorCode = "NO_CANDIDATE"
	CodeNotFound    ErrorCode = "NOT_FOUND"
)

type CustomHttpError struct {
	HttpCode int
	Code     ErrorCode `json:"code"`
	Message  string    `json:"message"`
}

func (e *CustomHttpError) Error() string {
	return string(e.Code) + ": " + e.Message
}

func NewCustomHttpError(httpCode int, code ErrorCode, msg string) *CustomHttpError {
	return &CustomHttpError{
		HttpCode: httpCode,
		Code:     code,
		Message:  msg,
	}
}

func HttpErrTeamExists() *CustomHttpError {
	return NewCustomHttpError(http.StatusBadRequest, CodeTeamExists, "team_name already exists")
}

func HttpErrPRExists() *CustomHttpError {
	return NewCustomHttpError(http.StatusConflict, CodePRExists, "PR id already exists")
}

func HttpErrPRMerged() *CustomHttpError {
	return NewCustomHttpError(http.StatusConflict, CodePRMerged, "cannot reassign on merged PR")
}

func HttpErrNotAssigned() *CustomHttpError {
	return NewCustomHttpError(http.StatusConflict, CodeNotAssigned, "reviewer is not assigned to this PR")
}

func HttpErrNoCandidate() *CustomHttpError {
	return NewCustomHttpError(http.StatusConflict, CodeNoCandidate, "no active replacement candidate in team")
}

func HttpErrNotFound() *CustomHttpError {
	return NewCustomHttpError(http.StatusNotFound, CodeNotFound, "resource not found")
}
