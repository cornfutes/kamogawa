package handler

import (
	"kamogawa/core"

	"github.com/gin-gonic/gin"
)

// TODO: branch if logged in
func Blog1(c *gin.Context) {
	core.HTMLWithGlobalState(c, "blog_1.tmpl", gin.H{})
}
