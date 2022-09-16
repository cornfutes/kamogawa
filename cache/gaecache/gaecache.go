package gaecache

import (
	"fmt"
	"gorm.io/gorm"
	"kamogawa/types/gcp/gaetypes"
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
		services = append(services, gaetypes.GAEService{Name: serviceDB.Name, Id: serviceDB.Id, Split: gaetypes.GAEServiceTrafficAllocations{}})
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
		serviceDBs = append(serviceDBs, gaetypes.GAEServiceDB{Name: v.Name, Id: v.Id, ProjectId: projectId})
	}

	for _, serviceDB := range serviceDBs {
		db.FirstOrCreate(&serviceDB)
	}
}

func ReadVersionsCache(db *gorm.DB, serviceId string) (*gaetypes.GAEListVersionsResponse, error) {
	var versionDBs []gaetypes.GAEVersionDB
	result := db.Where("service_id = ?", serviceId).Find(&versionDBs)
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
		versions = append(versions, gaetypes.GAEVersion{Name: versionDB.Name, Id: versionDB.Id, ServingStatus: versionDB.ServingStatus})
	}

	return &gaetypes.GAEListVersionsResponse{Versions: versions}, nil
}

func WriteVersionsCache(db *gorm.DB, serviceId string, resp *gaetypes.GAEListVersionsResponse) {
	if resp == nil {
		panic("Unexpected list of projects")
	}
	if len(resp.Versions) == 0 {
		return
	}

	versionDBs := make([]gaetypes.GAEVersionDB, 0, len(resp.Versions))

	for _, v := range resp.Versions {
		versionDBs = append(versionDBs, gaetypes.GAEVersionDB{Name: v.Name, Id: v.Id, ServingStatus: v.ServingStatus, ServiceId: serviceId})
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
		instances = append(instances, gaetypes.GAEInstance{Name: instanceDB.Name, Id: instanceDB.Id, VMName: instanceDB.VMName})
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
		instanceDBs = append(instanceDBs, gaetypes.GAEInstanceDB{Name: v.Name, Id: v.Id, VMName: v.VMName, VersionId: versionId})
	}

	for _, instanceDB := range instanceDBs {
		db.FirstOrCreate(&instanceDB)
	}
}
