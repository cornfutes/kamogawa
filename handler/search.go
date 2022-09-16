package handler

import (
	"kamogawa/core"

	"github.com/gin-gonic/gin"
)

type SearchResult struct {
	Text string
	Link string
}

func getFakeData(q string) []SearchResult {
	searchResults := [2]SearchResult{}
	searchResults[0] = SearchResult{
		Text: "/projects/kanazawa2/versions/20220830t021415",
		Link: "localhost",
	}
	searchResults[1] = SearchResult{
		Text: "/projects/kanazawa/versions/20220829t195029",
		Link: "localhost",
	}
	return searchResults[:]
}

// TODO(david): implement
func getRealData(q string) []SearchResult {

	return nil
}

func Search(c *gin.Context) {
	q := c.Query("q")

	searchResults := getRealData(q)
	if searchResults == nil {
		searchResults = getFakeData(q)
	}

	// Renders HTML
	core.HTMLWithGlobalState(c, "search.html", gin.H{
		"Query":   q,
		"Results": searchResults,
	})
}
