package response

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"gohub/pkg/logger"
)

func TestDefaultMessage(t *testing.T) {
	if got := defaultMessage("default"); got != "default" {
		t.Fatalf("unexpected default message: %s", got)
	}
	if got := defaultMessage("default", "custom"); got != "custom" {
		t.Fatalf("unexpected override message: %s", got)
	}
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestAbort404(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Abort404(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	logFile := filepath.Join(t.TempDir(), "test.log")
	logger.InitLogger(logFile, 1, 1, 1, false, "single", "error")

	BadRequest(c, errTest{})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

type errTest struct{}

func (errTest) Error() string { return "boom" }
