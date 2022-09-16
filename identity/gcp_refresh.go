package identity

import (
	"bytes"
	"database/sql"
	"fmt"
	"kamogawa/types"
	"log"
	"time"

	"net/http"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResponseRefreshToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func GCPRefresh(c *gin.Context, db *gorm.DB) ResponseRefreshToken {
	var email, _ = c.Get(IdentityContextkey)

	fmt.Printf("Revoking for email: %v\n", email)
	var user types.User
	db.First(&user, "email = ?", email)

	postBody, err := json.Marshal(map[string]string{
		"client_id":     "<retroactively_redacted>",
		"client_secret": "<retroactively_redacted>",
		"refresh_token": user.RefreshToken.String,
		"grant_type":    "refresh_token",
	})
	if err != nil {
		log.Fatalf("An Error occured while refreshing %v", err)
	}
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/json", reqBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var responseRefreshToken ResponseRefreshToken
	err = json.NewDecoder(resp.Body).Decode(&responseRefreshToken)
	if err != nil {
		panic(err)
	}
	log.Println(responseRefreshToken)

	return responseRefreshToken
}

func CheckDBAndRefreshToken(c *gin.Context, user types.User, db *gorm.DB) {
	if time.Now().Sub(user.UpdatedAt).Seconds() > 3300 {
		respRefreshtoken := GCPRefresh(c, db)
		user.AccessToken = &sql.NullString{String: respRefreshtoken.AccessToken, Valid: true}
		db.Save(&user)
	}
}
