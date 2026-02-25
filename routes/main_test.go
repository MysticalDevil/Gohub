package routes_test

import (
	"os"
	"path/filepath"
	"testing"
)

func repoRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return ""
		}
		wd = parent
	}
}

func TestMain(m *testing.M) {
	originalWD, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	if root := repoRoot(); root != "" {
		_ = os.Chdir(root)
	}

	code := m.Run()

	_ = os.Chdir(originalWD)
	os.Exit(code)
}
