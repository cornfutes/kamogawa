package gcecache

import (
	"fmt"
	"gorm.io/gorm"
	"kamogawa/types"
	"kamogawa/types/gcp/coretypes"
	"kamogawa/types/gcp/gcetypes"
	"kamogawa/utility"
)

func ReadProjectsCache(db *gorm.DB, user types.User) (*gcetypes.ListProjectResponse, error) {
	var projectDBs []coretypes.ProjectDB

	result := db.Joins(
		" INNER JOIN project_auths"+
			" ON project_auths.project_id = project_dbs.project_id"+
			" AND project_auths.gmail = ? ", user.Gmail.String,
	).Order("project_dbs.name, project_dbs.project_id").Find(&projectDBs)
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

func ReadProjectsCache2(db *gorm.DB, user types.User) []coretypes.ProjectDB {
	var projectDBs []coretypes.ProjectDB

	result := db.Joins(
		" INNER JOIN project_auths"+
			" ON project_auths.project_id = project_dbs.project_id"+
			" AND project_auths.gmail = ? ", user.Gmail.String,
	).Order("project_dbs.name, project_dbs.project_id").Find(&projectDBs)
	if result.Error != nil {
		panic("Query failed")
	}

	if len(projectDBs) == 0 {
		fmt.Printf("Cache miss\n")
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

	projectDBs := make([]coretypes.ProjectDB, 0, len(resp.Projects))
	for _, v := range resp.Projects {
		projectDBs = append(projectDBs, coretypes.ProjectToProjectDB(&v))
	}

	utility.BatchUpsertOrElseSingleProjectDB(db, projectDBs)
	WriteProjectsAuth(db, user, resp)
}

func WriteProjectsCache2(db *gorm.DB, user types.User, projectDBs []coretypes.ProjectDB) {
	if len(projectDBs) == 0 {
		return
	}

	utility.BatchUpsertOrElseSingleProjectDB(db, projectDBs)
	WriteProjectsAuth2(db, user, projectDBs)
}

func WriteProjectsAuth(db *gorm.DB, user types.User, resp *gcetypes.ListProjectResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Projects) == 0 {
		return
	}

	projectAuths := make([]coretypes.ProjectAuth, 0, len(resp.Projects))
	for _, v := range resp.Projects {
		projectAuths = append(projectAuths, coretypes.ProjectAuth{Gmail: user.Gmail.String, ProjectId: v.ProjectId})
	}

	utility.BatchUpsertOrElseSingleProjectAuth(db, projectAuths)
}

func WriteProjectsAuth2(db *gorm.DB, user types.User, projects []coretypes.ProjectDB) {
	if len(projects) == 0 {
		return
	}

	projectAuths := make([]coretypes.ProjectAuth, 0, len(projects))
	for _, v := range projects {
		projectAuths = append(projectAuths, coretypes.ProjectAuth{Gmail: user.Gmail.String, ProjectId: v.ProjectId})
	}

	utility.BatchUpsertOrElseSingleProjectAuth(db, projectAuths)
}

func ReadInstancesCache(db *gorm.DB, user types.User, projectId string) (*gcetypes.GCEAggregatedInstances, error) {
	var gceInstanceDBs []gcetypes.GCEInstanceDB
	result := db.Joins(
		" INNER JOIN gce_instance_auths"+
			" ON gce_instance_auths.id = gce_instance_dbs.id"+
			" AND gce_instance_auths.gmail = ?"+
			" AND gce_instance_dbs.project_id = ?", user.Gmail.String, projectId,
	).Order("gce_instance_dbs.name, gce_instance_dbs.id").Find(&gceInstanceDBs)
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
	utility.BatchUpsertOrElseSingleGCEInstanceDB(db, gceInstanceDBs)
	WriteInstancesAuth(db, user, gceInstanceDBs)
}

func WriteInstancesAuth(db *gorm.DB, user types.User, gceInstanceDBs []gcetypes.GCEInstanceDB) {
	if len(gceInstanceDBs) == 0 {
		return
	}

	instanceAuths := make([]gcetypes.GCEInstanceAuth, 0, len(gceInstanceDBs))
	for _, v := range gceInstanceDBs {
		instanceAuths = append(instanceAuths, gcetypes.GCEInstanceAuth{Gmail: user.Gmail.String, Id: v.Id})
	}

	utility.BatchUpsertOrElseSingleGCEInstanceAuth(db, instanceAuths)
}
