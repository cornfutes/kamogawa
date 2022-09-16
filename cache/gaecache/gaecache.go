package gaecache

import (
	"fmt"
	"gorm.io/gorm"
	"kamogawa/types"
	"kamogawa/types/gcp/gaetypes"
)

func ReadServicesCache(db *gorm.DB, user types.User, projectId string) (*gaetypes.GAEListServicesResponse, error) {
	var gaeServiceDBs []gaetypes.GAEServiceDB
	result := db.Joins(
		" INNER JOIN gae_service_auths"+
			" ON gae_service_auths.id = gae_service_dbs.id"+
			" AND gae_service_auths.gmail = ?"+
			" AND gae_service_dbs.project_id = ?", user.Gmail.String, projectId).Find(&gaeServiceDBs)
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

	for _, serviceDB := range serviceDBs {
		db.FirstOrCreate(&serviceDB)
	}

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
	for _, v := range gaeServiceAuths {
		db.FirstOrCreate(&v)
	}
}

func ReadVersionsCache(db *gorm.DB, projectId string, serviceName string) (*gaetypes.GAEListVersionsResponse, error) {
	var versionDBs []gaetypes.GAEVersionDB
	result := db.Where("service_id = ?", gaetypes.ToServiceId(projectId, serviceName)).Find(&versionDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(versionDBs) == 0 {
		fmt.Printf("Cache miss\n")
		return nil, fmt.Errorf("Cache miss")
	}
	fmt.Printf("Cache hit\n")

	versions := make([]gaetypes.GAEVersion, 0, len(versionDBs))
	for _, versionDB := range versionDBs {
		versions = append(versions, gaetypes.GAEVersion{
			// GAE API is misleading
			Id:            versionDB.Name,
			Name:          versionDB.Id,
			ServingStatus: versionDB.ServingStatus,
		})
	}

	return &gaetypes.GAEListVersionsResponse{Versions: versions}, nil
}

func WriteVersionsCache(db *gorm.DB, projectId string, serviceName string, resp *gaetypes.GAEListVersionsResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Versions) == 0 {
		return
	}

	versionDBs := make([]gaetypes.GAEVersionDB, 0, len(resp.Versions))

	for _, v := range resp.Versions {
		versionDBs = append(versionDBs, gaetypes.GAEVersionDB{
			// GAE API is misleading
			Name:          v.Id,   // The ID is the name, not its unique resource name e.g. 20220830t021415
			Id:            v.Name, // The name is the ID, the unique resource name e.g. apps/linear-cinema-360910/services/default/versions/20220830t021415
			ProjectId:     projectId,
			ServiceName:   serviceName,
			ServiceId:     gaetypes.ToServiceId(projectId, serviceName),
			ServingStatus: v.ServingStatus,
		})
	}

	for _, versionDB := range versionDBs {
		db.FirstOrCreate(&versionDB)
	}
}

func ReadInstancesCache(db *gorm.DB, projectId string, serviceName string, versionName string) (*gaetypes.GAEListInstancesResponse, error) {
	var instanceDBs []gaetypes.GAEInstanceDB
	result := db.Where("version_id = ?", gaetypes.ToVersionId(projectId, serviceName, versionName)).Find(&instanceDBs)
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

func WriteInstancesCache(db *gorm.DB, projectId string, serviceName string, versionName string, resp *gaetypes.GAEListInstancesResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Instances) == 0 {
		return
	}

	instanceDBs := make([]gaetypes.GAEInstanceDB, 0, len(resp.Instances))

	for _, v := range resp.Instances {
		instanceDBs = append(instanceDBs, gaetypes.GAEInstanceDB{
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

	for _, instanceDB := range instanceDBs {
		db.FirstOrCreate(&instanceDB)
	}
}
