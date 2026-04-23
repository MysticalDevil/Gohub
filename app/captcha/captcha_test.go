package captcha

import (
	"testing"

	"github.com/mojocn/base64Captcha"
	"github.com/stretchr/testify/require"
)

func TestCaptchaVerifyClearsAfterSuccess(t *testing.T) {
	store := NewMemoryStore()
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	c := base64Captcha.NewCaptcha(driver, store)

	id, _, _, err := c.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, id)

	answer := store.Get(id, false)
	require.NotEmpty(t, answer)

	// First verification with clear=true should succeed.
	require.True(t, c.Verify(id, answer, true))

	// Second verification should fail because the code was cleared.
	require.False(t, c.Verify(id, answer, true))
}
