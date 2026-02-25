package routes_test

import (
	"context"
	"net/http"
	"testing"

	"gohub/pkg/auth"
	"gohub/tests"
)

func TestAuthSignupPhoneExist(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Phone: "00012345678"})

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/signup/phone/exist", map[string]any{
		"phone": user.Phone,
	}, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload map[string]any
	tests.DecodeJSON(t, rec, &payload)
	data, _ := payload["data"].(map[string]any)
	if data["exist"] != true {
		t.Fatalf("expected exist=true, got %v", data["exist"])
	}
}

func TestAuthSignupEmailExist(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{Email: "exists@testing.com"})

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/signup/email/exist", map[string]any{
		"email": user.Email,
	}, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload map[string]any
	tests.DecodeJSON(t, rec, &payload)
	data, _ := payload["data"].(map[string]any)
	if data["exist"] != true {
		t.Fatalf("expected exist=true, got %v", data["exist"])
	}
}

func TestAuthSignupUsingPhone(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/signup/using-phone", map[string]any{
		"phone":            "00012345678",
		"verify_code":      "123456",
		"name":             "phoneuser",
		"password":         "password123",
		"password_confirm": "password123",
	}, nil)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var payload map[string]any
	tests.DecodeJSON(t, rec, &payload)
	data, _ := payload["data"].(map[string]any)
	if data["token"] == "" {
		t.Fatalf("expected token in response")
	}
}

func TestAuthSignupUsingEmail(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/signup/using-email", map[string]any{
		"email":            "user@testing.com",
		"verify_code":      "123456",
		"name":             "emailuser",
		"password":         "password123",
		"password_confirm": "password123",
	}, nil)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var payload map[string]any
	tests.DecodeJSON(t, rec, &payload)
	data, _ := payload["data"].(map[string]any)
	if data["token"] == "" {
		t.Fatalf("expected token in response")
	}
}

func TestAuthVerifyCodes(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/verify-codes/captcha", map[string]any{}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var captchaPayload map[string]any
	tests.DecodeJSON(t, rec, &captchaPayload)
	data, _ := captchaPayload["data"].(map[string]any)
	if data["captcha_id"] == "" || data["captcha_image"] == "" {
		t.Fatalf("expected captcha fields in response")
	}

	rec = tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/verify-codes/phone", map[string]any{
		"phone":          "00012345678",
		"captcha_id":     "captcha_skip_test",
		"captcha_answer": "123456",
	}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	rec = tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/verify-codes/email", map[string]any{
		"email":          "user@testing.com",
		"captcha_id":     "captcha_skip_test",
		"captcha_answer": "123456",
	}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAuthLoginAndRefresh(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{
		Name:     "loginuser",
		Phone:    "00012345678",
		Email:    "login@testing.com",
		Password: "password123",
	})

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/login/using-phone", map[string]any{
		"phone":       user.Phone,
		"verify_code": "123456",
	}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	rec = tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/login/using-password", map[string]any{
		"login_id":       user.Name,
		"password":       "password123",
		"captcha_id":     "captcha_skip_test",
		"captcha_answer": "123456",
	}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var loginPayload map[string]any
	tests.DecodeJSON(t, rec, &loginPayload)
	data, _ := loginPayload["data"].(map[string]any)
	token, _ := data["token"].(string)
	if token == "" {
		t.Fatalf("expected token")
	}

	rec = tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/login/refresh-token", nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAuthPasswordReset(t *testing.T) {
	tests.ResetState(t)
	router := tests.NewRouter()

	user := tests.SeedUser(t, tests.UserParams{
		Phone:    "00012345678",
		Email:    "reset@testing.com",
		Password: "password123",
	})

	rec := tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/password-reset/using-phone", map[string]any{
		"phone":       user.Phone,
		"verify_code": "123456",
		"password":    "newpassword123",
	}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if _, err := auth.Attempt(context.Background(), user.Phone, "newpassword123"); err != nil {
		t.Fatalf("expected password reset to succeed: %v", err)
	}

	rec = tests.DoJSON(t, router, http.MethodPost, "/api/v1/auth/password-reset/using-email", map[string]any{
		"email":       user.Email,
		"verify_code": "123456",
		"password":    "newpassword456",
	}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if _, err := auth.Attempt(context.Background(), user.Email, "newpassword456"); err != nil {
		t.Fatalf("expected password reset to succeed: %v", err)
	}
}
