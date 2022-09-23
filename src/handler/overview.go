package handler

import (
	"strconv"
	"strings"
	"time"

	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"
	"kamogawa/types/gcp/gcetypes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Overview(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "overview.tmpl", gin.H{
				"Unauthorized": true,
				"PageName":     "gcp_overview",
				"Section":      "gae",
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)

		responseSuccess, responseError := gcp.GCPListProjects(db, user, useCache)
		if responseError != nil && responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
			core.HTMLWithGlobalState(c, "overview.tmpl", gin.H{
				"MissingScopes": true,
				"PageName":      "gcp_overview",
				"Section":       "gae",
			})
			return
		}

		var projectStrings []gcetypes.Project
		if user.Scope == nil || !user.Scope.Valid {
			projectStrings = []gcetypes.Project{}
		} else {
			projectStrings = responseSuccess.Projects
			for i, project := range projectStrings {
				if project.ProjectId == project.Name {
					projectStrings[i].ProjectId = "--same as Project Name--"
				}
			}
		}

		var hoursSinceLastRun int
		minutesSinceLastRun := 0
		nowHour := time.Now().Hour()
		nowMinute := time.Now().Minute()
		scheduledHour := 6
		if scheduledHour > nowHour {
			hoursSinceLastRun = 24 - nowHour + scheduledHour
			minutesSinceLastRun = nowMinute
		} else {
			hoursSinceLastRun = nowHour
			minutesSinceLastRun = nowMinute
		}
		lastRunTime := ""
		if hoursSinceLastRun > 0 {
			lastRunTime = strconv.Itoa(hoursSinceLastRun) + " hours "
		}
		lastRunTime += strconv.Itoa(minutesSinceLastRun) + " mins"
		nextRunHours := 24 - hoursSinceLastRun - 1
		nextRunMinutes := 60 - minutesSinceLastRun
		nextRunTime := ""
		if nextRunHours > 0 {
			nextRunTime = strconv.Itoa(nextRunHours) + " hours "
		}
		nextRunTime += strconv.Itoa(nextRunMinutes) + " mins"

		core.HTMLWithGlobalState(c, "overview.tmpl", gin.H{
			"HasProjects": len(projectStrings) > 0,
			"LastRunTime": lastRunTime,
			"NextRunTime": nextRunTime,
			"Projects":    projectStrings,
			"HasErrors":   responseError != nil && responseError.Error.Code > 0,
			"Error":       responseError,
			"PageName":    "gcp_overview",
			"Section":     "gae",
		})
	}
}
