package handler

import (
	"html/template"
	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"
	"kamogawa/types"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func GCE(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "gce.html", gin.H{
				"Unauthorized": true,
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)

		responseSuccess, responseError := gcp.GCPListProjects(db, user, useCache)
		if responseError != nil && responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
			core.HTMLWithGlobalState(c, "gce.html", gin.H{
				"MissingScopes": true,
			})
			return
		}

		var projectStrings []types.Project
		if user.Scope == nil || !user.Scope.Valid {
			projectStrings = []types.Project{}
		} else {
			projectStrings = responseSuccess.Projects
		}

		var htmlLines []string
		// Enumerate Projects for credentials
		for _, p := range projectStrings {
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

		core.HTMLWithGlobalState(c, "gce.html", gin.H{
			"AssetLines": template.HTML(strings.Join(htmlLines[:], "")),
		})
	}
}
