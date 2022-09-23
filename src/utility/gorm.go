package utility

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: generics
func BatchUpsertOrElseSingle[T any](db *gorm.DB, rows []T) {
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
