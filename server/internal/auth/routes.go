package auth

import (
	"context"
	"degrens/panel/internal/auth/apikeys"
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/discord"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/storage"
	"degrens/panel/internal/users"
	"degrens/panel/lib/errors"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	routes.Router
}

type SecuredAuthRouter struct {
	routes.Router
}

type NewAPIKeyBody struct {
	Comment string `json:"comment"`
	// In minutes
	Duration int64 `json:"duration"`
}

type DeleteAPIKeyBody struct {
	Keys []string `json:"keys"`
}

func NewAuthRouter(rg *gin.RouterGroup, logger *log.Logger) {
	router := &AuthRouter{
		routes.Router{
			RouterGroup: rg.Group("/auth"),
			Logger:      *logger,
		},
	}
	router.RegisterRoutes()
}

func (AR *AuthRouter) RegisterRoutes() {
	AR.RouterGroup.POST("/login", AR.loginHandler())
	AR.RouterGroup.POST("/logout", AR.logoutHandler())
	// TODO:  Check if this should be moved to the secured router
	AR.RouterGroup.POST("/refresh", AR.RefreshHandler)
	// Discord oauth callback
	AR.RouterGroup.GET("/callback", AR.discordCallbackHandler())
}

func (AR *AuthRouter) RefreshHandler(c *gin.Context) {
	userInfo, err := authinfo.GetUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to retrieve cookie",
		})
		return
	}
	if &userInfo == nil {
		// No session cookie
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No session cookie",
		})
		return
	}
	switch userInfo.AuthMethod {
	case authinfo.Discord:
		{
			isExpired := discord.IsTokenExpired(userInfo.ID)
			c.JSON(http.StatusOK, gin.H{
				"isExpired": isExpired,
			})
			if !isExpired {
				discord.RefreshToken(userInfo.ID)
			}
			return
		}
	default:
		{
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No refresh method for this auth method",
			})
			return
		}
	}
}

func (AR *AuthRouter) loginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authType := c.Query("type")
		switch authType {
		case "discord":
			state, err := discord.GenerateOAuthState(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Could not generate state",
				})
				return
			}
			// Set the state in the cookies
			storage.AddPublicCookie(c, "state", state)
			// Redirect to discord auth
			// Unable to redirect directly to discord, so we let frontend handle it
			c.JSON(http.StatusOK, gin.H{
				"url": discord.GetOAuthConf().AuthCodeURL(state),
			})
			return
		}
	}
}

func (AR *AuthRouter) logoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfoPtr, exists := c.Get("userInfo")
		userInfo := userInfoPtr.(*authinfo.AuthInfo)
		if exists == false {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get your details",
			})
			return
		}
		switch userInfo.AuthMethod {
		case authinfo.Discord:
			{
				discord.RemoveUserTokens(userInfo.ID)
				break
			}
		case authinfo.CFXToken:
			{
				// Send request to Cfx API to invalidate login token
				break
			}
		default:
			{
				AR.Logger.Error("Logout procedure failed to determine authentication type", "info", fmt.Sprintf("%+v", userInfo))
				break
			}
		}
		storage.RemoveCookie(c, "userInfo")
		c.String(http.StatusOK, "")
		return
	}
}

func (AR *AuthRouter) discordCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the state from the cookies
		var state string
		if err := storage.GetPublicCookie(c, "state", &state); &state == nil || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get authentication state",
			})
			return
		}
		// We check if the state is known to us
		if c.Query("state") != state {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid authentication state",
			})
			return
		}
		storage.RemoveCookie(c, "state")
		// Step 3: We exchange the code we got for an access token
		// Then we can use the access token to do actions, limited to scopes we requested
		token, err := discord.GetOAuthConf().Exchange(context.Background(), c.Query("code"))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "We encountered an error while trying to get your information",
			})
			AR.Logger.Error("Error while exchanging code for token", err)
			return
		}

		identity, err := discord.GetUserInfoViaToken(token)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get discord information, try again",
			})
			return
		}

		roles := discord.GetRegisterdRolesForIdentity(*identity)
		if len(roles) == 0 {
			c.Redirect(http.StatusTemporaryRedirect, "/errors/403")
			return
		}

		user := discord.UpdateUserInfo(*identity)
		discord.AssignTokenToUser(user.ID, token)

		// Set cookie with userId and roles
		authInfo := authinfo.AuthInfo{
			ID:         user.ID,
			Roles:      roles,
			AuthMethod: authinfo.Discord,
		}
		authInfo.Assign(c)

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func NewSecuredAuthRouter(rg *gin.RouterGroup, logger *log.Logger) {
	router := &SecuredAuthRouter{
		routes.Router{
			RouterGroup: rg.Group("/auth"),
			Logger:      *logger,
		},
	}
	router.RegisterRoutes()
}

