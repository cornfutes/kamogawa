package handler

import (
	"html/template"
	"kamogawa/aws"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kamogawa/core"
)

func EC2(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		useCache := len(c.Query("r")) == 0

		start := time.Now()

		ec2AggregatedInstances := aws.AWSListEC2Instances(db, useCache)

		var htmlLines []string
		cachedCalls := 10
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
