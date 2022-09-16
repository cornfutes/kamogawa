package types

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"primaryKey"`
	Password     string
	Gmail        *sql.NullString
	Scope        *sql.NullString
	AccessToken  *sql.NullString
	RefreshToken *sql.NullString
}

type ProjectDB struct {
	gorm.Model
	Project
	Email     string `gorm:"primaryKey:idx"`
	ProjectId string `gorm:"primaryKey:idx"`
	//Parent         ProjectParent
}

func (in *ProjectDB) ToProject() Project {
	return Project{
		ProjectNumber:  in.ProjectNumber,
		ProjectId:      in.ProjectId,
		LifeCycleState: in.LifeCycleState,
		Name:           in.Name,
		CreateTime:     in.CreateTime,
	}
}

func (in *Project) ToProjectDB(email string) ProjectDB {
	var out ProjectDB
	out.Email = email
	out.ProjectNumber = in.ProjectNumber
	out.ProjectId = in.ProjectId
	out.LifeCycleState = in.LifeCycleState
	out.Name = in.Name
	out.CreateTime = in.CreateTime
	return out
}

type GCEInstanceDB struct {
	gorm.Model
	GCEInstance
	Email     string `gorm:"primaryKey:idx"`
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
			gceInstanceDB.Email = user.Email
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
