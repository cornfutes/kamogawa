package coretypes

import (
	"fmt"
	"kamogawa/types/gcp/gcetypes"
	"time"

	"gorm.io/gorm"
)

type ProjectAuth struct {
	Gmail     string `gorm:"primaryKey:idx"`
	ProjectId string `gorm:"primaryKey:idx"`
}

type ProjectDB struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	gcetypes.Project
	ProjectId     string `gorm:"primaryKey:idx"`
	HasGCEEnabled bool
	//Parent         ProjectParent
}

func (in *ProjectDB) ToProject() gcetypes.Project {
	return gcetypes.Project{
		ProjectNumber:  in.ProjectNumber,
		ProjectId:      in.ProjectId,
		LifeCycleState: in.LifeCycleState,
		Name:           in.Name,
		CreateTime:     in.CreateTime,
	}
}

func (in *ProjectDB) ToSearchString() string {
	return fmt.Sprintf("Type: Project, Name: %v, Id: %v, Number: %v", in.Name, in.ProjectId, in.ProjectNumber)
}

func (in *ProjectDB) ToLink() string {
	return fmt.Sprintf("https://console.cloud.google.com/welcome?project=%v", in.ProjectId)
}

func ProjectToProjectDB(in *gcetypes.Project, hasGCEEnabled bool) ProjectDB {
	var out ProjectDB
	out.ProjectNumber = in.ProjectNumber
	out.ProjectId = in.ProjectId
	out.LifeCycleState = in.LifeCycleState
	out.Name = in.Name
	out.CreateTime = in.CreateTime
	out.HasGCEEnabled = hasGCEEnabled
	return out
}
