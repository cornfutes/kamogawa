package gcecache

import (
	"fmt"
	"gorm.io/gorm/clause"
	"kamogawa/types"
	"kamogawa/types/gcp/gcetypes"

	"gorm.io/gorm"
)

func ReadProjectsCache(db *gorm.DB, user types.User) (*gcetypes.ListProjectResponse, error) {
	var projectDBs []gcetypes.ProjectDB

	result := db.Joins(
		" INNER JOIN project_auths "+
			" ON project_auths.project_id = project_dbs.project_id"+
			" AND project_auths.gmail = ? ", user.Gmail.String).Find(&projectDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(projectDBs) == 0 {
		fmt.Printf("Cache miss\n")
		return nil, fmt.Errorf("Cache miss")
	}
	fmt.Printf("Cache hit\n")

	projects := make([]gcetypes.Project, 0, len(projectDBs))
	for _, projectDB := range projectDBs {
		projects = append(projects, projectDB.ToProject())
	}

	return &gcetypes.ListProjectResponse{Projects: projects}, nil
}

func ReadProjectsCache2(db *gorm.DB, user types.User) []gcetypes.ProjectDB {
	var projectDBs []gcetypes.ProjectDB

	result := db.Joins(
		" INNER JOIN project_auths "+
			" ON project_auths.project_id = project_dbs.project_id"+
			" AND project_auths.gmail = ? ", user.Gmail.String).Find(&projectDBs)
	if result.Error != nil {
		panic("Query failed")
	}

	if len(projectDBs) == 0 {
		fmt.Printf("Cache miss\n")
		return nil
	}
	fmt.Printf("Cache hit\n")

	return projectDBs
}

func WriteProjectsCache(db *gorm.DB, user types.User, resp *gcetypes.ListProjectResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Projects) == 0 {
		return
	}

	projectDBs := make([]gcetypes.ProjectDB, 0, len(resp.Projects))
	for _, v := range resp.Projects {
		projectDBs = append(projectDBs, gcetypes.ProjectToProjectDB(&v, true))
	}
	for _, projectDB := range projectDBs {
		db.FirstOrCreate(&projectDB)
	}

	WriteProjectsAuth(db, user, resp)
}

func WriteProjectsCache2(db *gorm.DB, user types.User, projects []gcetypes.ProjectDB) {
	if len(projects) == 0 {
		return
	}
	for _, projectDB := range projects {
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&projectDB)
	}

	WriteProjectsAuth2(db, user, projects)
}

func WriteProjectsAuth(db *gorm.DB, user types.User, resp *gcetypes.ListProjectResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Projects) == 0 {
		return
	}

	projectAuths := make([]gcetypes.ProjectAuth, 0, len(resp.Projects))
	for _, v := range resp.Projects {
		projectAuths = append(projectAuths, gcetypes.ProjectAuth{Gmail: user.Gmail.String, ProjectId: v.ProjectId})
	}
	for _, v := range projectAuths {
		db.FirstOrCreate(&v)
	}
}

func WriteProjectsAuth2(db *gorm.DB, user types.User, projects []gcetypes.ProjectDB) {
	if len(projects) == 0 {
		return
	}

	projectAuths := make([]gcetypes.ProjectAuth, 0, len(projects))
	for _, v := range projects {
		projectAuths = append(projectAuths, gcetypes.ProjectAuth{Gmail: user.Gmail.String, ProjectId: v.ProjectId})
	}
	for _, v := range projectAuths {
		db.FirstOrCreate(&v)
	}
}

func ReadInstancesCache(db *gorm.DB, projectId string) (*gcetypes.GCEAggregatedInstances, error) {
	var gceInstanceDBs []gcetypes.GCEInstanceDB
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

	gceAggregatedInstances := gcetypes.GCEInstanceDBToGCEAggregatedInstances(gceInstanceDBs)
	return &gceAggregatedInstances, nil
}

func WriteInstancesCache(db *gorm.DB, user types.User, projectId string, resp *gcetypes.GCEAggregatedInstances) {
	if resp == nil {
		panic("Unexpected list of instances")
	}
	if len(resp.Zones) == 0 {
		return
	}

	gceInstanceDBs := gcetypes.GCEAggregatedInstancesToGCEInstanceDB(user, projectId, resp)

	for _, gceInstanceDB := range gceInstanceDBs {
		db.FirstOrCreate(&gceInstanceDB)
	}

}
