package routes_test

import (
	"net/http"
	"testing"

	"gohub/tests"
)

func TestCategoriesIndex(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	_ = tests.SeedCategory(t, tests.CategoryParams{Name: "cat1"})
	_ = tests.SeedCategory(t, tests.CategoryParams{Name: "cat2"})

	rec := tests.DoJSON(t, router, http.MethodGet, "/api/v1/categories", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCategoriesStoreUpdateDelete(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "catuser"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/categories", map[string]any{
		"name":        "cat",
		"description": "desc",
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	category := tests.SeedCategory(t, tests.CategoryParams{Name: "catx"})

	rec = tests.DoJSON(t, router, http.MethodPut, "/api/v1/categories/"+category.GetStringID(), map[string]any{
		"name":        "caty",
		"description": "new desc",
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	rec = tests.DoJSON(t, router, http.MethodDelete, "/api/v1/categories/"+category.GetStringID(), nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCategoriesDeleteNotFound(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "catuser"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodDelete, "/api/v1/categories/999999", nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
