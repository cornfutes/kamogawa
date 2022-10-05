package handler

import (
	"fmt"
	"kamogawa/aws"
	"net/url"
	"strings"
	"time"

	"kamogawa/core"
	"kamogawa/identity"
	"kamogawa/types"
	"kamogawa/types/gcp/coretypes"
	"kamogawa/types/gcp/gaetypes"
	"kamogawa/types/gcp/gcetypes"

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

const (
	SERPPageSize          = 10
	resultLimit           = 50
	minQueryLength int    = 1
	minQueryError  string = "Query must be 4 characters or more."
	maxQueryLength int    = 80
	maxQueryError  string = "Query must be 80 characters or less."
)

func Search(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := identity.CheckSessionForUser(c, db)
		// TODO: Handle non-GCP users
		if user.AccessToken == nil || !user.Gmail.Valid {
			core.HTMLWithGlobalState(c, "search.tmpl", gin.H{
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
		isGcp := c.Query("p") == "gcp"
		isAws := c.Query("p") == "aws"
		isAzure := c.Query("p") == "azure"
		isMore := c.Query("p") == "more"
		isAll := !isGcp && !isAws && !isAzure && !isMore

		start := time.Now()
		var allSearchResults []SearchResult
		if isAll || isGcp {
			allSearchResults = getSearchResults(db, user, q)
		} else {
			allSearchResults = nil
		}
		duration := time.Since(start)
		var results []SearchResult
		if len(allSearchResults) > SERPPageSize {
			results = allSearchResults[:SERPPageSize]
		} else {
			results = allSearchResults
		}

		fmt.Printf("%v\n", isAll)
		numTotalResults := len(allSearchResults)
		core.HTMLWithGlobalState(c, "search.tmpl", gin.H{
			"HasFilter":         originalQ != q,
			"Error":             nil,
			"IsRegex":           queryIsRegex(q),
			"Query":             originalQ,
			"HasResults":        results != nil,
			"Results":           results,
			"CountTotalResults": numTotalResults,
			"Duration":          duration,
			"IsSearch":          "yes",
			"SearchUrlBase":     "/search?q=" + url.QueryEscape(originalQ),
			"SearchUrlGCP":      "/search?p=gcp&q=" + url.QueryEscape(originalQ),
			"SearchUrlAWS":      "/search?p=aws&q=" + url.QueryEscape(originalQ),
			"SearchUrlAzure":    "/search?p=azure&q=" + url.QueryEscape(originalQ),
			"SearchUrlMore":     "/search?p=more&q=" + url.QueryEscape(originalQ),
			"IsAll":             isAll,
			"IsGCP":             isGcp,
			"IsAWS":             isAws,
			"IsAzure":           isAzure,
			"IsMore":            isMore,
		})
	}
}

func validateQuery(c *gin.Context, q string) bool {
	if len(q) < minQueryLength {
		core.HTMLWithGlobalState(c, "search.tmpl", gin.H{
			"Error":      minQueryError,
			"Query":      q,
			"HasResults": false,
			"Results":    nil,
			"IsSearch":   "yes",
		})
		return false
	}
	if len(q) > maxQueryLength {
		core.HTMLWithGlobalState(c, "search.tmpl", gin.H{
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

	r, err = SearchGCEInstances(db, user, word)
	if err == nil {
		searchResults = append(searchResults, r...)
	}

	r, err = searchGAEServices(db, user, word)
	if err == nil {
		searchResults = append(searchResults, r...)
	}

	r, err = searchGAEVersions(db, user, word)
	if err == nil {
		searchResults = append(searchResults, r...)
	}
	r, err = searchEC2Instances(db, user, word)
	if err == nil {
		searchResults = append(searchResults, r...)
	}
	// }
	return searchResults
}

func searchProjects(db *gorm.DB, user types.User, q string) ([]SearchResult, error) {
	var projectDBs []coretypes.ProjectDB
	result := db.Raw(
		" SELECT project_dbs.* "+
			" FROM project_dbs"+
			" INNER JOIN project_auths "+
			" ON project_auths.project_id = project_dbs.project_id"+
			" AND project_auths.gmail = ?"+
			" AND (project_dbs.name || ' ' || project_dbs.project_id) ILIKE ?"+
			" ORDER BY project_dbs.name, project_dbs.project_id"+
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

func SearchGCEInstances(db *gorm.DB, user types.User, q string) ([]SearchResult, error) {
	var gceInstanceDBs []gcetypes.GCEInstanceDB
	result := db.Raw(""+
		" SELECT gce_instance_dbs.* "+
		" FROM gce_instance_dbs"+
		" INNER JOIN gce_instance_auths "+
		" ON gce_instance_auths.id = gce_instance_dbs.id"+
		" AND gce_instance_auths.gmail = ?"+
		" AND (gce_instance_dbs.id || ' ' || gce_instance_dbs.name || ' ' || gce_instance_dbs.status || ' ' || gce_instance_dbs.project_id || ' ' || gce_instance_dbs.zone) ILIKE ?"+
		" ORDER BY gce_instance_dbs.name, gce_instance_dbs.id"+
		" LIMIT ?", user.Gmail.String, fmt.Sprintf("%%%v%%", q), resultLimit).Find(&gceInstanceDBs)
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

func searchGAEServices(db *gorm.DB, user types.User, q string) ([]SearchResult, error) {
	var gaeServiceDBs []gaetypes.GAEServiceDB
	result := db.Raw(""+
		" SELECT gae_service_dbs.* "+
		" FROM gae_service_dbs"+
		" INNER JOIN gae_service_auths "+
		" ON gae_service_auths.id = gae_service_dbs.id"+
		" AND gae_service_auths.gmail = ?"+
		" AND (gae_service_dbs.name || ' ' || gae_service_dbs.id || ' ' || gae_service_dbs.project_id) ILIKE ?"+
		" ORDER BY gae_service_dbs.name, gae_service_dbs.id"+
		" LIMIT ?", user.Gmail.String, fmt.Sprintf("%%%v%%", q), resultLimit).Find(&gaeServiceDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(gaeServiceDBs) == 0 {
		fmt.Printf("No results found\n")
		return nil, fmt.Errorf("No results found")
	}

	searchResults := make([]SearchResult, 0, len(gaeServiceDBs))
	for _, v := range gaeServiceDBs {
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

func searchGAEVersions(db *gorm.DB, user types.User, q string) ([]SearchResult, error) {
	var gaeVersionDBs []gaetypes.GAEVersionDB
	result := db.Raw(""+
		" SELECT gae_version_dbs.* "+
		" FROM gae_version_dbs"+
		" INNER JOIN gae_version_auths "+
		" ON gae_version_auths.id = gae_version_dbs.id"+
		" AND gae_version_auths.gmail = ?"+
		" AND (gae_version_dbs.name || ' ' || gae_version_dbs.id || ' ' || gae_version_dbs.project_id || ' ' || gae_version_dbs.service_name || ' ' || gae_version_dbs.service_id || ' ' || gae_version_dbs.serving_status) ILIKE ?"+
		" ORDER BY gae_version_dbs.name, gae_version_dbs.id"+
		" LIMIT ?", user.Gmail.String, fmt.Sprintf("%%%v%%", q), resultLimit).Find(&gaeVersionDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(gaeVersionDBs) == 0 {
		fmt.Printf("No results found\n")
		return nil, fmt.Errorf("No results found")
	}

	searchResults := make([]SearchResult, 0, len(gaeVersionDBs))
	for _, v := range gaeVersionDBs {
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

func searchEC2Instances(db *gorm.DB, user types.User, q string) ([]SearchResult, error) {
	ec2AggregatedInstances := aws.AWSListEC2Instances(db, user, true)
	searchResults := make([]SearchResult, 0)
	for _, ec2Zone := range ec2AggregatedInstances.Zones {
		for _, instance := range ec2Zone.Instances {
			if strings.Contains(instance.Id, q) {
				searchResults = append(searchResults,
					SearchResult{
						Text:     fmt.Sprintf("Type: EC2 Instance, Name: %v, Id: %v, Zone: %v", instance.Name, instance.Id, ec2Zone.Zone),
						Link:     "https://us-west-2.console.aws.amazon.com/ec2/home?region=" + ec2Zone.Zone[:len(ec2Zone.Zone)-2] + "#Instances:instanceId=" + instance.Id,
						Name:     instance.Name,
						Provider: "AWS",
						Product:  "EC2",
						Kind:     "Instance",
					})
			}
		}
	}

	return searchResults, nil
}
