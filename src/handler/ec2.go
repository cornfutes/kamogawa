package handler

import (
	"fmt"
	"html/template"
	"kamogawa/aws"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"kamogawa/core"
	"kamogawa/identity"

	"gorm.io/gorm"
)

func EC2(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		user := identity.CheckSessionForUser(c, db)
		if user.AccessToken == nil {
			core.HTMLWithGlobalState(c, "ec2.tmpl", gin.H{
				"NumCachedCalls": 0,
				"Unauthorized":   true,
				"PageName":       "aws_ec2_overview",
				"Section":        "aws",
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)
		if user.Scope == nil || !user.Scope.Valid {
			panic("Missing scope")
		}
		fmt.Printf("User %v\n", user)

		start := time.Now()

		ec2AggregatedInstances := aws.AWSListEC2Instances(db, user, useCache)

		var htmlLines []string
		cachedCalls := 0
		for _, zone := range ec2AggregatedInstances.Zones {
			if len(zone.Instances) == 0 {
				continue
			}

			htmlLines = append(htmlLines, "<li>"+zone.Zone+" ( Zone ) <ul>")
			for _, instance := range zone.Instances {
				htmlLines = append(htmlLines, "<li>"+instance.Name+" ( Instance )</li>")
			}
			htmlLines = append(htmlLines, "</ul></li>")
		}

		duration := time.Since(start)
		core.HTMLWithGlobalState(c, "ec2.tmpl", gin.H{
			"Duration":       duration,
			"NumCachedCalls": cachedCalls,
			"AssetLines":     template.HTML(strings.Join(htmlLines[:], "")),
			"PageName":       "aws_ec2_overview",
			"Section":        "aws",
		})
	}
}
