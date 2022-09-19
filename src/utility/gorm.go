package utility

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"kamogawa/types/gcp/coretypes"
	"kamogawa/types/gcp/gaetypes"
	"kamogawa/types/gcp/gcetypes"
)

// TODO: generics
func BatchUpsertOrElseSingleGAEServiceAuth(db *gorm.DB, rows []gaetypes.GAEServiceAuth) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleGAEServiceDB(db *gorm.DB, rows []gaetypes.GAEServiceDB) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleGAEVersionAuth(db *gorm.DB, rows []gaetypes.GAEVersionAuth) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleGAEVersionDB(db *gorm.DB, rows []gaetypes.GAEVersionDB) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleGAEInstanceAuth(db *gorm.DB, rows []gaetypes.GAEInstanceAuth) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleGAEInstanceDB(db *gorm.DB, rows []gaetypes.GAEInstanceDB) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleGCEInstanceAuth(db *gorm.DB, rows []gcetypes.GCEInstanceAuth) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleGCEInstanceDB(db *gorm.DB, rows []gcetypes.GCEInstanceDB) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleGCPProjectAPI(db *gorm.DB, rows []coretypes.GCPProjectAPI) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleProjectAuth(db *gorm.DB, rows []coretypes.ProjectAuth) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}
func BatchUpsertOrElseSingleProjectDB(db *gorm.DB, rows []coretypes.ProjectDB) {
	if result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rows); result.Error != nil {
		for _, row := range rows {
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&row)
		}
	}
}