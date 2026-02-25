package routes_test

import (
	"net/http"
	"testing"

	"gohub/tests"
)

func TestLinksIndex(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	_ = tests.SeedLink(t, tests.LinkParams{Name: "link1"})

	rec := tests.DoJSON(t, router, http.MethodGet, "/api/v1/links", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
