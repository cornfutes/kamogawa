package gcpcache

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"kamogawa/types"
	"kamogawa/types/gcp/coretypes"
)

func ReadGCPProjectAPIsCache(db *gorm.DB, user types.User, projectDB coretypes.ProjectDB) []coretypes.GCPProjectAPI {
	var gcpProjectAPIs []coretypes.GCPProjectAPI

	result := db.Joins(
		" INNER JOIN project_auths"+
			" ON project_auths.project_id = gcp_project_apis.project_id"+
			" AND project_auths.gmail = ?"+
			" AND gcp_project_apis.project_id = ?", user.Gmail.String, projectDB.ProjectId,
	).Find(&gcpProjectAPIs)
	if result.Error != nil {
		panic("Query failed")
	}

	if len(gcpProjectAPIs) == 0 {
		fmt.Printf("Cache miss\n")
	}
	fmt.Printf("Cache hit\n")

	return gcpProjectAPIs
}

func WriteGCPProjectAPIsCache(db *gorm.DB, gcpProjectAPIs []coretypes.GCPProjectAPI) {
	if len(gcpProjectAPIs) == 0 {
		return
	}
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&gcpProjectAPIs)
}
