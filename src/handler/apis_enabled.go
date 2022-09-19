package handler

import (
	"fmt"
	"html/template"
	"kamogawa/types/gcp/coretypes"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"

	"gorm.io/gorm"
)

func APIsEnabled(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "apis_enabled.tmpl", gin.H{
				"NumCachedCalls": 0,
				"Unauthorized":   true,
				"PageName":       "gcp_apis_enabled",
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)
		if user.Scope == nil || !user.Scope.Valid {
			panic("Missing scope")
		}
		fmt.Printf("User %v\n", user)

		var start = time.Now()

		projectDBs, responseError := gcp.GCPListProjects(db, user, useCache)
		if responseError != nil && responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
			core.HTMLWithGlobalState(c, "apis_enabled.tmpl", gin.H{
				"MissingScopes": true,
			})
			return
		}

		var htmlLines []string
		var cachedCalls = 0
		for _, p := range projectDBs {
			cachedCalls++
			gcpProjectAPIs, _ := gcp.GCPListProjectAPIs(db, user, p, useCache)
			resp := coretypes.JSONToGCPServiceUsageServicesListResponse(gcpProjectAPIs[0].API)

			htmlLines = append(htmlLines, "<li>"+p.ProjectId+" ( Project ) <ul>")
			for _, s := range resp.Services {
				htmlLines = append(htmlLines, "<li>"+s.Config.Title+" ("+s.Config.Name+")</li>")
			}
			htmlLines = append(htmlLines, "</ul>")
		}

		duration := time.Since(start)
		core.HTMLWithGlobalState(c, "apis_enabled.tmpl", gin.H{
			"Duration":       duration,
			"NumCachedCalls": cachedCalls,
			"AssetLines":     template.HTML(strings.Join(htmlLines[:], "")),
			"PageName":       "gcp_apis_enabled",
		})
	}
}
