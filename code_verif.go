package codeverif

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/patrickmn/go-cache"
)

// New generates a new manager instance.
func New(codeLength int, expiry time.Duration, maxAttempts int) *VerifCode {
	return &VerifCode{
		codeLength:  codeLength,
		expiry:      expiry,
		maxAttempts: maxAttempts,
		cache:       cache.New(expiry, time.Hour),
	}
}

// RequestCode generates a code for a given userID and save it for expiry duration.
func (thiz *VerifCode) RequestCode(userID string) (string, error) {
	thiz.mu.Lock()
	defer thiz.mu.Unlock()

	code, err := generateNumericCode(thiz.codeLength)
	if err != nil {
		return "", err
	}
	rec := &record{
		Code:          code,
		Attempts:      0,
		LastAttemptAt: time.Unix(0, 0),
	}
	err = thiz.cache.Add(userID, rec, thiz.expiry)
	if err != nil {
		return "", ErrThrottled
	}

	return code, nil
}

// VerifyCode checks the submitted code. On success it deletes the stored code.
func (thiz *VerifCode) VerifyCode(userID, code string) error {
	thiz.mu.Lock()
	defer thiz.mu.Unlock()

	val, expiration, ok := thiz.cache.GetWithExpiration(userID)
	if !ok {
		return ErrCodeExpired
	}
	rec := val.(*record) //nolint:forcetypeassert

	rec.Attempts++
	rec.LastAttemptAt = time.Now().UTC()
	d := time.Until(expiration)
	if d <= 0 {
		return ErrCodeExpired
	}
	err := thiz.cache.Replace(userID, rec, time.Until(expiration))
	if err != nil {
		return err
	}
	if rec.Attempts > thiz.maxAttempts {
		return ErrTooManyAttempts
	}
	if rec.Code != code {
		return ErrInvalidCode
	}
	thiz.cache.Delete(userID)

	return nil
}

// Unexported methods

func generateNumericCode(n int) (string, error) {
	if n <= 2 {
		return "", ErrInvalidLength
	}
	maximum := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil) // 10^n
	num, err := rand.Int(rand.Reader, maximum)
	if err != nil {
		return "", err
	}
	format := fmt.Sprintf("%%0%dd", n)
	return fmt.Sprintf(format, num.Int64()), nil
}
