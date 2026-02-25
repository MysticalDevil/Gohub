package helpers

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	require.True(t, Empty(""))
	require.False(t, Empty("x"))
	require.True(t, Empty([]string{}))
	require.False(t, Empty([]string{"a"}))
	require.True(t, Empty(0))
	require.False(t, Empty(1))
	require.True(t, Empty(false))
	require.False(t, Empty(true))
}

func TestMicrosecondStr(t *testing.T) {
	require.Equal(t, "1.500ms", MicrosecondStr(1500*time.Microsecond))
}

func TestRandomNumber(t *testing.T) {
	s := RandomNumber(8)
	require.Len(t, s, 8)
	require.Empty(t, strings.Trim(s, "0123456789"))
}

func TestRandomString(t *testing.T) {
	s := RandomString(12)
	require.Len(t, s, 12)
}

func TestFirstElement(t *testing.T) {
	require.Equal(t, "", FirstElement([]string{}))
	require.Equal(t, "a", FirstElement([]string{"a"}))
}
