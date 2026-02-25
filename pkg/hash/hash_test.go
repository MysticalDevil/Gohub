package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcryptHashAndCheck(t *testing.T) {
	password := "secret123"
	hashed := BcryptHash(password)
	require.Len(t, hashed, 60)
	require.True(t, BcryptCheck(password, hashed))
	require.False(t, BcryptCheck("wrong", hashed))
}

func TestBcryptIsHashed(t *testing.T) {
	require.True(t, BcryptIsHashed(BcryptHash("x")))
	require.False(t, BcryptIsHashed("plain"))
}
