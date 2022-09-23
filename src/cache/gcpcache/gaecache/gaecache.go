package gaecache

import (
	"fmt"
	"gorm.io/gorm"
	"kamogawa/types"
	"kamogawa/types/gcp/gaetypes"
	"kamogawa/utility"
)

func ReadServicesCache(db *gorm.DB, user types.User, projectId string) (*gaetypes.GAEListServicesResponse, error) {
	var gaeServiceDBs []gaetypes.GAEServiceDB
	result := db.Joins(
		" INNER JOIN gae_service_auths"+
			" ON gae_service_auths.id = gae_service_dbs.id"+
			" AND gae_service_auths.gmail = ?"+
			" AND gae_service_dbs.project_id = ?", user.Gmail.String, projectId,
	).Order("gae_service_dbs.name, gae_service_dbs.id").Find(&gaeServiceDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(gaeServiceDBs) == 0 {
		fmt.Printf("Cache miss\n")
		return nil, fmt.Errorf("Cache miss")
	}
	fmt.Printf("Cache hit\n")

	services := make([]gaetypes.GAEService, 0, len(gaeServiceDBs))
	for _, serviceDB := range gaeServiceDBs {
		services = append(services,
			gaetypes.GAEService{
				// GCP nomenclature is misleading.
				Name:  serviceDB.Id,
				Id:    serviceDB.Name,
				Split: gaetypes.GAEServiceTrafficAllocations{},
			})
	}

	return &gaetypes.GAEListServicesResponse{Services: services}, nil
}

func WriteServicesCache(db *gorm.DB, user types.User, projectId string, resp *gaetypes.GAEListServicesResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Services) == 0 {
		return
	}

	serviceDBs := make([]gaetypes.GAEServiceDB, 0, len(resp.Services))
	for _, v := range resp.Services {
		serviceDBs = append(serviceDBs, gaetypes.GAEServiceDB{
			// GAE API is misleading
			Name:      v.Id,   // The ID is the name, not its unique resource name e.g. default
			Id:        v.Name, // The name is the ID, the unique resource name e.g. apps/linear-cinema-360910/services/default
			ProjectId: projectId,
		})
	}

	utility.BatchUpsertOrElseSingle(db, serviceDBs)
	WriteServicesAuth(db, user, serviceDBs)
}

func WriteServicesAuth(db *gorm.DB, user types.User, gaeServiceDBs []gaetypes.GAEServiceDB) {
	if len(gaeServiceDBs) == 0 {
		return
	}

	gaeServiceAuths := make([]gaetypes.GAEServiceAuth, 0, len(gaeServiceDBs))
	for _, v := range gaeServiceDBs {
		gaeServiceAuths = append(gaeServiceAuths, gaetypes.GAEServiceAuth{Gmail: user.Gmail.String, Id: v.Id})
	}

	utility.BatchUpsertOrElseSingle(db, gaeServiceAuths)
}

