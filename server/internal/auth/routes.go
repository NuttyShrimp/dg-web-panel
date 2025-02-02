package auth

import (
	"context"
	"degrens/panel/internal/auth/apikeys"
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/cfxauth"
	"degrens/panel/internal/auth/cfxtoken"
	"degrens/panel/internal/auth/discord"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/storage"
	"degrens/panel/internal/users"
	"degrens/panel/lib/errors"
	"degrens/panel/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func NewAuthRouter(rg *gin.RouterGroup) {
	router := &AuthRouter{
		routes.Router{
			RouterGroup: rg.Group("/auth"),
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
	if userInfo.ID == 0 {
		// No session cookie
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No session cookie",
		})
		return
	}
	switch userInfo.AuthMethod {
	case authinfo.Discord:
		isExpired := discord.IsTokenExpired(userInfo.ID)
		c.JSON(http.StatusOK, gin.H{
			"isExpired": isExpired,
		})
		if !isExpired {
			discord.RefreshToken(userInfo.ID)
		}
		return
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No refresh method for this auth method",
		})
		return
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
		case "cfxtoken":
			tokenHeader := c.GetHeader("Authorization")
			if tokenHeader == "" {
				// No header given, yeet req
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to get a valid token from the request",
				})
				return
			}

			authTokens := strings.Split(tokenHeader, " ")

			if len(authTokens) != 2 || authTokens[0] != "Bearer" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to get a valid token from the request",
				})
				return
			}

			err := cfxauth.AuthorizeToken(c, authTokens[1])
			if err != nil {
				logrus.WithError(err).WithField("token", tokenHeader).Error("Failed to authorize a cfx token")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Failed to use cfx token to create session",
				})
				return
			}
			logrus.Debug("Successfully authorized cfx token")
			c.JSON(200, gin.H{})
		}
	}
}

func (AR *AuthRouter) logoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfoPtr, exists := c.Get("userInfo")
		userInfo := userInfoPtr.(*authinfo.AuthInfo)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get your details",
			})
			return
		}
		switch userInfo.AuthMethod {
		case authinfo.Discord:
			discord.RemoveUserTokens(userInfo.ID)
		case authinfo.CFXToken:
			// Send request to Cfx API to invalidate login token
			err := cfxtoken.RemoveToken(userInfo.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to get your details",
				})
			}
		default:
			logrus.WithField("info", fmt.Sprintf("%+v", userInfo)).Error("Logout procedure failed to determine authentication type")
		}
		storage.RemoveCookie(c, "userInfo")
		c.String(http.StatusOK, "")
	}
}

func (AR *AuthRouter) discordCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the state from the cookies
		var state string
		if err := storage.GetPublicCookie(c, "state", &state); state == "" || err != nil {
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
			logrus.WithError(err).Error("Error while exchanging code for token")
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
		new_authinfo := authinfo.GetAuthInfoFromUser(user)
		err = new_authinfo.Assign(c)
		if err != nil {
			logrus.Error(err)
			c.Redirect(http.StatusTemporaryRedirect, "/errors/500")
		}

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func NewSecuredAuthRouter(rg *gin.RouterGroup) {
	router := &SecuredAuthRouter{
		routes.Router{
			RouterGroup: rg.Group("/auth"),
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
			logrus.WithField("target", target).Error("Failed to parse role in /auth/role request")
			ctx.JSON(400, models.RouteErrorMessage{
				Title:       "Request error",
				Description: "Encountered an error while trying to check access to a secure page",
			})
			return
		}
		hasAccess, err := users.HasRoleAccess(ctx, target)
		if err != nil {
			logrus.WithError(err).Error("Failed to check role in /auth/role request")
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
		if !exists {
			logrus.Error("Failed to get userinfo while fetching API keys")
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
			logrus.WithError(err).Error("Failed to parse body in APIKey creation request")
			ctx.JSON(500, errors.BodyParsingFailed)
			return
		}
		ctxUserInfo, exists := ctx.Get("userInfo")
		userInfo := ctxUserInfo.(*authinfo.AuthInfo)
		if !exists {
			logrus.Error("Failed to get userinfo in request trying to make an API key")
			ctx.JSON(403, errors.Unauthorized)
			return
		}
		var key string
		key, err = apikeys.CreateAPIKey(userInfo.ID, body.Comment, time.Duration(body.Duration)*time.Minute)
		if err != nil {
			logrus.WithField("body", body).WithError(err).Error("Failed to create API key")
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Server error",
				Description: "We encountered an error in the process of creating your API key",
			})
		}
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
			logrus.WithError(err).Error("Failed to parse body in APIKey deletion request")
			ctx.JSON(500, errors.BodyParsingFailed)
			return
		}
		ctxUserInfo, exists := ctx.Get("userInfo")
		userInfo := ctxUserInfo.(*authinfo.AuthInfo)
		if !exists {
			logrus.Error("Failed to get userinfo in request trying to make an API key")
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
		apikeys.DeleteAPIKeys(userInfo.ID, body.Keys)
		ctx.JSON(200, gin.H{})
	}
}

func (SAR *SecuredAuthRouter) fetchAllAPIKeys() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxUserInfo, exists := ctx.Get("userInfo")
		userInfo := ctxUserInfo.(*authinfo.AuthInfo)
		if !exists {
			logrus.Error("Failed to get userinfo in request to fetch all API keys")
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
