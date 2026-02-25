package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	appconfig "gohub/config"
	pkgconfig "gohub/pkg/config"
)

func initTestConfig(t *testing.T) {
	t.Helper()

	if err := os.WriteFile(".env", []byte("APP_ENV=local\nAPP_URL=http://localhost:3000\nTIMEZONE=UTC\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(".env")
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
