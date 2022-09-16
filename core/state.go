package core

import (
	"kamogawa/config"
	"kamogawa/identity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HTMLWithGlobalState(c *gin.Context, page string, obj map[string]interface{}) {
	value, _ := c.Get(identity.IdentityContextKey)
	obj["IsLoggedIn"] = value
	obj["EapUrl"] = config.EapUrl
	obj["ContactEmail"] = config.ContactEmail
	obj["BrandName"] = config.BrandName
	c.HTML(http.StatusOK, page, obj)
}
