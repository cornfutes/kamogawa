package handler

import (
	"fmt"
	"kamogawa/core"
	"kamogawa/types/gcp/gaetypes"
	"kamogawa/types/gcp/gcetypes"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SearchResult struct {
	Text string
	Link string

	Product  string
	Provider string
	Kind     string
	Name     string
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

func getRealData(db *gorm.DB, q string) []SearchResult {
	var searchResults []SearchResult

	// This will be buggy. duplicate results?
	for _, word := range strings.Fields(q) {
		r, err := SearchInstances(db, word)
		if err == nil {
			searchResults = append(searchResults, r...)
		}

		r, err = searchProjects(db, word)
		if err == nil {
			searchResults = append(searchResults, r...)
		}

		r, err = searchGAEServices(db, word)
		if err == nil {
			searchResults = append(searchResults, r...)
		}

		r, err = searchGAEVersions(db, word)
		if err == nil {
			searchResults = append(searchResults, r...)
		}
	}
	return searchResults
}

func searchGAEServices(db *gorm.DB, q string) ([]SearchResult, error) {
	var x []gaetypes.GAEServiceDB
	result := db.Raw(""+
		" SELECT * "+
		" FROM gae_service_dbs"+
		" WHERE name || ' ' || id || ' ' || project_id"+
		" ILIKE ?"+
		" LIMIT 50", fmt.Sprintf("%%%v%%", q)).Find(&x)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(x) == 0 {
		fmt.Printf("No results found\n")
		return nil, fmt.Errorf("No results found")
	}

	searchResults := make([]SearchResult, 0, len(x))
	for _, v := range x {
		searchResults = append(searchResults,
			SearchResult{
				Text:     v.ToSearchString(),
				Link:     v.ToLink(),
				Name:     v.Name,
				Provider: "GCP",
				Product:  "GAE",
				Kind:     "Service",
			})
	}

	return searchResults, nil

}

func searchGAEVersions(db *gorm.DB, q string) ([]SearchResult, error) {
	var x []gaetypes.GAEVersionDB
	result := db.Raw(""+
		" SELECT * "+
		" FROM gae_version_dbs"+
		" WHERE id "+
		" ILIKE ?"+
		" LIMIT 50", fmt.Sprintf("%%%v%%", q)).Find(&x)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(x) == 0 {
		fmt.Printf("No results found\n")
		return nil, fmt.Errorf("No results found")
	}

	searchResults := make([]SearchResult, 0, len(x))
	for _, v := range x {
		searchResults = append(searchResults,
			SearchResult{
				Text:     v.ToSearchString(),
				Link:     v.ToLink(),
				Name:     v.Name,
				Provider: "GCP",
				Product:  "GAE",
				Kind:     "Version",
			})
	}

	return searchResults, nil

}

func SearchInstances(db *gorm.DB, q string) ([]SearchResult, error) {
	var gceInstanceDBs []gcetypes.GCEInstanceDB
	result := db.Raw(""+
		" SELECT * "+
		" FROM gce_instance_dbs"+
		" WHERE name || ' ' || id || ' ' || project_id || ' ' || zone"+
		" ILIKE ?"+
		" LIMIT 50", fmt.Sprintf("%%%v%%", q)).Find(&gceInstanceDBs)
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
		searchResults = append(searchResults,
			SearchResult{
				Text:     v.ToSearchString(),
				Link:     v.ToLink(),
				Name:     v.Name,
				Provider: "GCP",
				Product:  "GCE",
				Kind:     "VM",
			})
	}

	return searchResults, nil
}

func searchProjects(db *gorm.DB, q string) ([]SearchResult, error) {
	var projectDBs []gcetypes.ProjectDB
	result := db.Raw(""+
		" SELECT * "+
		" FROM project_dbs"+
		" WHERE (name || ' ' || project_id || ' ' || project_number"+
		" ILIKE ?)"+
		" AND email = '1337gamer@gmail.com'"+
		" LIMIT 50", fmt.Sprintf("%%%v%%", q)).Find(&projectDBs)
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
		searchResults = append(searchResults,
			SearchResult{
				Text:     v.ToSearchString(),
				Link:     v.ToLink(),
				Name:     v.Name,
				Provider: "GCP",
				Product:  "",
				Kind:     "Project",
			})
	}

	return searchResults, nil
}

var minQueryLength int = 1
var minQueryError string = "Query must be 4 characters or more."
var maxQueryLength int = 80
var maxQueryError string = "Query must be 80 characters or less."

func Search(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		originalQ := c.Query("q")
		q := strings.Split(originalQ, ":::")[0]

		if !validateQuery(c, q) {
			return
		}

		searchResults := getRealData(db, q)

		core.HTMLWithGlobalState(c, "search.html", gin.H{
			"HasFilter":  originalQ != q,
			"Error":      nil,
			"IsRegex":    queryIsRegex(q),
			"Query":      originalQ,
			"HasResults": searchResults != nil,
			"Results":    searchResults,
			"IsSearch":   true,
		})
	}
}

func validateQuery(c *gin.Context, q string) bool {
	if len(q) < minQueryLength {
		core.HTMLWithGlobalState(c, "search.html", gin.H{
			"Error":      minQueryError,
			"Query":      q,
			"HasResults": false,
			"Results":    nil,
			"IsSearch":   true,
		})
		return false
	}
	if len(q) > maxQueryLength {
		core.HTMLWithGlobalState(c, "search.html", gin.H{
			"Error":      maxQueryError,
			"Query":      q,
			"HasResults": false,
			"Results":    nil,
			"IsSearch":   true,
		})
		return false
	}

	return true
}

// For displaying warning that we ignore regex right now.
func queryIsRegex(q string) bool {
	isRegex := false
	for i := 0; i < len(q); i++ {
		char := string(q[i])
		if char == "*" || char == "\\" || char == "." || char == "+" || char == "^" || char == "[" || char == "]" {
			isRegex = true
			break
		}
	}
	return isRegex
}
