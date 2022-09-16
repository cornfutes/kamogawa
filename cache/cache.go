package cache

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"kamogawa/types"
)

func ReadInstancesCache(db *gorm.DB, projectId string) (*types.GCEAggregatedInstances, error) {
	var gceInstanceDBs []types.GCEInstanceDB
	result := db.Where("project_id = ?", projectId).Find(&gceInstanceDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(gceInstanceDBs) == 0 {
		fmt.Printf("Cache miss\n")
		return nil, fmt.Errorf("Cache miss")
	}
	fmt.Printf("Cache hit\n")

	gceAggregatedInstances := types.GCEInstanceDBToGCEAggregatedInstances(gceInstanceDBs)
	return &gceAggregatedInstances, nil
}

func WriteInstancesCache(db *gorm.DB, user types.User, projectId string, resp types.GCEAggregatedInstances) {
	if len(resp.Zones) == 0 {
		return
	}

	gceInstanceDBs := types.GCEAggregatedInstancesToGCEInstanceDB(user, projectId, resp)

	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&gceInstanceDBs)
	if result.Error != nil {
		fmt.Printf("Batch insert failed, trying individual\n")
		for _, gceInstanceDB := range gceInstanceDBs {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&gceInstanceDB)
		}
	}
}

func ReadProjectsCache(db *gorm.DB, user types.User) (*types.ListProjectResponse, error) {
	var projectDBs []types.ProjectDB
	result := db.Where("email = ?", user.Email).Find(&projectDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(projectDBs) == 0 {
		fmt.Printf("Cache miss\n")
		return nil, fmt.Errorf("Cache miss")
	}
	fmt.Printf("Cache hit\n")

	projects := make([]types.Project, 0, len(projectDBs))
	for _, projectDB := range projectDBs {
		projects = append(projects, projectDB.ToProject())
	}

	return &types.ListProjectResponse{Projects: projects}, nil
}

func WriteProjectsCache(db *gorm.DB, user types.User, resp *types.ListProjectResponse) {
	if len(resp.Projects) == 0 {
		return
	}

	projectDBs := make([]types.ProjectDB, 0, len(resp.Projects))

	for _, v := range resp.Projects {
		projectDBs = append(projectDBs, v.ToProjectDB(user.Email))
	}

	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&projectDBs)
	if result.Error != nil {
		fmt.Printf("Batch insert failed, trying individual\n")
		for _, projectDB := range projectDBs {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&projectDB)
		}
	}
}
