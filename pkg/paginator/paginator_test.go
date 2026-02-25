package paginator

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

func TestFormatBaseURL(t *testing.T) {
	initTestConfig(t)

	p := &Paginator{}
	require.Equal(t, "/topics?page=", p.formatBaseURL("/topics"))
	require.Equal(t, "/topics?order=asc&page=", p.formatBaseURL("/topics?order=asc"))
}

func TestGetPageLink(t *testing.T) {
	initTestConfig(t)

	p := &Paginator{
		BaseURL: "/topics?page=",
		Sort:    "created_at",
		Order:   "desc",
		PerPage: 20,
	}
	require.Equal(t, "/topics?page=2&sort=created_at&order=desc&per_page=20", p.getPageLink(2))
}

func TestPrevNextLinks(t *testing.T) {
	initTestConfig(t)

	p := &Paginator{
		BaseURL:   "/topics?page=",
		Sort:      "id",
		Order:     "asc",
		PerPage:   10,
		TotalPage: 3,
		Page:      2,
	}
	require.Equal(t, "/topics?page=1&sort=id&order=asc&per_page=10", p.getPrevPageURL())
	require.Equal(t, "/topics?page=3&sort=id&order=asc&per_page=10", p.getNextPageURL())
}

func TestPrevNextLinksBounds(t *testing.T) {
	initTestConfig(t)

	p := &Paginator{BaseURL: "/topics?page=", Sort: "id", Order: "asc", PerPage: 10, TotalPage: 1, Page: 1}
	require.Equal(t, "", p.getPrevPageURL())
	require.Equal(t, "", p.getNextPageURL())
}
