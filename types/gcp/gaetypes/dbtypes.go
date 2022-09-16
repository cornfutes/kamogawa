package gaetypes

import "fmt"

type GAEServiceDB struct {
	Name      string
	Id        string `gorm:"primaryKey"`
	ProjectId string `gorm:"primaryKey"`
}

type GAEVersionDB struct {
	Name          string
	Id            string `gorm:"primaryKey"`
	ServingStatus string
	ServiceId     string `gorm:"primaryKey"`
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
	return fmt.Sprintf("google.com/gae/service/%v", in.Id)
}

func (in *GAEVersionDB) ToSearchString() string {
	return fmt.Sprintf("GAE version by the name %v", in.Id)
}

func (in *GAEVersionDB) ToLink() string {
	return fmt.Sprintf("google.com/gae/version/%v", in.Id)
}
