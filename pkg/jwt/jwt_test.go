package jwt

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwtPkg "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	appconfig "gohub/config"
	pkgconfig "gohub/pkg/config"
)

func initJWTTestConfig(t *testing.T) {
	t.Helper()

	env := "APP_ENV=local\nAPP_DEBUG=false\nAPP_KEY=unit-test-key\nAPP_NAME=Gohub\nJWT_EXPIRE_TIME=1\nJWT_MAX_REFRESH_TIME=5\nTIMEZONE=UTC\n"
	tmpDir, err := os.MkdirTemp("", "gohub-jwt-test-")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte(env), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	if err := os.Setenv("APP_ENV_PATH", envPath); err != nil {
		t.Fatalf("set APP_ENV_PATH: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tmpDir)
		_ = os.Unsetenv("APP_ENV_PATH")
	})

	appconfig.Initialize()
	pkgconfig.InitConfig("")
}

func newGinContextWithToken(t *testing.T, token string) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c.Request = req
	return c
}

func TestIssueAndParseToken(t *testing.T) {
	initJWTTestConfig(t)

	jwt := NewJWT()
	tokenString := jwt.IssueToken("1", "alice")
	require.NotEmpty(t, tokenString)

	c := newGinContextWithToken(t, tokenString)
	claims, err := jwt.ParseToken(c)
	require.NoError(t, err)
	require.Equal(t, "1", claims.UserID)
	require.Equal(t, "alice", claims.UserName)
}

func TestParseTokenMissingHeader(t *testing.T) {
	initJWTTestConfig(t)

	jwt := NewJWT()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := jwt.ParseToken(c)
	require.ErrorIs(t, err, ErrHeaderEmpty)
}

func TestParseTokenMalformedHeader(t *testing.T) {
	initJWTTestConfig(t)

	jwt := NewJWT()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token abc")
	c.Request = req

	_, err := jwt.ParseToken(c)
	require.ErrorIs(t, err, ErrHeaderMalformed)
}

func TestRefreshTokenExpiredWithinMaxRefresh(t *testing.T) {
	initJWTTestConfig(t)

	jwt := NewJWT()
	now := time.Now().UTC()
	claims := CustomClaims{
		UserID:       "1",
		UserName:     "bob",
		ExpireAtTime: now.Add(-time.Minute).Unix(),
		RegisteredClaims: jwtPkg.RegisteredClaims{
			IssuedAt:  jwtPkg.NewNumericDate(now.Add(-time.Minute)),
			ExpiresAt: jwtPkg.NewNumericDate(now.Add(-time.Minute)),
			Issuer:    "Gohub",
		},
	}

	expiredToken, err := jwt.createToken(claims)
	require.NoError(t, err)

	c := newGinContextWithToken(t, expiredToken)
	newToken, err := jwt.RefreshToken(c)
	require.NoError(t, err)
	require.NotEmpty(t, newToken)
}

func TestRefreshTokenExpiredBeyondMaxRefresh(t *testing.T) {
	initJWTTestConfig(t)

	jwt := NewJWT()
	now := time.Now().UTC()
	claims := CustomClaims{
		UserID:       "1",
		UserName:     "bob",
		ExpireAtTime: now.Add(-time.Hour).Unix(),
		RegisteredClaims: jwtPkg.RegisteredClaims{
			IssuedAt:  jwtPkg.NewNumericDate(now.Add(-time.Hour)),
			ExpiresAt: jwtPkg.NewNumericDate(now.Add(-time.Hour)),
			Issuer:    "Gohub",
		},
	}

	expiredToken, err := jwt.createToken(claims)
	require.NoError(t, err)

	c := newGinContextWithToken(t, expiredToken)
	_, err = jwt.RefreshToken(c)
	require.ErrorIs(t, err, ErrTokenExpiredMaxRefresh)
}
