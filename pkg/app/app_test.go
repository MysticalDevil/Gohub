package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	appconfig "gohub/config"
	pkgconfig "gohub/pkg/config"
)

func initTestConfig(t *testing.T) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "gohub-app-test-")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("APP_ENV=local\nAPP_URL=http://localhost:3000\nTIMEZONE=UTC\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	if err := os.Setenv("APP_ENV_PATH", envPath); err != nil {
		t.Fatalf("set APP_ENV_PATH: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tmpDir)
		_ = os.Unsetenv("APP_ENV_PATH")
	})

	appconfig.Initialize()
	pkgconfig.InitConfig("")
}

func TestEnvironmentFlags(t *testing.T) {
	initTestConfig(t)

	require.True(t, IsLocal())
	require.False(t, IsProduction())
	require.False(t, IsTesting())
}

func TestURLs(t *testing.T) {
	initTestConfig(t)

	require.Equal(t, "http://localhost:3000/health", URL("/health"))
	require.Equal(t, "http://localhost:3000/v1/topics", V1URL("topics"))
}
