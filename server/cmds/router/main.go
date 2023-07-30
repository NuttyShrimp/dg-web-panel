package router

import (
	"degrens/panel/internal/admin"
	"degrens/panel/internal/auth"
	"degrens/panel/internal/cfx/characters"
	"degrens/panel/internal/config"
	"degrens/panel/internal/staff"
	"degrens/panel/internal/state"
	"degrens/panel/internal/users"
	"degrens/panel/lib/ratelimiter"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRouter(conf *config.Config) *gin.Engine {
	// Create a new gin Router
	r := gin.New()
	// TODO: set proxy when deploying
	err := r.SetTrustedProxies(nil)
	if err != nil {
		logrus.Fatal("Failed to set the trusted proxies", "error", err)
	}

	// Middlewares
	r.Use(
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
		}),
		cors.New(cors.Config{
			AllowOrigins:     conf.Server.Cors.Origins,
			AllowCredentials: true,
			AllowWebSockets:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "sentry-trace", "baggage", "X-Api-Key", "Upgrade"},
			MaxAge:           12 * time.Hour,
		}),
		ratelimiter.RateLimit(conf.Server.ReqPerSeq),
		gin.Logger(),
		gin.CustomRecovery(func(c *gin.Context, err any) {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(time.Second * 5)
			c.AbortWithStatus(http.StatusInternalServerError)
		}),
	)

	apiRG := r.Group("/api")

	// Register routes
	auth.NewAuthRouter(apiRG)

	securedapiRG := r.Group("/api", auth.NewMiddleWare())
	auth.NewSecuredAuthRouter(securedapiRG)
	users.NewUserRouter(securedapiRG)
	staff.NewStaffRouter(securedapiRG)
	admin.NewDevRouter(securedapiRG)
	characters.NewCharacterRouter(securedapiRG)
	state.NewStateRouter(securedapiRG)

	if conf.Server.Env != "production" {
		r.GET("/", func(ctx *gin.Context) {
			ctx.Redirect(307, "http://localhost:3001")
		})
	}

	return r
}
