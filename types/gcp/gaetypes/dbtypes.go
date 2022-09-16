package gaetypes

import (
	"fmt"
	"strings"
)

type GAEServiceDB struct {
	Name      string
	Id        string `gorm:"primaryKey"`
	ProjectId string `gorm:"primaryKey"`
}

type GAEVersionDB struct {
	Name          string
	Id            string `gorm:"primaryKey"`
	ServingStatus string
	ParentId      string `gorm:"primaryKey"`
}

type GAEInstanceDB struct {
	Name      string
	Id        string `gorm:"primaryKey"`
	VMName    string
	VersionId string `gorm:"primaryKey"`
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
	var parentPath = strings.Split(in.ParentId, "/services/")
	var projectId = strings.Split(parentPath[0], "apps/")[1]
	var serviceId = parentPath[1]
	return fmt.Sprintf("https://console.cloud.google.com/appengine/versions?serviceId=%v&project=%v", serviceId, projectId)
}
