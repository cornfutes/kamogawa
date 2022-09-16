package handler

import (
	"html/template"
	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"
	"kamogawa/types/gcp/gcetypes"
	"strings"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func Functions(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "functions.html", gin.H{
				"Unauthorized": true,
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)

		responseSuccess, responseError := gcp.GCPListProjects(db, user, useCache)
		if responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
			core.HTMLWithGlobalState(c, "functions.html", gin.H{
				"MissingScopes": true,
			})
			return
		}

		var projectStrings []gcetypes.Project
		if user.Scope == nil || !user.Scope.Valid {
			projectStrings = []gcetypes.Project{}
		} else {
			projectStrings = responseSuccess.Projects
		}

		var htmlLines []string
		// Enumerate Projects for credentials
		for _, p := range projectStrings {
			responseSuccess, responseError := gcp.GCFListFunctions(user, p.ProjectId, nil)
			if responseError != nil && responseError.Error.Code > 0 {
				// Shortcircuit on first API call with missing scope to GCF.
				if responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
					core.HTMLWithGlobalState(c, "functions.html", gin.H{
						"MissingScopes": true,
					})
					return
				}

				if responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Cloud Functions API has not been used in project") {
					htmlLines = append(htmlLines, "<li>Cloud Functions API has not been abled in project</li>")
				} else {
					htmlLines = append(htmlLines, "<li>Error: \""+responseError.Error.Message+"\"</li>")
				}
			} else {
				if len(responseSuccess.Functions) == 0 {
					htmlLines = append(htmlLines, "<li>project"+p.ProjectId+" has no functions</li>")
				} else {
					for _, f := range responseSuccess.Functions {
						htmlLines = append(htmlLines, "<li>id: "+f.Name+"</li>")
					}
				}
			}
		}

		core.HTMLWithGlobalState(c, "functions.html", gin.H{
			"AssetLines": template.HTML(strings.Join(htmlLines[:], "")),
		})
	}
}
