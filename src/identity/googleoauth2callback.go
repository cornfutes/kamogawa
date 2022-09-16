package identity

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"kamogawa/config"
	"kamogawa/types"
	"log"
	"net/http"

	"github.com/ederinbay/GoogleIdTokenVerifier"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccessTokenResponse struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    float64 `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	IdToken      string  `json:"id_token"`
	Scope        string  `json:"scope"`
	TokenType    string  `json:"token_type"`
}

type UserInfoResponse struct {
	Email string `json:"email"`
}

var Tokens = make(map[string]string)

func GoogleOAuth2Callback(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Printf("o  \n")

		code := c.Query("code")

		// TODO: fix oauth in prod.
		postBody, err := json.Marshal(map[string]string{
			"code":          code,
			"client_id":     config.GCPClientId,
			"client_secret": config.GCPClientSecret,
			"redirect_uri":  config.RedirectUri,
			"grant_type":    "authorization_code",
		})
		if err != nil {
			log.Fatalf("An Error Occured exchanging %v", err)
		}

		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post("https://oauth2.googleapis.com/token", "application/json", responseBody)
		//Handle Error
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()

		var token AccessTokenResponse
		err = json.NewDecoder(resp.Body).Decode(&token)
		if err != nil {
			panic(err)
		}
		Tokens["access_token"] = token.AccessToken
		Tokens["refresh_token"] = token.RefreshToken
		fmt.Printf("RefreshToken: %v \n", token.RefreshToken)

		var ti GoogleIdTokenVerifier.TokenInfo = *GoogleIdTokenVerifier.Verify(token.IdToken, config.GCPClientId)
		gmail := ti.Email

		var email, exists = c.Get(IdentityContextKey)
		if !exists {
			panic("Unexpected exrror lookup context")
		}
		var user types.User
		err = db.First(&user, "email = ?", email).Error
		if err != nil {
			panic("Invalid DB state")
		}
		db.Model(&user).Updates(types.User{
			Gmail:        &sql.NullString{String: gmail, Valid: true},
			Scope:        &sql.NullString{String: token.Scope, Valid: true},
			AccessToken:  &sql.NullString{String: token.AccessToken, Valid: true},
			RefreshToken: &sql.NullString{String: token.RefreshToken, Valid: true},
		})
		c.Redirect(http.StatusFound, "/authorization")
	}
}
