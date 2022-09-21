package identity

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"kamogawa/config"
	"kamogawa/types"

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
	email, _ := c.Get(IdentityContextKey)

	fmt.Printf("Refresh acess token for: %v\n", email)
	var user types.User
	db.First(&user, "email = ?", email)

	postBody, err := json.Marshal(map[string]string{
		"client_id":     config.GCPClientId,
		"client_secret": config.GCPClientSecret,
		"refresh_token": user.RefreshToken.String,
		"grant_type":    "refresh_token",
	})
	if err != nil {
		log.Fatalf("An error occured while refreshing %v", err)
	}
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/json", reqBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("Successfully refreshed access token")
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
	// TODO: What if we update the user for another reason, resetting the clock?
	if time.Since(user.UpdatedAt).Seconds() > 3300 || len(user.AccessToken.String) == 0 {
		responseRefreshToken, errorRefreshToken := GCPRefresh(c, db)
		if errorRefreshToken != nil {
			// The refresh token did not look like a refresh token.
			if errorRefreshToken.Error == "invalid_grant" && errorRefreshToken.ErrorDescription == "Bad Request" {
				panic(errorRefreshToken)
			}
			// Was a legit refresh token
			if errorRefreshToken.Error == "invalid_grant" && errorRefreshToken.ErrorDescription == "Token has been expired or revoked." {
				panic(errorRefreshToken)
			}
			// Empty string for refresh token.
			if errorRefreshToken.Error == "invalid_request" && errorRefreshToken.ErrorDescription == "Missing required parameter: refresh_token" {
				panic(errorRefreshToken)
			}
			panic("Error: '" + errorRefreshToken.Error + "', description: '" + errorRefreshToken.ErrorDescription + "'")
		}
		fmt.Printf("Saving new access token\n")
		user.AccessToken = &sql.NullString{String: responseRefreshToken.AccessToken, Valid: true}
		db.Save(&user)
	}
}
