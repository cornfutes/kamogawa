package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kamogawa/core"
	"kamogawa/identity"
	"kamogawa/types"
	"kamogawa/types/gcp/gaetypes"
	"kamogawa/types/gcp/gcetypes"
	"strings"
	"time"
)

type SearchResult struct {
	Text string
	Link string

	Product  string
	Provider string
	Kind     string
	Name     string
}

const SERPPageSize = 10
const resultLimit = 50
const minQueryLength int = 1
const minQueryError string = "Query must be 4 characters or more."
const maxQueryLength int = 80
const maxQueryError string = "Query must be 80 characters or less."

func Search(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := identity.CheckSessionForUser(c, db)
		// TODO: Handle non-GCP users
		if user.AccessToken == nil || !user.Gmail.Valid {
			core.HTMLWithGlobalState(c, "search.html", gin.H{
				"Unauthorized": true,
			})
			return
		}
		identity.CheckDBAndRefreshToken(c, user, db)

		originalQ := c.Query("q")
		q := strings.Split(originalQ, ":::")[0]

		if !validateQuery(c, q) {
			return
		}

		start := time.Now()
		allSearchResults := getSearchResults(db, user, q)
		duration := time.Since(start)
		var results []SearchResult
		if len(allSearchResults) > SERPPageSize {
			results = allSearchResults[:SERPPageSize]
		} else {
			results = allSearchResults
		}
		numTotalResults := len(allSearchResults)
		core.HTMLWithGlobalState(c, "search.html", gin.H{
			"HasFilter":         originalQ != q,
			"Error":             nil,
			"IsRegex":           queryIsRegex(q),
			"Query":             originalQ,
			"HasResults":        results != nil,
			"Results":           results,
			"CountTotalResults": numTotalResults,
			"Duration":          duration,
			"IsSearch":          "yes",
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
			"IsSearch":   "yes",
		})
		return false
	}
	if len(q) > maxQueryLength {
		core.HTMLWithGlobalState(c, "search.html", gin.H{
			"Error":      maxQueryError,
			"Query":      q,
			"HasResults": false,
			"Results":    nil,
			"IsSearch":   "yes",
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

func getSearchResults(db *gorm.DB, user types.User, q string) []SearchResult {
	var searchResults []SearchResult

	// TODO: renable multiworld. duplicate results
	word := strings.Fields(q)[0]
	// for _, word := range strings.Fields(q) {
	r, err := searchProjects(db, user, word)
	if err == nil {
		searchResults = append(searchResults, r...)
	}

	r, err = SearchGCEInstances(db, word)
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
	// }
	return searchResults
}

func searchProjects(db *gorm.DB, user types.User, q string) ([]SearchResult, error) {
	var projectDBs []gcetypes.ProjectDB
	result := db.Raw(""+
		" SELECT project_dbs.* "+
		" FROM ("+
		"   project_dbs"+
		"   INNER JOIN project_auths "+
		"   ON project_auths.project_id = project_dbs.project_id"+
		"   AND project_auths.gmail = ?"+
		" ) WHERE (project_dbs.name || ' ' || project_dbs.project_id"+
		" ILIKE ?)"+
		" LIMIT ?", user.Gmail.String, fmt.Sprintf("%%%v%%", q), resultLimit).Find(&projectDBs)
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

func SearchGCEInstances(db *gorm.DB, q string) ([]SearchResult, error) {
	var gceInstanceDBs []gcetypes.GCEInstanceDB
	result := db.Raw(""+
		" SELECT * "+
		" FROM gce_instance_dbs"+
		" WHERE name || ' ' || id || ' ' || project_id || ' ' || zone"+
		" ILIKE ?"+
		" LIMIT ?", fmt.Sprintf("%%%v%%", q), resultLimit).Find(&gceInstanceDBs)
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

func searchGAEServices(db *gorm.DB, q string) ([]SearchResult, error) {
	var results []gaetypes.GAEServiceDB
	result := db.Raw(""+
		" SELECT * "+
		" FROM gae_service_dbs"+
		" WHERE name || ' ' || id || ' ' || project_id"+
		" ILIKE ?"+
		" LIMIT ?", fmt.Sprintf("%%%v%%", q), resultLimit).Find(&results)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(results) == 0 {
		fmt.Printf("No results found\n")
		return nil, fmt.Errorf("No results found")
	}

	searchResults := make([]SearchResult, 0, len(results))
	for _, v := range results {
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
		" LIMIT ?", fmt.Sprintf("%%%v%%", q), resultLimit).Find(&x)
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
