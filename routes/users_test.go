package routes_test

import (
	"net/http"
	"testing"

	"gohub/pkg/auth"
	"gohub/tests"
)

func TestUsersCurrentUserUnauthorized(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	rec := tests.DoJSON(t, router, http.MethodGet, "/api/v1/user", nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestUsersCurrentUser(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "currentuser"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodGet, "/api/v1/user", nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUsersIndex(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	_ = tests.SeedUser(t, tests.UserParams{Name: "userone"})
	_ = tests.SeedUser(t, tests.UserParams{Name: "usertwo"})

	rec := tests.DoJSON(t, router, http.MethodGet, "/api/v1/users", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUsersUpdateProfile(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "profileuser"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodPut, "/api/v1/users", map[string]any{
		"name":         "newname",
		"city":         "beijing",
		"introduction": "hello world",
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUsersUpdateEmail(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "emailuser"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodPut, "/api/v1/users/email", map[string]any{
		"email":       "updated@testing.com",
		"verify_code": "123456",
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUsersUpdatePhone(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "phoneuser"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodPut, "/api/v1/users/phone", map[string]any{
		"phone":       "00012345679",
		"verify_code": "123456",
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUsersUpdatePassword(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "passuser", Password: "password123"})
	token := tests.IssueToken(user)

	rec := tests.DoJSON(t, router, http.MethodPut, "/api/v1/users/password", map[string]any{
		"password":             "password123",
		"new_password":         "password456",
		"new_password_confirm": "password456",
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if _, err := auth.Attempt(user.Name, "password456"); err != nil {
		t.Fatalf("expected new password to work: %v", err)
	}
}

func TestUsersUpdateAvatar(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Name: "avataruser"})
	token := tests.IssueToken(user)
	filePath := tests.CreatePNGFile(t)

	rec := tests.DoMultipart(t, router, http.MethodPut, "/api/v1/users/avatar", nil, map[string]string{
		"avatar": filePath,
	}, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
