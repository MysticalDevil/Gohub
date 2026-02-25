package paginator

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	appconfig "gohub/config"
	pkgconfig "gohub/pkg/config"
)

func initTestConfig(t *testing.T) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "gohub-paginator-test-")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("APP_ENV=local\n"), 0o644); err != nil {
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

func TestGetLimitDefaults(t *testing.T) {
	initTestConfig(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	p := &Paginator{ctx: c}
	require.Equal(t, 10, p.getLimit(0))
}

func TestGetLimitFromQuery(t *testing.T) {
	initTestConfig(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?limit=25", nil)

	p := &Paginator{ctx: c}
	require.Equal(t, 25, p.getLimit(0))
}

func TestGetOffset(t *testing.T) {
	initTestConfig(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?offset=30", nil)

	p := &Paginator{ctx: c}
	require.Equal(t, 30, p.getOffset())
}

func TestGetOffsetNegative(t *testing.T) {
	initTestConfig(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?offset=-5", nil)

	p := &Paginator{ctx: c}
	require.Equal(t, 0, p.getOffset())
}
