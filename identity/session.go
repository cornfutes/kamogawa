package identity

import (
	"kamogawa/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CheckSessionForUser(c *gin.Context, db *gorm.DB) types.User {
	var email, exists = c.Get(IdentityContextkey)
	if !exists {
		panic("Unexpected")
	}
	var user types.User
	err := db.First(&user, "email = ?", email).Error
	if err != nil {
		panic("Unvalid DB state")
	}
	return user
}
