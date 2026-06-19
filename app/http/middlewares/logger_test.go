package middlewares

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeBodyRedactsSensitiveFields(t *testing.T) {
	input := `{"password":"secret123","token":"abc","email":"user@example.com","nested":{"verify_code":"123456","code":"456"}}`
	result := sanitizeBody(input)

	require.Contains(t, result, "[REDACTED]")
	require.NotContains(t, result, "secret123")
	require.NotContains(t, result, `"token":"abc"`)
	require.NotContains(t, result, "123456")
	require.Contains(t, result, "user@example.com")
}

func TestSanitizeBodyHandlesNestedArrays(t *testing.T) {
	input := `[{"password":"p1"},{"access_token":"at"}]`
	result := sanitizeBody(input)
	require.Contains(t, result, "[REDACTED]")
	require.NotContains(t, result, `"p1"`)
	require.NotContains(t, result, `"at"`)
}

func TestSanitizeBodyIgnoresNonJSON(t *testing.T) {
	input := `plain text body`
	require.Equal(t, input, sanitizeBody(input))
}

func TestSanitizeBodyEmpty(t *testing.T) {
	require.Equal(t, "", sanitizeBody(""))
}

func TestSanitizeRequestForLogRedactsHeadersAndOmitsBody(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodPost,
		"/login",
		io.NopCloser(strings.NewReader(`{"password":"secret","verify_code":"123456"}`)),
	)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Cookie", "session=abc")
	req.Header.Set("Content-Type", "application/json")

	result := sanitizeRequestForLog(req)

	require.Contains(t, result, "Authorization: [REDACTED]")
	require.Contains(t, result, "Cookie: [REDACTED]")
	require.NotContains(t, result, "Bearer token")
	require.NotContains(t, result, "session=abc")
	require.NotContains(t, result, "secret")
	require.NotContains(t, result, "123456")
}
