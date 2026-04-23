package verifycode

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	_ "gohub/config"
	"gohub/pkg/config"
	"gohub/pkg/mail"
)

type mockMailDriver struct {
	shouldFail bool
}

func (m *mockMailDriver) Send(email mail.Email, config map[string]string) bool {
	return !m.shouldFail
}

func setupConfig(t *testing.T) {
	t.Helper()
	env := "APP_ENV=testing\nAPP_DEBUG=false\nAPP_KEY=unit-test-key\nAPP_NAME=Gohub\nAPP_URL=http://localhost:3000\nAPP_PORT=3000\nTIMEZONE=UTC\nVERIFY_CODE_LENGTH=6\nVERIFY_CODE_EXPIRE=15\n"

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte(env), 0o644); err != nil {
		t.Fatalf("write env failed: %v", err)
	}

	_ = os.Setenv("APP_ENV_PATH", envPath)
	_ = os.Setenv("APP_ENV", "testing")

	config.InitConfig("")
}

func TestVerifyCodeCheckAnswerClearsAfterSuccess(t *testing.T) {
	store := NewMemoryStore()
	store.Set("test@example.com", "123456")

	// First verification with clear=true should succeed.
	require.True(t, store.Verify("test@example.com", "123456", true))

	// Second verification should fail because the code was cleared.
	require.False(t, store.Verify("test@example.com", "123456", true))
}

func TestSendEmailReturnsErrorOnFailure(t *testing.T) {
	setupConfig(t)

	originalDriver := mail.NewMailer().Driver
	mail.NewMailer().SetDriver(&mockMailDriver{shouldFail: true})
	t.Cleanup(func() {
		mail.NewMailer().SetDriver(originalDriver)
	})

	vc := &VerifyCode{Store: NewMemoryStore()}
	err := vc.SendEmail("fail@example.com")
	require.Error(t, err)
}

func TestSendEmailSucceeds(t *testing.T) {
	setupConfig(t)

	originalDriver := mail.NewMailer().Driver
	mail.NewMailer().SetDriver(&mockMailDriver{shouldFail: false})
	t.Cleanup(func() {
		mail.NewMailer().SetDriver(originalDriver)
	})

	vc := &VerifyCode{Store: NewMemoryStore()}
	err := vc.SendEmail("ok@example.com")
	require.NoError(t, err)
}
