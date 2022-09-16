package types

import (
	"fmt"

	"gorm.io/gorm"
)

type GCEInstanceDB struct {
	gorm.Model
	GCEInstance
	Id        string `gorm:"primaryKey:idx"`
	ProjectId string
	Zone      string
}

func (in *GCEInstanceDB) ToGCEInstance() GCEInstance {
	return GCEInstance{
		Id:     in.Id,
		Name:   in.Name,
		Status: in.Status,
	}
}

func (in *GCEInstanceDB) ToSearchString() string {
	return fmt.Sprintf("Type: Instance, Name: %v, Id: %v, Project Id: %v, Zone: %v", in.Name, in.Id, in.ProjectId, in.Zone)
}

func (in *GCEInstanceDB) ToLink() string {
	return fmt.Sprintf("https://console.cloud.google.com/compute/instancesDetail/%v/instances/%v?project=%v", in.Zone, in.Name, in.ProjectId)
}

func GCEInstanceDBToGCEAggregatedInstances(in []GCEInstanceDB) GCEAggregatedInstances {
	zoneMetadataMap := make(map[string][]GCEInstance)
	for _, v := range in {
		zoneMetadataMap[v.Zone] = append(zoneMetadataMap[v.Zone], v.ToGCEInstance())
	}

	var out GCEAggregatedInstances
	for zone, instances := range zoneMetadataMap {
		out.Zones = append(out.Zones, ZoneMetadata{Zone: zone, Instances: instances})
	}

	return out
}

func GCEAggregatedInstancesToGCEInstanceDB(user User, projectId string, in GCEAggregatedInstances) []GCEInstanceDB {
	var out []GCEInstanceDB
	for _, zone := range in.Zones {
		for _, instance := range zone.Instances {
			var gceInstanceDB GCEInstanceDB
			gceInstanceDB.Id = instance.Id
			gceInstanceDB.ProjectId = projectId
			gceInstanceDB.Zone = zone.Zone
			gceInstanceDB.Name = instance.Name
			gceInstanceDB.Status = instance.Status

			out = append(out, gceInstanceDB)
		}
	}
	return out
}
