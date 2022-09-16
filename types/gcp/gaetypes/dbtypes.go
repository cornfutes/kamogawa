package gaetypes

type GAEServiceDB struct {
	Name      string `json:"name"`
	Id        string `json:"id,gorm:"primaryKey"`
	ProjectId string `json:"id,gorm:"primaryKey"`
}

type GAEVersionDB struct {
	Name          string `json:"name"`
	Id            string `json:"id",gorm:"primaryKey"`
	ServingStatus string `json:"servingStatus"`
	ServiceId     string `json:"id,gorm:"primaryKey"`
}

type GAEInstanceDB struct {
	Name      string `json:"name"`
	Id        string `json:"id,gorm:"primaryKey""`
	VMName    string `json:"vmName"`
	VersionId string `json:"id,gorm:"primaryKey"`
}
