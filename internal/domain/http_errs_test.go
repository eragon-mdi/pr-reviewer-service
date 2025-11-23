package domain

import (
	"net/http"
	"testing"
)

func TestCustomHttpError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *CustomHttpError
		want string
	}{
		{
			name: "team exists error",
			err:  HttpErrTeamExists(),
			want: "TEAM_EXISTS: team_name already exists",
		},
		{
			name: "pr exists error",
			err:  HttpErrPRExists(),
			want: "PR_EXISTS: PR id already exists",
		},
		{
			name: "pr merged error",
			err:  HttpErrPRMerged(),
			want: "PR_MERGED: cannot reassign on merged PR",
		},
		{
			name: "not assigned error",
			err:  HttpErrNotAssigned(),
			want: "NOT_ASSIGNED: reviewer is not assigned to this PR",
		},
		{
			name: "no candidate error",
			err:  HttpErrNoCandidate(),
			want: "NO_CANDIDATE: no active replacement candidate in team",
		},
		{
			name: "not found error",
			err:  HttpErrNotFound(),
			want: "NOT_FOUND: resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("CustomHttpError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCustomHttpError(t *testing.T) {
	tests := []struct {
		name     string
		httpCode int
		code     ErrorCode
		msg      string
		wantCode ErrorCode
		wantMsg  string
	}{
		{
			name:     "custom error",
			httpCode: http.StatusBadRequest,
			code:     CodeTeamExists,
			msg:      "test message",
			wantCode: CodeTeamExists,
			wantMsg:  "test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCustomHttpError(tt.httpCode, tt.code, tt.msg)
			if err.HttpCode != tt.httpCode {
				t.Errorf("CustomHttpError.HttpCode = %v, want %v", err.HttpCode, tt.httpCode)
			}
			if err.Code != tt.wantCode {
				t.Errorf("CustomHttpError.Code = %v, want %v", err.Code, tt.wantCode)
			}
			if err.Message != tt.wantMsg {
				t.Errorf("CustomHttpError.Message = %v, want %v", err.Message, tt.wantMsg)
			}
		})
	}
}

func TestHttpErrFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() *CustomHttpError
		wantCode int
		wantErr  ErrorCode
	}{
		{
			name:     "HttpErrTeamExists",
			fn:       HttpErrTeamExists,
			wantCode: http.StatusBadRequest,
			wantErr:  CodeTeamExists,
		},
		{
			name:     "HttpErrPRExists",
			fn:       HttpErrPRExists,
			wantCode: http.StatusConflict,
			wantErr:  CodePRExists,
		},
		{
			name:     "HttpErrPRMerged",
			fn:       HttpErrPRMerged,
			wantCode: http.StatusConflict,
			wantErr:  CodePRMerged,
		},
		{
			name:     "HttpErrNotAssigned",
			fn:       HttpErrNotAssigned,
			wantCode: http.StatusConflict,
			wantErr:  CodeNotAssigned,
		},
		{
			name:     "HttpErrNoCandidate",
			fn:       HttpErrNoCandidate,
			wantCode: http.StatusConflict,
			wantErr:  CodeNoCandidate,
		},
		{
			name:     "HttpErrNotFound",
			fn:       HttpErrNotFound,
			wantCode: http.StatusNotFound,
			wantErr:  CodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err.HttpCode != tt.wantCode {
				t.Errorf("%s().HttpCode = %v, want %v", tt.name, err.HttpCode, tt.wantCode)
			}
			if err.Code != tt.wantErr {
				t.Errorf("%s().Code = %v, want %v", tt.name, err.Code, tt.wantErr)
			}
		})
	}
}
