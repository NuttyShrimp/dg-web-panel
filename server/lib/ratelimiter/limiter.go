package ratelimiter

import (
	"github.com/didip/tollbooth/v6"
	"github.com/gin-gonic/gin"
)

func RateLimit(max float64) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(max, nil)
	// Switch around if not behind proxy
	lmt.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}
