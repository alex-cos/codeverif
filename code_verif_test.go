package codeverif_test

import (
	"sync"
	"testing"
	"time"

	"github.com/alex-cos/codeverif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNominal(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(4, 5*time.Second, 2)

	code, err := cv.RequestCode("userID")
	require.NoError(t, err)
	assert.NotEmpty(t, code)

	err = cv.VerifyCode("userID", code)
	assert.NoError(t, err)
}

func TestRequestCodeThrottled(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(4, 5*time.Second, 2)

	_, err := cv.RequestCode("userID")
	require.NoError(t, err)

	_, err = cv.RequestCode("userID")
	assert.ErrorIs(t, err, codeverif.ErrThrottled)
}

func TestVerifyCodeExpired(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(4, 100*time.Millisecond, 2)
	code, err := cv.RequestCode("userID")
	require.NoError(t, err)

	time.Sleep(150 * time.Millisecond)

	err = cv.VerifyCode("userID", code)
	assert.ErrorIs(t, err, codeverif.ErrCodeExpired)
}

func TestVerifyCodeNotFound(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(4, 5*time.Second, 2)

	err := cv.VerifyCode("nonexistent", "1234")
	assert.ErrorIs(t, err, codeverif.ErrCodeExpired)
}

func TestVerifyInvalidCode(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(4, 5*time.Second, 3)
	code, err := cv.RequestCode("userID")
	require.NoError(t, err)

	err = cv.VerifyCode("userID", "0000")
	assert.ErrorIs(t, err, codeverif.ErrInvalidCode)

	err = cv.VerifyCode("userID", code)
	assert.NoError(t, err)
}

func TestVerifyTooManyAttempts(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(4, 5*time.Second, 2)
	code, err := cv.RequestCode("userID")
	require.NoError(t, err)

	err = cv.VerifyCode("userID", "0000")
	assert.ErrorIs(t, err, codeverif.ErrInvalidCode)

	err = cv.VerifyCode("userID", "0000")
	assert.ErrorIs(t, err, codeverif.ErrInvalidCode)

	err = cv.VerifyCode("userID", "0000")
	assert.ErrorIs(t, err, codeverif.ErrTooManyAttempts)

	err = cv.VerifyCode("userID", code)
	assert.ErrorIs(t, err, codeverif.ErrTooManyAttempts)
}

func TestCodeFormat(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(6, 5*time.Second, 2)
	code, err := cv.RequestCode("userID")
	require.NoError(t, err)

	assert.Len(t, code, 6)
	for _, c := range code {
		assert.True(t, c >= '0' && c <= '9')
	}
}

func TestInvalidCodeLength(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(1, 5*time.Second, 2)
	_, err := cv.RequestCode("userID")
	assert.ErrorIs(t, err, codeverif.ErrInvalidLength)
}

func TestConcurrentVerify(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(4, 5*time.Second, 100)
	code, err := cv.RequestCode("userID")
	require.NoError(t, err)

	var wg sync.WaitGroup
	var successCount int
	var mu sync.Mutex

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := cv.VerifyCode("userID", code)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, 1, successCount)
}
