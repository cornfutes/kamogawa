package handler

import (
	"kamogawa/core"

	"github.com/gin-gonic/gin"
)

func Export(c *gin.Context) {
	core.HTMLWithGlobalState(c, "export.html", gin.H{})
}
