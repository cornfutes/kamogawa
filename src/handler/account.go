package handler

import (
	"kamogawa/core"
	"kamogawa/identity"

	"github.com/gin-gonic/gin"
)

func Account(c *gin.Context) {
	email, _ := c.Get(identity.IdentityContextKey)

	core.HTMLWithGlobalState(c, "account.tmpl", gin.H{
		"Email": email,
	})
}
