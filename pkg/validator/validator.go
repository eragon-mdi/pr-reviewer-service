package validator

import (
	"context"
	"errors"
	"regexp"

	"github.com/go-playground/validator"
)

func init() {
	SetValidator(New())
}

var global *validator.Validate

const (
	ErrInvalidFormat      = "Invalid format"
	ErrFieldRequired      = "Field is required"
	ErrFieldExceedsMaxLen = "Field exceeds maximum length"
	ErrFieldBelowMinLen   = "Field is below minimum length"
	ErrFieldExceedsMaxVal = "Field exceeds maximum value"
	ErrFieldBelowMinVal   = "Field is below minimum value"
	ErrUnknownValidation  = "Unknown validation error"
)

var mapValidationErrors = map[string]string{
	"tag":      ErrInvalidFormat,
	"required": ErrFieldRequired,
	"max":      ErrFieldExceedsMaxLen,
	"min":      ErrFieldBelowMinLen,
	"lt":       ErrFieldExceedsMaxVal,
	"lte":      ErrFieldExceedsMaxVal,
	"gt":       ErrFieldBelowMinVal,
	"gte":      ErrFieldBelowMinVal,
	"uuid":     "Invalid UUID format",
	"oneof":    "Invalid value, must be one of allowed options",
}

func New() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("tag", validateTag)

	return v
}

func SetValidator(v *validator.Validate) {
	global = v
}

func Validator() *validator.Validate {
	return global
}

func validateTag(fl validator.FieldLevel) bool {
	re, _ := regexp.Compile(`^#[a-z0-9_\-]+$`)
	return re.MatchString(fl.Field().String())
}

func Validate(ctx context.Context, structure any) error {
	return parseValidationErrors(Validator().StructCtx(ctx, structure))
}

func parseValidationErrors(err error) error {
	if err == nil {
		return nil
	}

	vErrors, ok := err.(validator.ValidationErrors)
	if !ok || len(vErrors) == 0 {
		return nil
	}

	validationError := vErrors[0]
	vErrDescription, ok := mapValidationErrors[validationError.Tag()]
	if !ok {
		vErrDescription = ErrUnknownValidation
	}

	return errors.New(vErrDescription + ": " + validationError.Namespace())
}
