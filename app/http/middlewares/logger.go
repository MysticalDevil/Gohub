// Package middlewares Storage system middleware
package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

var sensitiveFields = []string{
	"password", "new_password", "new_password_confirm", "old_password", "current_password",
	"token", "access_token", "refresh_token",
	"verify_code", "verification_code", "code", "captcha_answer",
}

func isSensitiveField(key string) bool {
	for _, f := range sensitiveFields {
		if f == key {
			return true
		}
	}
	return false
}

func sanitizeJSONValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		for k, v2 := range val {
			if isSensitiveField(strings.ToLower(k)) {
				val[k] = "[REDACTED]"
			} else {
				val[k] = sanitizeJSONValue(v2)
			}
		}
		return val
	case []any:
		for i, v2 := range val {
			val[i] = sanitizeJSONValue(v2)
		}
		return val
	default:
		return v
	}
}

func sanitizeBody(bodyStr string) string {
	bodyStr = strings.TrimSpace(bodyStr)
	if bodyStr == "" {
		return bodyStr
	}

	var data any
	if err := json.Unmarshal([]byte(bodyStr), &data); err != nil {
		return bodyStr
	}

	data = sanitizeJSONValue(data)
	b, err := json.Marshal(data)
	if err != nil {
		return bodyStr
	}
	return string(b)
}

// Logger Log request
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get response content
		w := &responseBodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: c.Writer,
		}
		c.Writer = w

		// Get request data
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Set start time
		start := time.Now()
		c.Next()

		// Logic to start logging
		cost := time.Since(start)
		responseStatus := c.Writer.Status()

		logFields := []slog.Attr{
			slog.Int("status", responseStatus),
			slog.String("request", c.Request.Method+""+c.Request.URL.String()),
			slog.String("query", c.Request.URL.RawQuery),
			slog.String("ip", c.ClientIP()),
			slog.String("user-agent", c.Request.UserAgent()),
			slog.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			slog.String("time", helpers.MicrosecondStr(cost)),
		}

		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			// Request content
			logFields = append(logFields, slog.String("request_body", sanitizeBody(string(requestBody))))

			// Response content
			logFields = append(logFields, slog.String("response_body", sanitizeBody(w.body.String())))
		}

		if responseStatus > 400 && responseStatus <= 499 {
			logger.Warn("HTTP Warning "+cast.ToString(responseStatus), logFields...)
		} else if responseStatus >= 500 && responseStatus <= 599 {
			logger.Error("HTTP error "+cast.ToString(responseStatus), logFields...)
		} else {
			logger.Info("HTTP Access Log", logFields...)
		}
	}
}
