package handler

import (
	"kamogawa/core"
	"kamogawa/identity"

	"github.com/gin-gonic/gin"
)

func Account(c *gin.Context) {
	email, _ := c.Get(identity.IdentityContextkey)

	core.HTMLWithGlobalState(c, "account.html", gin.H{
		"Email": email,
	})
}
