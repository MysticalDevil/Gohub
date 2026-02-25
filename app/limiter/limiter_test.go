package limiter

import "testing"

func TestRouteToKeyString(t *testing.T) {
	cases := map[string]string{
		"/api/v1/topics":          "-api-v1-topics",
		"/api/v1/topics/:id":      "-api-v1-topics-_id",
		"/users/:user_id/links":   "-users-_user_id-links",
		"/health":                 "-health",
		"/":                       "-",
		"/api/v1/topics/:id/edit": "-api-v1-topics-_id-edit",
	}
	for input, expected := range cases {
		if got := routeToKeyString(input); got != expected {
			t.Fatalf("routeToKeyString(%q) = %q, want %q", input, got, expected)
		}
	}
}
