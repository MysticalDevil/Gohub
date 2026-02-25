package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPutAndExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	if Exists(path) {
		t.Fatalf("did not expect file to exist")
	}
	if err := Put([]byte("hello"), path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !Exists(path) {
		t.Fatalf("expected file to exist")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected file content: %s", string(data))
	}
}

func TestNameWithoutExtension(t *testing.T) {
	if got := NameWithoutExtension("photo.png"); got != "photo" {
		t.Fatalf("unexpected name: %s", got)
	}
	if got := NameWithoutExtension("archive.tar.gz"); got != "archive.tar" {
		t.Fatalf("unexpected name: %s", got)
	}
}
