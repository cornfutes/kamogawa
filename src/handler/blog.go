package handler

import (
	"kamogawa/core"

	"github.com/gin-gonic/gin"
)

// TODO: branch if logged in
func Blog(c *gin.Context) {
	core.HTMLWithGlobalState(c, "blog.tmpl", gin.H{})
}
