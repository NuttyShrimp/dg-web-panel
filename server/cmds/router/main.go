package router

import (
	"degrens/panel/internal/admin"
	"degrens/panel/internal/auth"
	"degrens/panel/internal/cfx/characters"
	"degrens/panel/internal/config"
	"degrens/panel/internal/staff"
	"degrens/panel/internal/users"
	"degrens/panel/lib/log"
	"degrens/panel/lib/ratelimiter"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(config *config.Config, logger *log.Logger) *gin.Engine {
	// Create a new gin Router
	r := gin.New()
	// TODO: set proxy when deploying
	r.SetTrustedProxies(nil)

	// Middlewares
	r.Use(
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
		}),
		cors.New(cors.Config{
			AllowOrigins:     config.Server.Cors.Origins,
			AllowCredentials: true,
			AllowWebSockets:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "sentry-trace", "baggage", "X-Api-Key"},
			MaxAge:           12 * time.Hour,
		}),
		ratelimiter.RateLimit(config.Server.ReqPerSeq),
		gin.Logger(),
		gin.Recovery(),
	)

	apiRG := r.Group("/api")

	// Register routes
	auth.NewAuthRouter(apiRG, logger)

	securedapiRG := r.Group("/api", auth.NewMiddleWare(logger))
	auth.NewSecuredAuthRouter(securedapiRG, logger)
	users.NewUserRouter(securedapiRG, logger)
	staff.NewStaffRouter(securedapiRG, logger)
	admin.NewDevRouter(securedapiRG, logger)
	characters.NewCharacterRouter(securedapiRG, logger)

	return r
}
