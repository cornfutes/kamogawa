package coretypes

import (
	"bytes"
	"fmt"
	"kamogawa/types/gcp/gcetypes"
	"time"

	"gorm.io/datatypes"
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
	ProjectId string `gorm:"primaryKey:idx"`
	//Parent         ProjectParent
}

type GCPProjectAPI struct {
	ProjectId string `gorm:"primaryKey:idx"`
	API       datatypes.JSON
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

func ProjectToProjectDB(in *gcetypes.Project) ProjectDB {
	var out ProjectDB
	out.ProjectNumber = in.ProjectNumber
	out.ProjectId = in.ProjectId
	out.LifeCycleState = in.LifeCycleState
	out.Name = in.Name
	out.CreateTime = in.CreateTime
	return out
}

func (in *GCPProjectAPI) IsEnabled(title string) bool {
	return bytes.Contains(in.API, []byte(title))
}
