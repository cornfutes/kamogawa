package gaetypes

import "fmt"

type GAEServiceDB struct {
	Name      string `json:"name"`
	Id        string `json:"id,gorm:"primaryKey"`
	ProjectId string `json:"id,gorm:"primaryKey"`
}

type GAEVersionDB struct {
	Name          string `json:"name"`
	Id            string `json:"id",gorm:"primaryKey"`
	ServingStatus string `json:"servingStatus"`
	ParentId      string `json:"id,gorm:"primaryKey"`
}

type GAEInstanceDB struct {
	Name      string `json:"name"`
	Id        string `json:"id,gorm:"primaryKey""`
	VMName    string `json:"vmName"`
	VersionId string `json:"id,gorm:"primaryKey"`
}

func (in *GAEServiceDB) ToSearchString() string {
	return fmt.Sprintf("GAE service by the name %v", in.Id)
}

func (in *GAEServiceDB) ToLink() string {
	return fmt.Sprintf("google.com/gae/service/%v", in.Id)
}

func (in *GAEVersionDB) ToSearchString() string {
	return fmt.Sprintf("GAE version by the name %v", in.Id)
}

func (in *GAEVersionDB) ToLink() string {
	return fmt.Sprintf("google.com/gae/version/%v", in.Id)
}
