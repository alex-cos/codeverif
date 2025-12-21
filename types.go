package codeverif

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// VerifCode exposes public methods to request and verify codes.
type VerifCode struct {
	codeLength  int           // number of digits in numeric code
	expiry      time.Duration // how long codes are valid
	maxAttempts int           // allowed wrong attempts before lock
	cache       *cache.Cache
}

// record represents a stored code.
type record struct {
	Code          string
	Attempts      int
	LastAttemptAt time.Time
}
