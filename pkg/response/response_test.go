package response

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gohub/pkg/logger"
)

func TestDefaultMessage(t *testing.T) {
	require.Equal(t, "default", defaultMessage("default"))
	require.Equal(t, "custom", defaultMessage("default", "custom"))
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAbort404(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Abort404(c)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	logFile := filepath.Join(t.TempDir(), "test.log")
	logger.InitLogger(logFile, 1, 1, 1, false, "single", "error")

	BadRequest(c, errTest{})
	require.Equal(t, http.StatusBadRequest, w.Code)
}

type errTest struct{}

func (errTest) Error() string { return "boom" }
