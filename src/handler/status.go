package handler

import (
	"kamogawa/core"

	"github.com/gin-gonic/gin"
)

func Status(c *gin.Context) {
	core.HTMLWithGlobalState(c, "status.tmpl", gin.H{})
}
