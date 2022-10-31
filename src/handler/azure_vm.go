package handler

import (
	"html/template"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kamogawa/core"
)

func AzureVM(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		time.Sleep(5 * time.Millisecond)

		var htmlLines []string
		cachedCalls := 10

		htmlLines = append(htmlLines, "<li>us-west-1a ( Zone ) <ul>")
		htmlLines = append(htmlLines, "<li>shimogawa (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "<li>akari (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "<li>ichiban (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "<li>ichiro (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "</ul></li>")

		htmlLines = append(htmlLines, "<li>us-west-1b ( Zone ) <ul>")
		htmlLines = append(htmlLines, "<li>kaze (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "<li>oni (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "<li>moku (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "<li>mizu (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "<li>oto (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "<li>kumo (Virtual Machine)</li>")
		htmlLines = append(htmlLines, "</ul></li>")

		duration := time.Since(start)
		core.HTMLWithGlobalState(c, "azure_vm.tmpl", gin.H{
			"Duration":       duration,
			"NumCachedCalls": cachedCalls,
			"AssetLines":     template.HTML(strings.Join(htmlLines[:], "")),
			"PageName":       "azure_vm_overview",
			"Section":        "azure",
		})
	}
}
