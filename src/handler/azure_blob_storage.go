package handler

import (
	"html/template"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kamogawa/core"
)

func AzureBlobStorage(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		time.Sleep(5 * time.Millisecond)

		var htmlLines []string
		cachedCalls := 10

		htmlLines = append(htmlLines, "<li>Images ( Container ) <ul>")
		htmlLines = append(htmlLines, "<li>img001.png (Blob)</li>")
		htmlLines = append(htmlLines, "<li>img002.png (Blob)</li>")
		htmlLines = append(htmlLines, "<li>img003.png (Blob)</li>")
		htmlLines = append(htmlLines, "<li>img004.png (Blob)</li>")
		htmlLines = append(htmlLines, "<li>img005.png (Blob)</li>")
		htmlLines = append(htmlLines, "</ul></li>")

		htmlLines = append(htmlLines, "<li>Videos ( Container ) <ul>")
		htmlLines = append(htmlLines, "<li>mov1.avi (Blob)</li>")
		htmlLines = append(htmlLines, "<li>mov2.avi (Blob)</li>")
		htmlLines = append(htmlLines, "<li>mov3.avi (Blob)</li>")
		htmlLines = append(htmlLines, "<li>mov4.avi (Blob)</li>")
		htmlLines = append(htmlLines, "<li>mov5.avi (Blob)</li>")
		htmlLines = append(htmlLines, "</ul></li>")

		duration := time.Since(start)
		core.HTMLWithGlobalState(c, "azure_blob_storage.tmpl", gin.H{
			"Duration":       duration,
			"NumCachedCalls": cachedCalls,
			"AssetLines":     template.HTML(strings.Join(htmlLines[:], "")),
			"PageName":       "azure_blob_storage_overview",
			"Section":        "azure",
		})
	}
}
