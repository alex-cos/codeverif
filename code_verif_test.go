package codeverif_test

import (
	"testing"
	"time"

	"github.com/alex-cos/codeverif"
	"github.com/stretchr/testify/assert"
)

func TestNominal(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(4, 5*time.Second, 2)

	code, err := cv.RequestCode("userID")
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	err = cv.VerifyCode("userID", code)
	assert.NoError(t, err)
}

func TestWithError(t *testing.T) {
	t.Parallel()

	cv := codeverif.New(4, 5*time.Second, 2)

	code, err := cv.RequestCode("userID")
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	err = cv.VerifyCode("userID", "12345")
	assert.Error(t, err)

	err = cv.VerifyCode("xxxxxx", code)
	assert.Error(t, err)

	code, err = cv.RequestCode("userID")
	assert.Error(t, err)
	assert.Empty(t, code)

	cv = codeverif.New(4, time.Second, 2)
	code, err = cv.RequestCode("userID")
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	time.Sleep(1100 * time.Millisecond)

	err = cv.VerifyCode("userID", code)
	assert.Error(t, err)

	cv = codeverif.New(4, time.Second, 2)
	code, err = cv.RequestCode("userID")
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	err = cv.VerifyCode("userID", "12345")
	assert.Error(t, err)
	err = cv.VerifyCode("userID", "12345")
	assert.Error(t, err)
	err = cv.VerifyCode("userID", code)
	assert.Error(t, err)
}
