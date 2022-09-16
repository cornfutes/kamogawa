package handler

import (
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

func Overview(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "overview.html", gin.H{
				"Unauthorized": true,
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)

		responseSuccess, responseError := gcp.GCPListProjects(db, user, useCache)
		if responseError != nil && responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
			core.HTMLWithGlobalState(c, "overview.html", gin.H{
				"MissingScopes": true,
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
		var minutesSinceLastRun = 0
		var nowHour = time.Now().Hour()
		var nowMinute = time.Now().Minute()
		var scheduledHour = 6
		if scheduledHour > nowHour {
			hoursSinceLastRun = 24 - nowHour + scheduledHour
			minutesSinceLastRun = nowMinute
		} else {
			hoursSinceLastRun = nowHour
			minutesSinceLastRun = nowMinute
		}
		var lastRunTime = ""
		if hoursSinceLastRun > 0 {
			lastRunTime = strconv.Itoa(hoursSinceLastRun) + " hours "
		}
		lastRunTime += strconv.Itoa(minutesSinceLastRun) + " mins"
		var nextRunHours = 24 - hoursSinceLastRun - 1
		var nextRunMinutes = 60 - minutesSinceLastRun
		var nextRunTime = ""
		if nextRunHours > 0 {
			nextRunTime = strconv.Itoa(nextRunHours) + " hours "
		}
		nextRunTime += strconv.Itoa(nextRunMinutes) + " mins"

		core.HTMLWithGlobalState(c, "overview.html", gin.H{
			"HasProjects": len(projectStrings) > 0,
			"LastRunTime": lastRunTime,
			"NextRunTime": nextRunTime,
			"Projects":    projectStrings,
			"HasErrors":   responseError != nil && responseError.Error.Code > 0,
			"Error":       responseError,
		})
	}
}
