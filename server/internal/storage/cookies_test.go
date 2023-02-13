package storage_test

import (
	"degrens/panel/internal/storage"
	"degrens/panel/lib/tests"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
)

var Codecs []securecookie.Codec
var env *tests.Env

func TestMain(m *testing.M) {
	env = tests.LoadBareEnv()
	storage.InitCookieStore(env.Config, env.Logger)
	storage.InitStateTokenStorage()
	Codecs = securecookie.CodecsFromPairs([]byte(env.Config.Server.SessionSecret))
	m.Run()
}

func TestAddPublicCookie(t *testing.T) {
	testContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	wasSuccessfull := storage.AddPublicCookie(testContext, "test", true)
	if !wasSuccessfull {
		t.Error("A cookie should not fail to set on a valid context")
		return
	}
	header := testContext.Writer.Header().Get("Set-Cookie")
	securedValue, _ := securecookie.EncodeMulti("test", true, Codecs...)
	assert.Equal(t, fmt.Sprint("test=", url.QueryEscape(securedValue), "; Path=/; Domain=localhost; Max-Age=2592000; SameSite=Lax"), header)
}

func TestInvalidPublicCookie(t *testing.T) {
	testContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	wasSuccessfull := storage.AddPublicCookie(testContext, "test", nil)
	if wasSuccessfull {
		t.Error("Setting a cookie should fail when setting it to a invalid value")
		return
	}
}
