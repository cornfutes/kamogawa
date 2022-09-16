package identity

import (
	"kamogawa/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleLogout(c *gin.Context) {
	// The way to delete and nullify a cookie is to set an expiration in the past.
	// -1 indicates -1 second from now.
	// TODO: evaluate if there may be bug due to timezone skew
	const maxAge = -1
	c.SetCookie(SessionCookieKey, "", maxAge, "/", config.Host, false, true)

	c.Redirect(http.StatusFound, "/loggedout")
}
