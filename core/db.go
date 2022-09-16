package core

import (
	"kamogawa/types"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&types.User{}, &types.ProjectDB{}, &types.GCEInstanceDB{})

	return db
}
