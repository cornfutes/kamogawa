package gaecache

import (
	"fmt"
	"kamogawa/types/gcp/gaetypes"
	"strings"

	"gorm.io/gorm"
)

func ReadServicesCache(db *gorm.DB, projectId string) (*gaetypes.GAEListServicesResponse, error) {
	var serviceDBs []gaetypes.GAEServiceDB
	result := db.Where("project_id = ?", projectId).Find(&serviceDBs)
	if result.Error != nil {
		fmt.Printf("Query failed\n")
		return nil, fmt.Errorf("Query failed")
	}

	if len(serviceDBs) == 0 {
		fmt.Printf("Cache miss\n")
		return nil, fmt.Errorf("Cache miss")
	}
	fmt.Printf("Cache hit\n")

	services := make([]gaetypes.GAEService, 0, len(serviceDBs))
	for _, serviceDB := range serviceDBs {
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

func WriteServicesCache(db *gorm.DB, projectId string, resp *gaetypes.GAEListServicesResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Services) == 0 {
		return
	}

	serviceDBs := make([]gaetypes.GAEServiceDB, 0, len(resp.Services))

	for _, v := range resp.Services {
		serviceDBs = append(serviceDBs, gaetypes.GAEServiceDB{
			// GAE API is misleading. The ID is the name of the service, not its unique resource path.
			Name: v.Id,
			// The name of the service is the unique resource name
			Id:        v.Name,
			ProjectId: projectId,
		})
	}

	for _, serviceDB := range serviceDBs {
		db.FirstOrCreate(&serviceDB)
	}
}

func ReadVersionsCache(db *gorm.DB, projectId string, serviceId string) (*gaetypes.GAEListVersionsResponse, error) {
	var versionDBs []gaetypes.GAEVersionDB
	result := db.Where("service_id = ?", fmt.Sprintf("apps/%v/services/%v", projectId, serviceId)).Find(&versionDBs)
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
			Id:            versionDB.Name,
			Name:          versionDB.Id,
			ServingStatus: versionDB.ServingStatus,
		})
	}

	return &gaetypes.GAEListVersionsResponse{Versions: versions}, nil
}

func WriteVersionsCache(db *gorm.DB, projectId string, serviceId string, resp *gaetypes.GAEListVersionsResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Versions) == 0 {
		return
	}

	versionDBs := make([]gaetypes.GAEVersionDB, 0, len(resp.Versions))

	for _, v := range resp.Versions {
		versionDBs = append(versionDBs, gaetypes.GAEVersionDB{
			Name:          v.Id,
			Id:            v.Name,
			ServingStatus: v.ServingStatus,
			ServiceId:     strings.Split(v.Name, "/versions/")[0],
		})
	}

	for _, versionDB := range versionDBs {
		db.FirstOrCreate(&versionDB)
	}
}

func ReadInstancesCache(db *gorm.DB, versionId string) (*gaetypes.GAEListInstancesResponse, error) {
	var instanceDBs []gaetypes.GAEInstanceDB
	result := db.Where("version_id = ?", versionId).Find(&instanceDBs)
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
			Name:   instanceDB.Id,
			Id:     instanceDB.Name,
			VMName: instanceDB.VMName})
	}

	return &gaetypes.GAEListInstancesResponse{Instances: instances}, nil
}

func WriteInstancesCache(db *gorm.DB, versionId string, resp *gaetypes.GAEListInstancesResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Instances) == 0 {
		return
	}

	instanceDBs := make([]gaetypes.GAEInstanceDB, 0, len(resp.Instances))

	for _, v := range resp.Instances {
		instanceDBs = append(instanceDBs, gaetypes.GAEInstanceDB{Name: v.Id, Id: v.Name, VMName: v.VMName, VersionId: versionId})
	}

	for _, instanceDB := range instanceDBs {
		db.FirstOrCreate(&instanceDB)
	}
}
