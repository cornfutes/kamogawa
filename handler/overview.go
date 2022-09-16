package handler

import (
	"database/sql"
	"fmt"
	"kamogawa/core"
	"kamogawa/gcp"
	"kamogawa/identity"
	"kamogawa/types"
	"strconv"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func mockGCPListProjects(user types.User) (*types.ListProjectResponse, *types.ErrorResponse) {
	var errorResponse *types.ErrorResponse
	var listProjectResponse *types.ListProjectResponse = &types.ListProjectResponse{}

	var projects = make([]types.Project, 5)
	projects[0] = types.Project{
		Name:      "diceduckmonk1",
		ProjectId: "diceduckmonk1",
	}
	projects[1] = types.Project{
		Name:      "diceduckmonk",
		ProjectId: "diceduckmonk",
	}
	projects[2] = types.Project{
		Name:      "kanazawa2",
		ProjectId: "linear-cinema-360910",
	}
	projects[3] = types.Project{
		Name:      "kanazawa",
		ProjectId: "kanazawa-360910",
	}
	projects[4] = types.Project{
		Name:      "tokyo",
		ProjectId: "tokyo-360910",
	}
	projects[4] = types.Project{
		Name:      "zerotheta1337",
		ProjectId: "zerotheta1337",
	}
	listProjectResponse.Projects = projects

	return listProjectResponse, errorResponse
}

func Overview(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		var email, exists = c.Get(identity.IdentityContextkey)
		if !exists {
			panic("Unexpected")
		}
		var user types.User
		err := db.First(&user, "email = ?", email).Error
		if err != nil {
			panic("Unvalid DB state")
		}
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "overview.html", gin.H{
				"Unauthorized": true,
			})
			return
		}

		responseSuccess, responseError := gcp.GCPListProjects(db, user, useCache)
		if responseError != nil && responseError.Error.Code == 403 && strings.HasPrefix(responseError.Error.Message, "Request had insufficient authentication scopes.") {
			core.HTMLWithGlobalState(c, "overview.html", gin.H{
				"MissingScopes": true,
			})
			return
		}
		// TODO: this bakes in assumption that responseError not nil IFF 401.
		if responseError != nil && responseError.Error.Code > 0 {
			fmt.Printf("Retrying project fetch by first refreshing access token")
			respRefreshtoken := identity.GCPRefresh(c, db)
			user.AccessToken = &sql.NullString{String: respRefreshtoken.AccessToken, Valid: true}
			db.Save(&user)
			responseSuccess, _ = gcp.GCPListProjects(db, user, useCache)
		}

		var projectStrings []types.Project
		if user.Scope == nil || !user.Scope.Valid {
			projectStrings = []types.Project{}
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
