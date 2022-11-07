package handler

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"

	"gorm.io/gorm"
)

func GCE(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "gce.tmpl", gin.H{
				"NumCachedCalls": 0,
				"Unauthorized":   true,
				"PageName":       "gcp_gce_overview",
				"Section":        "gcp",
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
			core.HTMLWithGlobalState(c, "gce.tmpl", gin.H{
				"MissingScopes": true,
				"PageName":      "gcp_gce_overview",
				"Section":       "gcp",
			})
			return
		}

		var htmlLines []string
		var cachedCalls = 0
		for _, p := range projectDBs {
			responseSuccessAPI, _ := gcp.GCPListProjectAPIs(db, user, p, useCache)
			if len(responseSuccessAPI) == 0 || !responseSuccessAPI[0].IsEnabled("Compute Engine API") {
				cachedCalls++
				htmlLines = append(htmlLines, "<li>"+p.ProjectId+" ( Project ) <ul><li>Compute Engine API has not been enabled on project.</li></ul>")
				continue
			}

			responseSuccess, responseError := gcp.GCEListInstances(db, user, p.ProjectId, useCache)
			htmlLines = append(htmlLines, "<li>"+p.ProjectId+" ( Project ) <ul>")
			if responseError.Error.Code > 0 {
				// Shortcircuit if missing GCE scope.
				// TODO: refactor to do oonce utside loop
				if responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
					core.HTMLWithGlobalState(c, "gce.tmpl", gin.H{
						"MissingScopes": true,
						"PageName":      "gcp_gce_overview",
						"Section":       "gcp",
					})
					return
				} else {
					htmlLines = append(htmlLines, "<li>Unknown error with code: "+strconv.Itoa(responseError.Error.Code)+"</li>")
				}
			} else {
				if responseSuccess == nil || len(responseSuccess.Zones) == 0 {
					htmlLines = append(htmlLines, "<li>There are no instances in this project.</li>")
				} else {
					for _, zone := range responseSuccess.Zones {
						htmlLines = append(htmlLines, "<li>"+zone.Zone+" ( Zone ) <ul>")
						for _, instance := range zone.Instances {
							htmlLines = append(htmlLines, "<li>"+instance.Name+" ( Instance )</li>")
						}
						htmlLines = append(htmlLines, "</ul></li>")
					}
				}
			}
			htmlLines = append(htmlLines, "</ul>")
		}

		duration := time.Since(start)
		core.HTMLWithGlobalState(c, "gce.tmpl", gin.H{
			"Duration":       duration,
			"NumCachedCalls": cachedCalls,
			"AssetLines":     template.HTML(strings.Join(htmlLines[:], "")),
			"PageName":       "gcp_gce_overview",
			"Section":        "gcp",
		})
	}
}
