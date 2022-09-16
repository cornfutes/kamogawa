package handler

import (
	"fmt"
	"kamogawa/core"
	"kamogawa/types"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
func getRealData(db *gorm.DB, q string) []SearchResult {
	var searchResults []SearchResult

	for _, word := range strings.Fields(q) {
		r, err := SearchInstances(db, word)
		if err == nil {
			searchResults = append(searchResults, r...)
		}

		r, err = searchProjects(db, word)
		if err == nil {
			searchResults = append(searchResults, r...)
		}
	}
	return searchResults
}

func SearchInstances(db *gorm.DB, q string) ([]SearchResult, error) {
	var gceInstanceDBs []types.GCEInstanceDB
	result := db.Raw(""+
		" SELECT * "+
		" FROM gce_instance_dbs"+
		" WHERE name || ' ' || id || ' ' || project_id || ' ' || zone"+
		" ILIKE ?", fmt.Sprintf("%%%v%%", q)).Find(&gceInstanceDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(gceInstanceDBs) == 0 {
		fmt.Printf("No results found\n")
		return nil, fmt.Errorf("No results found")
	}

	searchResults := make([]SearchResult, 0, len(gceInstanceDBs))
	for _, v := range gceInstanceDBs {
		searchResults = append(searchResults, SearchResult{Text: v.ToSearchString(), Link: v.ToLink()})
	}

	return searchResults, nil
}

func searchProjects(db *gorm.DB, q string) ([]SearchResult, error) {
	var projectDBs []types.ProjectDB
	result := db.Raw(""+
		" SELECT * "+
		" FROM project_dbs"+
		" WHERE name || ' ' || project_id || ' ' || project_number"+
		" ILIKE ?", fmt.Sprintf("%%%v%%", q)).Find(&projectDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(projectDBs) == 0 {
		fmt.Printf("No results found\n")
		return nil, fmt.Errorf("No results found")
	}

	searchResults := make([]SearchResult, 0, len(projectDBs))
	for _, v := range projectDBs {
		searchResults = append(searchResults, SearchResult{Text: v.ToSearchString(), Link: v.ToLink()})
	}

	return searchResults, nil
}

func Search(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		q := c.Query("q")

		searchResults := getRealData(db, q)
		if searchResults == nil {
			searchResults = getFakeData(q)
		}

		// Renders HTML
		core.HTMLWithGlobalState(c, "search.html", gin.H{
			"Query":   q,
			"Results": searchResults,
		})
	}
}
