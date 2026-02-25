// Package middlewares Storage system middleware
package middlewares

import (
	"bytes"
	"io"
	"log/slog"
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
			logFields = append(logFields, slog.String("request_body", string(requestBody)))

			// Response content
			logFields = append(logFields, slog.String("response_body", w.body.String()))
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
