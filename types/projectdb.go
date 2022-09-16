package types

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ProjectDB struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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

func (in *ProjectDB) ToSearchString() string {
	return fmt.Sprintf("Type: Project, Name: %v, Id: %v, Number: %v", in.Name, in.ProjectId, in.ProjectNumber)
}

func (in *ProjectDB) ToLink() string {
	return fmt.Sprintf("https://console.cloud.google.com/welcome?project=%v", in.ProjectId)
}
