package handler

import (
	"fmt"
	"html/template"
	"kamogawa/cache/gcecache"
	"kamogawa/config"
	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"
	"kamogawa/types/gcp/gcetypes"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func GCE(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "gce.html", gin.H{
				"NumCachedCalls": 0,
				"Unauthorized":   true,
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)
		if user.Scope == nil || !user.Scope.Valid {
			panic("Missing scope")
		}
		fmt.Printf("User %v\n", user)

		var start = time.Now()
		var responseSuccess *gcetypes.ListProjectResponse
		var projectDBs []gcetypes.ProjectDB = nil
		if config.CacheEnabled && useCache {
			projectDBs = gcecache.ReadProjectsCache2(db, user)
		}
		if projectDBs == nil {
			var responseError *gcetypes.ErrorResponse
			responseSuccess, responseError = gcp.GCPListProjectsMain(db, user)
			if responseError != nil && responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
				core.HTMLWithGlobalState(c, "gce.html", gin.H{
					"MissingScopes": true,
				})
				return
			}
			projectDBs = make([]gcetypes.ProjectDB, 0, len(responseSuccess.Projects))
			for _, p := range responseSuccess.Projects {
				projectDBs = append(projectDBs, gcetypes.ProjectToProjectDB(user.Email, &p, 0))
			}
			fmt.Printf("len %v\n", projectDBs)
		}

		var htmlLines []string
		var cachedCalls = 0
		for i, p := range projectDBs {
			if projectDBs[i].HasGCEEnabled == -1 {
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
					core.HTMLWithGlobalState(c, "gce.html", gin.H{
						"MissingScopes": true,
					})
					return
				}

				if responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Compute Engine API has not been used in project") {
					htmlLines = append(htmlLines, "<li>Compute Engine API has not been enabled on project.</li>")
					projectDBs[i].HasGCEEnabled = -1
				} else {
					htmlLines = append(htmlLines, "<li>Unknown error with code: "+strconv.Itoa(responseError.Error.Code)+"</li>")
				}
			} else {
				if len(responseSuccess.Zones) == 0 {
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
		core.HTMLWithGlobalState(c, "gce.html", gin.H{
			"Duration":       duration,
			"NumCachedCalls": cachedCalls,
			"AssetLines":     template.HTML(strings.Join(htmlLines[:], "")),
		})

		if config.CacheEnabled {
			fmt.Printf("writing cache %v\n", len(projectDBs))
			gcecache.WriteProjectsCache2(db, user, projectDBs)
		}
	}
}
