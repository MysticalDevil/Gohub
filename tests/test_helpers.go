package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"gohub/app/models/category"
	"gohub/app/models/link"
	"gohub/app/models/topic"
	"gohub/app/models/user"
	"gohub/bootstrap"
	appconfig "gohub/config"
	_ "gohub/database/migrations"
	"gohub/pkg/config"
	"gohub/pkg/database"
	"gohub/pkg/logger"
	"gohub/pkg/migrate"
	"gohub/pkg/redis"
)

var (
	setupOnce sync.Once
	setupErr  error
)

func SetupTestEnv(t *testing.T) {
	t.Helper()
	setupOnce.Do(func() {
		gin.SetMode(gin.TestMode)

		env := fmt.Sprintf(
			"APP_ENV=testing\nAPP_DEBUG=false\nAPP_KEY=unit-test-key\nAPP_NAME=Gohub\nAPP_URL=http://localhost:3000\nAPP_PORT=3000\nTIMEZONE=UTC\nDB_CONNECTION=sqlite\nDB_SQL_FILE=/tmp/gohub_test.db\nDB_MAX_IDLE_CONNECTIONS=1\nDB_MAX_OPEN_CONNECTIONS=1\nDB_MAX_LIFE_SECONDS=5\nREDIS_HOST=127.0.0.1\nREDIS_PORT=6379\nREDIS_PASSWORD=\nREDIS_MAIN_DB=1\nREDIS_CACHE_DB=0\nJWT_EXPIRE_TIME=1\nJWT_MAX_REFRESH_TIME=5\nVERIFY_CODE_LENGTH=6\nVERIFY_CODE_EXPIRE=15\n",
		)

		tmpDir, err := os.MkdirTemp("", "gohub-test-")
		if err != nil {
			setupErr = err
			return
		}
		t.Cleanup(func() {
			_ = os.RemoveAll(tmpDir)
			_ = os.Unsetenv("APP_ENV_PATH")
			_ = os.Unsetenv("CONSOLE_SILENT")
			_ = os.Unsetenv("APP_ENV")
		})

		envPath := filepath.Join(tmpDir, ".env")
		if err := os.WriteFile(envPath, []byte(env), 0o644); err != nil {
			setupErr = err
			return
		}

		if err := os.Setenv("APP_ENV_PATH", envPath); err != nil {
			setupErr = err
			return
		}
		_ = os.Setenv("CONSOLE_SILENT", "1")
		_ = os.Setenv("APP_ENV", "testing")

		appconfig.Initialize()
		config.InitConfig("")
		logger.InitLogger(
			config.GetString("log.filename"),
			config.GetInt("log.max_size"),
			config.GetInt("log.max_backup"),
			config.GetInt("log.max_age"),
			config.GetBool("log.compress"),
			config.GetString("log.type"),
			config.GetString("log.level"),
		)
		bootstrap.SetupDB()
		bootstrap.SetupCache()

		migrate.NewMigrator().Up()
	})

	if setupErr != nil {
		t.Fatalf("setup test env failed: %v", setupErr)
	}
}

func ResetState(t *testing.T) {
	t.Helper()
	SetupTestEnv(t)

	if config.GetString("database.connection") == "sqlite" {
		if err := database.DB.Migrator().DropTable(
			&user.User{},
			&category.Category{},
			&topic.Topic{},
			&link.Link{},
			&migrate.Migration{},
		); err != nil {
			t.Fatalf("reset db failed: %v", err)
		}
	} else {
		if err := database.DeleteAllTables(); err != nil {
			t.Fatalf("reset db failed: %v", err)
		}
	}
	migrate.NewMigrator().Up()

	if redis.Redis != nil {
		redis.Redis.FlushDB()
	}

	_ = os.RemoveAll("public/uploads")
}

func NewRouter() *gin.Engine {
	router := gin.New()
	bootstrap.SetupRoute(router)
	return router
}

func DoJSON(t *testing.T, router http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body failed: %v", err)
		}
		reader = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if headers == nil || headers["User-Agent"] == "" {
		req.Header.Set("User-Agent", "go-test")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func DoMultipart(
	t *testing.T,
	router http.Handler,
	method, path string,
	fields map[string]string,
	files map[string]string,
	headers map[string]string,
) *httptest.ResponseRecorder {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("write field failed: %v", err)
		}
	}

	for field, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatalf("open file failed: %v", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(field, filepath.Base(filePath))
		if err != nil {
			t.Fatalf("create form file failed: %v", err)
		}

		if _, err := io.Copy(part, file); err != nil {
			t.Fatalf("copy file failed: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart failed: %v", err)
	}

	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if headers == nil || headers["User-Agent"] == "" {
		req.Header.Set("User-Agent", "go-test")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func DecodeJSON(t *testing.T, rec *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), target); err != nil {
		t.Fatalf("decode json failed: %v", err)
	}
}
