package identity

import (
	"bytes"
	"database/sql"
	"fmt"
	"kamogawa/config"
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

type ErrorRefreshToken struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func GCPRefresh(c *gin.Context, db *gorm.DB) (*ResponseRefreshToken, *ErrorRefreshToken) {
	var email, _ = c.Get(IdentityContextKey)

	fmt.Printf("Revoking for email: %v\n", email)
	var user types.User
	db.First(&user, "email = ?", email)

	postBody, err := json.Marshal(map[string]string{
		"client_id":     config.GCPClientId,
		"client_secret": config.GCPClientSecret,
		"refresh_token": "1//06Ch0xWOX_kIWCgYIARAAGAYSNgF-L9IrRcMGwxK5-bEVrxNgUrnjS42vSWpgzL4JHZ4mE5WVxkNFyH_sZnhYh1ELaegt4B8z7Q",
		"grant_type":    "refresh_token",
	})
	if err != nil {
		log.Fatalf("An error occured while refreshing %v", err)
	}
	reqBody := bytes.NewBuffer(postBody)
	fmt.Printf("PostBody '%v' \n", reqBody)
	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/json", reqBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var responseSuccess ResponseRefreshToken
		if err := json.NewDecoder(resp.Body).Decode(&responseSuccess); err != nil {
			panic(err)
		}
		return &responseSuccess, nil
	} else {
		var responseError ErrorRefreshToken
		if err := json.NewDecoder(resp.Body).Decode(&responseError); err != nil {
			panic(err)
		}
		return nil, &responseError
	}
}

func CheckDBAndRefreshToken(c *gin.Context, user types.User, db *gorm.DB) {
	if time.Since(user.UpdatedAt).Seconds() > 3300 || len(user.AccessToken.String) == 0 {
		responseRefreshToken, errorRefreshToken := GCPRefresh(c, db)
		if errorRefreshToken != nil {
			if errorRefreshToken.Error == "invalid_grant" && errorRefreshToken.ErrorDescription == "Bad Request" {
				panic("Invalid refresh token")
			}
		}
		user.AccessToken = &sql.NullString{String: responseRefreshToken.AccessToken, Valid: true}
		db.Save(&user)
	}
}
