package handler

import (
	"kamogawa/core"

	"github.com/gin-gonic/gin"
)

func Reset(c *gin.Context) {
	core.HTMLWithGlobalState(c, "reset.html", gin.H{
		"Resetted": c.Query("status") == "1",
	})
}
