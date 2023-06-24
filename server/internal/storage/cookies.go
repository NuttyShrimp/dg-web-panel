package storage

import (
	"degrens/panel/internal/config"
	"degrens/panel/internal/db"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/sirupsen/logrus"
)

var cookieOptions CookieOptions
var logger *logrus.Entry

func InitCookieStore(conf *config.Config) {
	cookieOptions = CookieOptions{
		MaxAge: 86400 * 30,
		Domain: conf.Server.GetCookieHost(),
		Codecs: securecookie.CodecsFromPairs([]byte(conf.Server.SessionSecret)),
	}
	logger = logrus.WithField("module", "cookies")
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
		logger.WithError(err).Error("Could not generate UUID for hidden cookie")
		return false
	}
	// Store the value in redis
	err = db.Redis.Set(uuid, value)
	if err != nil {
		logger.WithError(err).Error("Could not store value in redis")
		return false
	}
	isSucces := setCookie(c, key, uuid)
	if !isSucces {
		err := db.Redis.Remove(uuid)
		if err != nil {
			logger.WithError(err).Error("Could not remove value from redis")
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
		logger.WithError(err).Error("Failed to decode cookie")
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
	// uuid is empty
	if uuid == "" {
		return nil
	}
	return db.Redis.Get(uuid, dst)
}

func RemoveCookie(c *gin.Context, key string) {
	c.SetCookie(key, "", -1, "/", cookieOptions.Domain, true, false)
}

// Helpers
func setCookie(c *gin.Context, name string, value interface{}) bool {
	// encode cookie
	encoded, err := securecookie.EncodeMulti(name, value, cookieOptions.Codecs...)
	if err != nil {
		logger.WithError(err).Error("Failed to encode cookie with key: " + name)
		return false
	}
	// Maybe change to strict mode but should add param then to manage this for eg. state
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(name, encoded, cookieOptions.MaxAge, "/", cookieOptions.Domain, true, false)
	return true
}
