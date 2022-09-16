package core

import (
	"kamogawa/types"
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
		&gcetypes.ProjectAuth{},
		&gcetypes.ProjectDB{},
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
