package gaetypes

import (
	"fmt"
	"time"
)

type GAEServiceAuth struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Gmail     string `gorm:"primaryKey:idx"`
	Id        string `gorm:"primaryKey:idx"`
}

type GAEServiceDB struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Id        string `gorm:"primaryKey"`
	ProjectId string `gorm:"primaryKey"`
}

type GAEVersionAuth struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Gmail     string `gorm:"primaryKey:idx"`
	Id        string `gorm:"primaryKey:idx"`
}

type GAEVersionDB struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Id            string `gorm:"primaryKey"`
	ProjectId     string
	ServiceName   string
	ServiceId     string
	ServingStatus string
}

type GAEInstanceAuth struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Gmail     string `gorm:"primaryKey:idx"`
	Id        string `gorm:"primaryKey:idx"`
}

type GAEInstanceDB struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Id          string `gorm:"primaryKey"`
	ProjectId   string
	ServiceName string
	VersionName string
	VersionId   string
	VMName      string
}

func (in *GAEServiceDB) ToSearchString() string {
	return fmt.Sprintf("GAE service by the name %v", in.Id)
}

func (in *GAEServiceDB) ToLink() string {
	return fmt.Sprintf("https://console.cloud.google.com/appengine/services?project=%v", in.ProjectId)
}

func (in *GAEVersionDB) ToSearchString() string {
	return fmt.Sprintf("GAE version by the name %v", in.Id)
}

func (in *GAEVersionDB) ToLink() string {
	return fmt.Sprintf("https://console.cloud.google.com/appengine/versions?project=%v&serviceId=%v", in.ProjectId, in.ServiceName)
}

func ToServiceId(projectId string, serviceName string) string {
	return fmt.Sprintf("apps/%v/services/%v", projectId, serviceName)
}

func ToVersionId(projectId string, serviceName string, versionName string) string {
	return fmt.Sprintf("apps/%v/services/%v/versions/%v", projectId, serviceName, versionName)
}