func ReadVersionsCache(db *gorm.DB, user types.User, projectId string, serviceName string) (*gaetypes.GAEListVersionsResponse, error) {
	var gaeVersionDBs []gaetypes.GAEVersionDB
	result := db.Joins(
		" INNER JOIN gae_version_auths"+
			" ON gae_version_auths.id = gae_version_dbs.id"+
			" AND gae_version_auths.gmail = ?"+
			" AND gae_version_dbs.service_id = ?", user.Gmail.String, gaetypes.ToServiceId(projectId, serviceName),
	).Order("gae_version_dbs.name, gae_version_dbs.id").Find(&gaeVersionDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(gaeVersionDBs) == 0 {
		fmt.Printf("Cache miss\n")
		return nil, fmt.Errorf("Cache miss")
	}
	fmt.Printf("Cache hit\n")

	gaeVersions := make([]gaetypes.GAEVersion, 0, len(gaeVersionDBs))
	for _, v := range gaeVersionDBs {
		gaeVersions = append(gaeVersions, gaetypes.GAEVersion{
			// GAE API is misleading
			Id:            v.Name,
			Name:          v.Id,
			ServingStatus: v.ServingStatus,
		})
	}

	return &gaetypes.GAEListVersionsResponse{Versions: gaeVersions}, nil
}

func WriteVersionsCache(db *gorm.DB, user types.User, projectId string, serviceName string, resp *gaetypes.GAEListVersionsResponse) {
	if resp == nil {
		panic("Unexpected list of versions")
	}
	if len(resp.Versions) == 0 {
		return
	}

	gaeVersionDBs := make([]gaetypes.GAEVersionDB, 0, len(resp.Versions))
	for _, v := range resp.Versions {
		gaeVersionDBs = append(gaeVersionDBs, gaetypes.GAEVersionDB{
			// GAE API is misleading
			Name:          v.Id,   // The ID is the name, not its unique resource name e.g. 20220830t021415
			Id:            v.Name, // The name is the ID, the unique resource name e.g. apps/linear-cinema-360910/services/default/versions/20220830t021415
			ProjectId:     projectId,
			ServiceName:   serviceName,
			ServiceId:     gaetypes.ToServiceId(projectId, serviceName),
			ServingStatus: v.ServingStatus,
		})
	}

	utility.BatchUpsertOrElseSingle(db, gaeVersionDBs)
	WriteVersionsAuth(db, user, gaeVersionDBs)
}

func WriteVersionsAuth(db *gorm.DB, user types.User, gaeVersionDBs []gaetypes.GAEVersionDB) {
	if len(gaeVersionDBs) == 0 {
		return
	}

	gaeVersionAuths := make([]gaetypes.GAEVersionAuth, 0, len(gaeVersionDBs))
	for _, v := range gaeVersionDBs {
		gaeVersionAuths = append(gaeVersionAuths, gaetypes.GAEVersionAuth{Gmail: user.Gmail.String, Id: v.Id})
	}

	utility.BatchUpsertOrElseSingle(db, gaeVersionAuths)
}

func ReadInstancesCache(db *gorm.DB, user types.User, projectId string, serviceName string, versionName string) (*gaetypes.GAEListInstancesResponse, error) {
	var instanceDBs []gaetypes.GAEInstanceDB
	result := db.Joins(
		" INNER JOIN gae_instance_auths"+
			" ON gae_instance_auths.id = gae_instance_dbs.id"+
			" AND gae_instance_auths.gmail = ?"+
			" AND gae_instance_dbs.version_id = ?", user.Gmail.String, gaetypes.ToVersionId(projectId, serviceName, versionName),
	).Order("gae_instance_dbs.name, gae_instance_dbs.id").Find(&instanceDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(instanceDBs) == 0 {
		fmt.Printf("Cache miss\n")
		return nil, fmt.Errorf("Cache miss")
	}
	fmt.Printf("Cache hit\n")

	instances := make([]gaetypes.GAEInstance, 0, len(instanceDBs))
	for _, instanceDB := range instanceDBs {
		instances = append(instances, gaetypes.GAEInstance{
			// GAE API is misleading
			Name:   instanceDB.Id,
			Id:     instanceDB.Name,
			VMName: instanceDB.VMName,
		})
	}

	return &gaetypes.GAEListInstancesResponse{Instances: instances}, nil
}

func WriteInstancesCache(db *gorm.DB, user types.User, projectId string, serviceName string, versionName string, resp *gaetypes.GAEListInstancesResponse) {
	if resp == nil {
		panic("Unexpected list of instances")
	}
	if len(resp.Instances) == 0 {
		return
	}

	gaeInstanceDBs := make([]gaetypes.GAEInstanceDB, 0, len(resp.Instances))
	for _, v := range resp.Instances {
		gaeInstanceDBs = append(gaeInstanceDBs, gaetypes.GAEInstanceDB{
			// GAE API is misleading
			Name:        v.Id,
			Id:          v.Name,
			ProjectId:   projectId,
			ServiceName: serviceName,
			VersionName: versionName,
			VersionId:   gaetypes.ToVersionId(projectId, serviceName, versionName),
			VMName:      v.VMName,
		})
	}

	utility.BatchUpsertOrElseSingle(db, gaeInstanceDBs)
	WriteInstancesAuth(db, user, gaeInstanceDBs)
}

func WriteInstancesAuth(db *gorm.DB, user types.User, gaeInstanceDBs []gaetypes.GAEInstanceDB) {
	if len(gaeInstanceDBs) == 0 {
		return
	}

	gaeInstanceAuths := make([]gaetypes.GAEInstanceAuth, 0, len(gaeInstanceDBs))
	for _, v := range gaeInstanceDBs {
		gaeInstanceAuths = append(gaeInstanceAuths, gaetypes.GAEInstanceAuth{Gmail: user.Gmail.String, Id: v.Id})
	}

	utility.BatchUpsertOrElseSingle(db, gaeInstanceAuths)
}
