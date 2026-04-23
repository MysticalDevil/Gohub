package middlewares

import (
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
