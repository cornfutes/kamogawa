package core

import (
	"net/http"

	"kamogawa/config"
	"kamogawa/identity"

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
	c.HTML(http.StatusOK, page, obj)
}
