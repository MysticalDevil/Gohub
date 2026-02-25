package verifycode

type Store interface {
	// Set verify code
	Set(id string, value string) bool

	// Get verify code
	Get(id string, clear bool) string

	// Verify Check verify code
	Verify(id, answer string, clear bool) bool
}
