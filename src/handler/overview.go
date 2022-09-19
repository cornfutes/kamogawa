package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"
	"strconv"
	"strings"
	"time"
)

func Overview(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "overview.tmpl", gin.H{
				"Unauthorized": true,
				"PageName":     "gcp_overview",
				"Section":      "gcp",
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)
		if user.Scope == nil || !user.Scope.Valid {
			panic("Missing scope")
		}

		projectDBs, responseError := gcp.GCPListProjects(db, user, useCache)
		if responseError != nil && responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
			core.HTMLWithGlobalState(c, "overview.tmpl", gin.H{
				"MissingScopes": true,
				"PageName":      "gcp_overview",
				"Section":      "gcp",
			})
			return
		}

		for i, p := range projectDBs {
			if p.ProjectId == p.Name {
				projectDBs[i].ProjectId = "--same as Project Name--"
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
			"HasProjects": len(projectDBs) > 0,
			"LastRunTime": lastRunTime,
			"NextRunTime": nextRunTime,
			"Projects":    projectDBs,
			"HasErrors":   responseError != nil && responseError.Error.Code > 0,
			"Error":       responseError,
			"PageName":    "gcp_overview",
			"Section":     "gcp",
		})
	}
}
