package handler

import (
	"kamogawa/config"
	"kamogawa/core"
	"kamogawa/identity"
	"kamogawa/types"
	"net/url"
	"strings"

	"sort"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var clientId = "38029384895-t1chm9en5703tjmkqfqpa1ki6020p8n6.apps.googleusercontent.com"
var redirectUri = config.RedirectUri
var googleAuthUrl = getUrlWithScopes(GCPScopes)

func getUrlWithScopes(scopes string) string {
	return "https://accounts.google.com/o/oauth2/v2/auth?" +
		"client_id=" + clientId + "&" +
		"redirect_uri=" + url.QueryEscape(redirectUri) + "&" +
		"response_type=code&" +
		"scope=" + url.QueryEscape(scopes) + "&" +
		"prompt=consent&" +
		"access_type=offline&" +
		"include_granted_scopes=true"
}

func Authorization(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var email, exists = c.Get(identity.IdentityContextKey)
		if !exists {
			panic("Unexpected")
		}
		var user types.User
		err := db.First(&user, "email = ?", email).Error
		if err != nil {
			panic("Invalid DB state")
		}

		var gmail string
		if user.Gmail == nil {
			gmail = ""
		} else {
			gmail = user.Gmail.String
		}
		var scopes []string
		if user.Scope == nil || !user.Scope.Valid {
			scopes = []string{}
		} else {
			scopes = strings.Fields(user.Scope.String)
			sort.Strings(scopes)
		}

		missingScopeToUrl := make(map[string]ScopeMetadata)
		for _, metadata := range AllScopes {
			found := false
			for _, hasScope := range scopes {
				if metadata.Scope == hasScope {
					found = true
				}
			}
			if !found {
				missingScopeToUrl[getUrlWithScopes(metadata.Scope)] = metadata
			}
		}
		core.HTMLWithGlobalState(c, "authorization.html", gin.H{
			"GmailExists":         user.Gmail != nil && user.Gmail.Valid,
			"Gmail":               gmail,
			"GCPDelegatedAuthUrl": googleAuthUrl,
			"Scopes":              scopes,
			"HasMissingScopes":    len(missingScopeToUrl) > 0,
			"MissingScopes":       missingScopeToUrl,
		})
	}
}
