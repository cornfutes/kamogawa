package core

import (
	"kamogawa/types"
	"kamogawa/types/gcp/coretypes"
	gaetypes "kamogawa/types/gcp/gaetypes"
	"kamogawa/types/gcp/gcetypes"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(
		&types.User{},
		&coretypes.ProjectAuth{},
		&coretypes.ProjectDB{},
		&coretypes.GCPProjectAPI{},
		&gcetypes.GCEInstanceAuth{},
		&gcetypes.GCEInstanceDB{},
		&gaetypes.GAEServiceAuth{},
		&gaetypes.GAEServiceDB{},
		&gaetypes.GAEVersionAuth{},
		&gaetypes.GAEVersionDB{},
		&gaetypes.GAEInstanceAuth{},
		&gaetypes.GAEInstanceDB{},
	)

	return db
}
