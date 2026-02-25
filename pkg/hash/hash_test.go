package hash

import "testing"

func TestBcryptHashAndCheck(t *testing.T) {
	password := "secret123"
	hashed := BcryptHash(password)
	if len(hashed) != 60 {
		t.Fatalf("unexpected hash length: %d", len(hashed))
	}
	if !BcryptCheck(password, hashed) {
		t.Fatalf("expected password to match hash")
	}
	if BcryptCheck("wrong", hashed) {
		t.Fatalf("expected password mismatch")
	}
}

func TestBcryptIsHashed(t *testing.T) {
	if !BcryptIsHashed(BcryptHash("x")) {
		t.Fatalf("expected hash to be detected")
	}
	if BcryptIsHashed("plain") {
		t.Fatalf("did not expect plain string to be treated as hash")
	}
}
