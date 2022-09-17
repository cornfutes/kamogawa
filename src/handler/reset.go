package handler

import (
	"kamogawa/core"

	"github.com/gin-gonic/gin"
)

func Reset(c *gin.Context) {
	core.HTMLWithGlobalState(c, "reset.tmpl", gin.H{
		"Resetted": c.Query("status") == "1",
	})
}
