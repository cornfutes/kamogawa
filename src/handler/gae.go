package handler

import (
	"html/template"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"

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
		if user.Scope == nil || !user.Scope.Valid {
			panic("Missing scope")
		}

		apiCallCount := 1
		projectDBs, responseError := gcp.GCPListProjects(db, user, useCache)
		if responseError != nil && responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
			core.HTMLWithGlobalState(c, "gae.tmpl", gin.H{
				"MissingScopes": true,
				"PageName":      "gcp_gae_overview",
				"Section":       "gcp",
			})
			return
		}

		var htmlLines []string
		// Enumerate Projects for credentials
		for _, p := range projectDBs {
			apiCallCount++
			htmlLines = append(htmlLines, "<li>"+p.ProjectId+" ( Project ) <ul>")

			responseSuccessAPI, _ := gcp.GCPListProjectAPIs(db, user, p, useCache)
			if !responseSuccessAPI[0].IsGAEEnabled {
				htmlLines = append(htmlLines, "<li>App Engine not initialized for this Project.</li>")
				htmlLines = append(htmlLines, "</ul>")
				continue
			}

			responseSuccessService, responseErrorService := gcp.GAEListServices(db, user, p.ProjectId, useCache)
			if responseErrorService.Error.Code == 404 {
				htmlLines = append(htmlLines, "<li>App engine state unknown for this Project.</li>")
				htmlLines = append(htmlLines, "</ul>")
				continue
			}
			if responseSuccessService == nil {
				htmlLines = append(htmlLines, "</ul>")
				continue
			}

			// Enumerate GAE Service(s) for Project
			for _, service := range responseSuccessService.Services {
				htmlLines = append(htmlLines, "<li>"+service.Id+" ( Service )<ul>")
				apiCallCount++
				responseSuccessVersion, responseErrorVersion := gcp.GAEListVersions(db, user, p.ProjectId, service.Id, useCache)
				if responseErrorVersion.Error.Code > 0 {
					htmlLines = append(htmlLines, "<li>")
					htmlLines = append(htmlLines, "Versions is an unknown state.")
					htmlLines = append(htmlLines, "</li>")
					htmlLines = append(htmlLines, "</ul></li>")
					continue
				}
				if responseSuccessVersion == nil {
					htmlLines = append(htmlLines, "</ul></li>")
					continue
				}

				// Enumerate GAE Version(s) for Service
				for _, version := range responseSuccessVersion.Versions {
					htmlLines = append(htmlLines, "<li>"+version.Id+" ( Version ) <ul>")
					responseSuccessInstance, responseErrorInstance := gcp.GAEListInstances(db, user, p.ProjectId, service.Id, version.Id, useCache)
					if responseErrorInstance.Error.Code > 0 {
						htmlLines = append(htmlLines, "<li>Instances are in unknown state.</li>")
						htmlLines = append(htmlLines, "</ul></li>")
						continue
					}
					if responseSuccessInstance == nil || len(responseSuccessInstance.Instances) == 0 {
						htmlLines = append(htmlLines, "<li>There are no Instances deployed for this version.</li>")
						htmlLines = append(htmlLines, "</ul></li>")
						continue
					}
					apiCallCount++
					// Enumerate GAE Version(s) for Version
					for _, instance := range responseSuccessInstance.Instances {
						htmlLines = append(htmlLines, "<li style=\"align-items-center;display:flex;\"><div style=\"white-space: nowrap; text-overflow: ellipsis; overflow: hidden; display: inline-block; width: 200px\">"+instance.Id+"</div>( Instance )</li>")
					}

					htmlLines = append(htmlLines, "</ul></li>")
				}

				htmlLines = append(htmlLines, "</ul></li>")
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
