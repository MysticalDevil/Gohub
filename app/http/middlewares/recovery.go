package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gohub/pkg/logger"
	"gohub/pkg/response"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

// Recovery Use zap.Error() to log panic and call stack
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get user request information
				httpRequest, _ := httputil.DumpRequest(c.Request, true)

				// When the link is interrupted,
				// it is normal behavior for the client to interrupt the connection,
				// and there is no need to record stack information.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") ||
							strings.Contains(errStr, "connect reset by peer") {
							brokenPipe = true
						}
					}
				}
				// In the event of a broken link
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					_ = c.Error(err.(error))
					c.Abort()
					// The link is broken, the status code cannot be written
					return
				}
				// If it is not a link break, start recording stack information
				logger.Error("recovery from panic",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.Stack("stacktrace"),
				)

				// return 500 status code
				response.Abort500(c)
			}
		}()
		c.Next()
	}
}
