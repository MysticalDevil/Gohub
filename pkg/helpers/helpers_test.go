package helpers

import (
	"strings"
	"testing"
	"time"
)

func TestEmpty(t *testing.T) {
	if !Empty("") {
		t.Fatalf("expected empty string")
	}
	if Empty("x") {
		t.Fatalf("expected non-empty string")
	}
	if !Empty([]string{}) {
		t.Fatalf("expected empty slice")
	}
	if Empty([]string{"a"}) {
		t.Fatalf("expected non-empty slice")
	}
	if !Empty(0) {
		t.Fatalf("expected empty int")
	}
	if Empty(1) {
		t.Fatalf("expected non-empty int")
	}
	if !Empty(false) {
		t.Fatalf("expected empty bool")
	}
	if Empty(true) {
		t.Fatalf("expected non-empty bool")
	}
}

func TestMicrosecondStr(t *testing.T) {
	if got := MicrosecondStr(1500 * time.Microsecond); got != "1.500ms" {
		t.Fatalf("unexpected format: %s", got)
	}
}

func TestRandomNumber(t *testing.T) {
	s := RandomNumber(8)
	if len(s) != 8 {
		t.Fatalf("unexpected length: %d", len(s))
	}
	if strings.Trim(s, "0123456789") != "" {
		t.Fatalf("expected numeric string, got %q", s)
	}
}

func TestRandomString(t *testing.T) {
	s := RandomString(12)
	if len(s) != 12 {
		t.Fatalf("unexpected length: %d", len(s))
	}
}

func TestFirstElement(t *testing.T) {
	if got := FirstElement([]string{}); got != "" {
		t.Fatalf("expected empty result")
	}
	if got := FirstElement([]string{"a"}); got != "a" {
		t.Fatalf("unexpected result: %s", got)
	}
}
