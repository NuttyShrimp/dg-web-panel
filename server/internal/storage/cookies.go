package storage

import (
	"degrens/panel/internal/config"
	"degrens/panel/internal/db"
	"degrens/panel/lib/log"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

var cookieOptions CookieOptions
var logger log.Logger

func InitCookieStore(config *config.Config, logger2 *log.Logger) {
	cookieOptions = CookieOptions{
		MaxAge: 86400 * 30,
		Domain: config.Server.Host,
		Codecs: securecookie.CodecsFromPairs([]byte(config.Server.SessionSecret)),
	}
	logger = (*logger2).With("module", "cookies")
}

type CookieOptions struct {
	MaxAge int
	Domain string
	Codecs []securecookie.Codec
}

// This will add a cokkie without hiding the actual value in redis
func AddPublicCookie(c *gin.Context, key string, value interface{}) bool {
	// Set cookie
	return setCookie(c, key, value)
}

func AddHiddenCookie(c *gin.Context, key string, value interface{}) bool {
	// Create an uuid that we will use to get/store the value in redis
	uuid, err := db.Redis.GenerateUUID()
	if err != nil {
		logger.Error("Could not generate UUID for hidden cookie", "error", err)
		return false
	}
	// Store the value in redis
	err = db.Redis.Set(uuid, value)
	if err != nil {
		logger.Error("Could not store value in redis", "error", err)
		return false
	}
	isSucces := setCookie(c, key, uuid)
	if !isSucces {
		err := db.Redis.Remove(uuid)
		if err != nil {
			logger.Error("Could not remove value from redis", "error", err)
		}
		return false
	}
	return true
}

func GetPublicCookie(c *gin.Context, key string, dst any) error {
	encodedCookie, err := c.Cookie(key)
	if err != nil {
		return errors.New("Cookie not found")
	}
	err = securecookie.DecodeMulti(key, encodedCookie, dst, cookieOptions.Codecs...)
	if err != nil {
		logger.Error("Failed to decode cookie ", err.Error())
		return errors.New("Failed to parse cookie")
	}
	return nil
}

func GetHiddenCookie(c *gin.Context, key string, dst any) error {
	var uuid string
	err := GetPublicCookie(c, key, &uuid)
	if err != nil {
		return err
	}
	// Check if pointer is nil
	if &uuid == nil {
		return nil
	}
	return db.Redis.Get(uuid, dst)
}

func RemoveCookie(c *gin.Context, key string) {
	c.SetCookie(key, "", -1, "/", cookieOptions.Domain, false, false)
}

// Helpers
func setCookie(c *gin.Context, name string, value interface{}) bool {
	// encode cookie
	encoded, err := securecookie.EncodeMulti(name, value, cookieOptions.Codecs...)
	if err != nil {
		logger.Error("Failed to encode cookie with key: "+name, "error", err.Error())
		return false
	}
	// Maybe change to strict mode but should add param then to manage this for eg. state
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, encoded, cookieOptions.MaxAge, "/", cookieOptions.Domain, false, false)
	return true
}
