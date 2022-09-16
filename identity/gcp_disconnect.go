package identity

import (
	"database/sql"
	"fmt"
	"kamogawa/types"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DisconnectGCP(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var email, _ = c.Get(IdentityContextKey)

		var user types.User
		db.First(&user, "email = ?", email)
		fmt.Printf("Disconnect gmail '%v' from '%v'\n", user.Gmail, email)

		user.Gmail = &sql.NullString{}
		user.Scope = &sql.NullString{}
		user.AccessToken = &sql.NullString{}
		user.RefreshToken = &sql.NullString{}
		db.Save(&user)

		c.Redirect(http.StatusFound, "/authorization")
	}
}
