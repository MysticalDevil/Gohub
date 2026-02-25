package routes_test

import (
	"fmt"
	"net/http"
	"testing"

	"gohub/tests"
)

func TestTopicsIndexAndShow(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "topicuser"})
	category := tests.SeedCategory(t, tests.CategoryParams{Name: "topiccat"})
	topic := tests.SeedTopic(t, user, category, tests.TopicParams{Title: "topic", Body: "topic body content"})

	rec := tests.DoJSON(t, router, http.MethodGet, "/api/v1/topics", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	rec = tests.DoJSON(t, router, http.MethodGet, "/api/v1/topics/"+topic.GetStringID(), nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestTopicsStoreUpdateDelete(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "owner"})
	category := tests.SeedCategory(t, tests.CategoryParams{Name: "cat1"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/topics", map[string]any{
		"title":       "new topic",
		"body":        "this is a new topic body",
		"category_id": category.GetStringID(),
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	topic := tests.SeedTopic(t, user, category, tests.TopicParams{Title: "old", Body: "old body content"})

	rec = tests.DoJSON(t, router, http.MethodPut, "/api/v1/topics/"+topic.GetStringID(), map[string]any{
		"title":       "updated",
		"body":        "updated body content",
		"category_id": category.GetStringID(),
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	rec = tests.DoJSON(t, router, http.MethodDelete, "/api/v1/topics/"+topic.GetStringID(), nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestTopicsUpdateForbidden(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	owner := tests.SeedUser(t, tests.UserParams{Name: "owner"})
	other := tests.SeedUser(t, tests.UserParams{Name: "other"})
	category := tests.SeedCategory(t, tests.CategoryParams{Name: "cat2"})
	topic := tests.SeedTopic(t, owner, category, tests.TopicParams{Title: "topic", Body: "topic body content"})

	otherToken := tests.IssueToken(other)

	rec := tests.DoJSON(t, router, http.MethodPut, "/api/v1/topics/"+topic.GetStringID(), map[string]any{
		"title":       "updated",
		"body":        "updated body content",
		"category_id": category.GetStringID(),
	}, map[string]string{
		"Authorization": "Bearer " + otherToken,
	})
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestTopicsDeleteNotFound(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "owner"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodDelete, "/api/v1/topics/999999", nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestTopicsStoreMissingCategory(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "owner"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/topics", map[string]any{
		"title":       "new topic",
		"body":        "this is a new topic body",
		"category_id": fmt.Sprintf("%d", 99999),
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}
