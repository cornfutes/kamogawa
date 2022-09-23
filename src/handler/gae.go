package handler

import (
	"html/template"
	"strconv"
	"strings"

	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"
	"kamogawa/types/gcp/gcetypes"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func GAE(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		email, _ := c.Get(identity.IdentityContextKey)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "gae.tmpl", gin.H{
				"Email":        email,
				"Unauthorized": true,
				"APICallCount": "-1",
				"PageName":     "gcp_gae_overview",
				"Section":      "gcp",
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)

		apiCallCount := 1
		responseSuccess, responseError := gcp.GCPListProjects(db, user, useCache)
		if responseError != nil && responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
			core.HTMLWithGlobalState(c, "gae.tmpl", gin.H{
				"MissingScopes": true,
				"PageName":      "gcp_gae_overview",
				"Section":       "gcp",
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
			apiCallCount++
			responseSuccessService, responseErrorService := gcp.GAEListServices(db, user, p.ProjectId, useCache)
			htmlLines = append(htmlLines, "<li>"+p.ProjectId+" ( Project ) <ul>")

			if responseErrorService.Error.Code == 404 {
				htmlLines = append(htmlLines, "<li>")
				if strings.HasPrefix(responseErrorService.Error.Message, "App does not exist.") {
					htmlLines = append(htmlLines, "App Engine not initialized for this Project.")
				} else {
					htmlLines = append(htmlLines, "App engine state unknown for this Project.")
				}
				htmlLines = append(htmlLines, "</li>")

			} else {
				// Enumerate GAE Service(s) for Project
				for _, service := range responseSuccessService.Services {
					htmlLines = append(htmlLines, "<li>"+service.Id+" ( Service )<ul>")
					apiCallCount++
					responseSuccessVersion, responseErrorVersion := gcp.GAEListVersions(db, user, p.ProjectId, service.Id, useCache)
					if responseErrorVersion.Error.Code > 0 {
						htmlLines = append(htmlLines, "<li>")
						htmlLines = append(htmlLines, "Versions is an unknown state.")
						htmlLines = append(htmlLines, "</li>")
					} else {
						// Enumerate GAE Version(s) for Service
						for _, version := range responseSuccessVersion.Versions {
							htmlLines = append(htmlLines, "<li>"+version.Id+" ( Version ) <ul>")
							responseSuccessInstance, responseErrorInstance := gcp.GAEListInstances(db, user, p.ProjectId, service.Id, version.Id, useCache)
							if responseErrorInstance.Error.Code > 0 {
								htmlLines = append(htmlLines, "<li>Instances are in unknown state.</li>")
							} else {
								if len(responseSuccessInstance.Instances) == 0 {
									htmlLines = append(htmlLines, "<li>There are no Instances deployed for this version.</li>")
								}
								apiCallCount++
								// Enumerate GAE Version(s) for Version
								for _, instance := range responseSuccessInstance.Instances {
									htmlLines = append(htmlLines, "<li style=\"align-items-center;display:flex;\"><div style=\"white-space: nowrap; text-overflow: ellipsis; overflow: hidden; display: inline-block; width: 200px\">"+instance.Id+"</div>( Instance )</li>")
								}
							}
							htmlLines = append(htmlLines, "</ul></li>")
						}
					}
					htmlLines = append(htmlLines, "</ul></li>")
				}
			}

			htmlLines = append(htmlLines, "</ul>")
		}

		core.HTMLWithGlobalState(c, "gae.tmpl", gin.H{
			"Email":        email,
			"AssetLines":   template.HTML(strings.Join(htmlLines[:], "")),
			"APICallCount": strconv.Itoa(apiCallCount),
			"PageName":     "gcp_gae_overview",
			"Section":      "gcp",
		})
	}
}
