package identity

import (
	"database/sql"
	"fmt"
	"kamogawa/types"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RevokeGCP(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var email, _ = c.Get(IdentityContextKey)

		fmt.Printf("Revoking for email: %v\n", email)
		var user types.User
		db.First(&user, "email = ?", email)
		fmt.Printf("Here %v\n", user.AccessToken.Valid)
		// TODO: Does any accesstoken work? What about refreshtoken?
		resp, err := http.Post("https://oauth2.googleapis.com/revoke?token="+user.AccessToken.String,
			"application/json", nil)
		if err != nil {
			c.Redirect(http.StatusFound, "/authorization?state=bad")
		}
		fmt.Printf("Successfully revoked token %v\n", resp)

		user.Gmail = &sql.NullString{}
		user.Scope = &sql.NullString{}
		user.AccessToken = &sql.NullString{}
		user.RefreshToken = &sql.NullString{}
		db.Save(&user)

		c.Redirect(http.StatusFound, "/authorization?state=revoked")
	}
}
