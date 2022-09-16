package handler

import (
	"kamogawa/core"

	"github.com/gin-gonic/gin"
)

// TODO: branch if logged in
func Demo(c *gin.Context) {
	core.HTMLWithGlobalState(c, "demo.html", gin.H{})
}
