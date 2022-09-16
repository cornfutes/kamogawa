package identity

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleReset(c *gin.Context) {
	c.Redirect(http.StatusFound, "/reset?status=1")
}
