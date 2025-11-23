package domain

import "errors"

var (
	ErrNotFound   = errors.New("no content found")
	ErrNoContent  = errors.New("no content")
	ErrValidation = errors.New("bad expression")
	ErrDuplicate  = errors.New("duplicate")
	ErrInternal   = errors.New("internal service err. try again later")
	ErrConflict   = errors.New("business conflict")
	ErrForbidden  = errors.New("err forbidden")
)