func (SAR *SecuredAuthRouter) RegisterRoutes() {
	SAR.RouterGroup.GET("/role", SAR.roleCheckHandler())

	SAR.RouterGroup.GET("/apikey", SAR.fetchAPIKeys())
	SAR.RouterGroup.POST("/apikey", SAR.handleAPIKeyCreation())
	SAR.RouterGroup.DELETE("/apikey", SAR.handleAPIKeyDeletion())

	SAR.RouterGroup.GET("/apikey/all", SAR.fetchAllAPIKeys())
}

func (SAR *SecuredAuthRouter) roleCheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		target := ctx.Query("role")
		if target == "" {
			SAR.Logger.Error("Failed to parse role in /auth/role request", "target", target)
			ctx.JSON(400, models.RouteErrorMessage{
				Title:       "Request error",
				Description: "Encountered an error while trying to check access to a secure page",
			})
			return
		}
		hasAccess, err := users.HasRoleAccess(ctx, target)
		if err != nil {
			SAR.Logger.Error("Failed to check role in /auth/role request", "error", err)
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Server error",
				Description: "We encountered an error while trying to check the allowance to a page",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"access": hasAccess,
		})
	}
}

func (SAR *SecuredAuthRouter) fetchAPIKeys() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxUserInfo, exists := ctx.Get("userInfo")
		userInfo := ctxUserInfo.(*authinfo.AuthInfo)
		if exists == false {
			SAR.Logger.Error("Failed to get userinfo while fetching API keys")
			ctx.JSON(403, errors.Unauthorized)
			return
		}
		apiKeys := apikeys.GetAPIKeys(userInfo.ID)
		ctx.JSON(200, gin.H{
			"keys": apiKeys,
		})
	}
}

func (SAR *SecuredAuthRouter) handleAPIKeyCreation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body NewAPIKeyBody
		err := ctx.Copy().ShouldBindJSON(&body)
		if err != nil {
			SAR.Logger.Error("Failed to parse body in APIKey creation request", "error", err)
			ctx.JSON(500, errors.BodyParsingFailed)
			return
		}
		ctxUserInfo, exists := ctx.Get("userInfo")
		userInfo := ctxUserInfo.(*authinfo.AuthInfo)
		if exists == false {
			SAR.Logger.Error("Failed to get userinfo in request trying to make an API key")
			ctx.JSON(403, errors.Unauthorized)
			return
		}
		var key string
		key, err = apikeys.CreateAPIKey(userInfo.ID, body.Comment, time.Duration(body.Duration)*time.Minute)
		ctx.JSON(200, gin.H{
			"key": key,
		})
	}
}

func (SAR *SecuredAuthRouter) handleAPIKeyDeletion() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body DeleteAPIKeyBody
		err := ctx.ShouldBindJSON(&body)
		if err != nil {
			SAR.Logger.Error("Failed to parse body in APIKey deletion request", "error", err)
			ctx.JSON(500, errors.BodyParsingFailed)
			return
		}
		ctxUserInfo, exists := ctx.Get("userInfo")
		userInfo := ctxUserInfo.(*authinfo.AuthInfo)
		if exists == false {
			SAR.Logger.Error("Failed to get userinfo in request trying to make an API key")
			ctx.JSON(403, errors.Unauthorized)
			return
		}
		if !users.DoesUserHaveRole(userInfo.Roles, "developer") {
			for _, key := range body.Keys {
				apiKey := apikeys.GetAPIKeyForUser(key, userInfo.ID)
				if apiKey == nil {
					ctx.JSON(404, models.RouteErrorMessage{
						Title:       "API key error",
						Description: "We could not find an API key tied to our account with the given key",
					})
					return
				}
			}
		}
		err = apikeys.DeleteAPIKeys(userInfo.ID, body.Keys)
		if err != nil {
			SAR.Logger.Error("Failed to delete API key", "error", err, "keys", body.Keys)
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Server error",
				Description: "We encountered an error in the process of deleting your API key",
			})
		}
		ctx.JSON(200, gin.H{})
	}
}

func (SAR *SecuredAuthRouter) fetchAllAPIKeys() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxUserInfo, exists := ctx.Get("userInfo")
		userInfo := ctxUserInfo.(*authinfo.AuthInfo)
		if exists == false {
			SAR.Logger.Error("Failed to get userinfo in request to fetch all API keys")
			ctx.JSON(403, errors.Unauthorized)
			return
		}
		if !users.DoesUserHaveRole(userInfo.Roles, "developer") {
			ctx.JSON(403, errors.Unauthorized)
			return
		}
		keys := apikeys.GetAllAPIKeys()
		ctx.JSON(200, gin.H{
			"keys": keys,
		})
	}
}
