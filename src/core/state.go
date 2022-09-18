package core

import (
	"fmt"
	"math/rand"
	"net/http"

	"kamogawa/config"
	"kamogawa/identity"
	"kamogawa/view"

	"github.com/gin-gonic/gin"
)

// TODO: have gin middleware that intercepts and adsd these binds.
func HTMLWithGlobalState(c *gin.Context, page string, obj map[string]interface{}) {
	value, _ := c.Get(identity.IdentityContextKey)
	obj["IsLoggedIn"] = value
	obj["EapUrl"] = config.EapUrl
	obj["ContactEmail"] = config.ContactEmail
	obj["BrandName"] = config.BrandName
	obj["SpaEnabled"] = config.SPAEnabled
	n := rand.Intn(2)
	fmt.Printf("Value %v\n", n)

	theme, err := c.Cookie(identity.CookieKeyTheme)
	if err != nil || theme == config.DefaultTheme || !view.IsValidTheme(theme) {
		c.HTML(http.StatusOK, page, obj)
	} else {
		c.HTML(http.StatusOK, page+theme, obj)
	}
}
