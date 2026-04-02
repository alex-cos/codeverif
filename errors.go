package codeverif

import "errors"

// Sentinel errors for programmatic error handling.
var (
	ErrThrottled       = errors.New("request throttled; wait before requesting another code")
	ErrCodeExpired     = errors.New("code expired")
	ErrTooManyAttempts = errors.New("too many attempts")
	ErrInvalidCode     = errors.New("invalid code")
	ErrInvalidLength   = errors.New("invalid code length")
)
