package health

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

var (
	// ErrInvalidArgument is returned when one or more arguments are invalid.
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrAlreadyExists    = errors.New("already exists")
	ErrBadRequest       = errors.New("bad request")
	ErrNotFound         = errors.New("not found")
	errBadRoute         = errors.New("bad route")
	ErrInvalidRequest   = errors.New("invalid params in request")
	ErrUserUnAuthorized = errors.New("user is not authorized")
	InternalError       = errors.New("internal error")
	ErrStatusForbidden  = errors.New("status is forbidden")
)

type ContextHTTPKey struct{}

type HTTPInfo struct {
	Method   string
	URL      string
	From     string
	Protocol string
}

type errorCode interface {
	Code() int
}

// getHTTPStatusCode returns http status code from error.
func getHTTPStatusCode(err error) int {
	if err == nil {
		return fasthttp.StatusOK
	}

	if e, ok := err.(errorCode); ok && e.Code() != 0 {
		return e.Code()
	}

	switch errors.Cause(err) {
	case ErrInvalidArgument:
		return fasthttp.StatusBadRequest
	case ErrAlreadyExists:
		return fasthttp.StatusBadRequest
	case ErrBadRequest:
		return fasthttp.StatusBadRequest
	case ErrNotFound:
		return fasthttp.StatusNotFound
	case ErrUserUnAuthorized:
		return fasthttp.StatusUnauthorized
	case InternalError:
		return fasthttp.StatusInternalServerError
	case ErrStatusForbidden:
		return fasthttp.StatusForbidden
	default:
		return fasthttp.StatusInternalServerError
	}
}
