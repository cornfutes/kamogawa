package handler

import (
	"html/template"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kamogawa/core"
)

func S3(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		time.Sleep(5 * time.Millisecond)

		var htmlLines []string
		cachedCalls := 4

		htmlLines = append(htmlLines, "<li>us-west-1a ( Zone ) <ul>")
		htmlLines = append(htmlLines, "<li>IDs (S3 Bucket)</li>")
		htmlLines = append(htmlLines, "<li>Addresses (S3 Bucket)</li>")
		htmlLines = append(htmlLines, "</ul></li>")

		htmlLines = append(htmlLines, "<li>us-west-1b ( Zone ) <ul>")
		htmlLines = append(htmlLines, "<li>Purchases (S3 Bucket)</li>")
		htmlLines = append(htmlLines, "<li>Payment-Methods (S3 Bucket)</li>")
		htmlLines = append(htmlLines, "</ul></li>")

		duration := time.Since(start)
		core.HTMLWithGlobalState(c, "s3.tmpl", gin.H{
			"Duration":       duration,
			"NumCachedCalls": cachedCalls,
			"AssetLines":     template.HTML(strings.Join(htmlLines[:], "")),
			"PageName":       "aws_s3_overview",
			"Section":        "aws",
		})
	}
}
