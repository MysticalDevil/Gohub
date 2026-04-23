package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testLink struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func TestGetObjectDeserializesCorrectly(t *testing.T) {
	InitWithCacheStore(NewMemoryStore())

	links := []testLink{
		{ID: 1, Name: "Google", URL: "https://google.com"},
		{ID: 2, Name: "GitHub", URL: "https://github.com"},
	}

	Set("links:test", links, time.Minute)

	var wanted []testLink
	GetObject("links:test", &wanted)

	require.Len(t, wanted, 2)
	require.Equal(t, "Google", wanted[0].Name)
	require.Equal(t, "GitHub", wanted[1].Name)
}

func TestGetObjectMissingKey(t *testing.T) {
	InitWithCacheStore(NewMemoryStore())

	var wanted []testLink
	GetObject("links:missing", &wanted)
	require.Empty(t, wanted)
}
