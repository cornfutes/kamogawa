package types

import (
	"database/sql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"primaryKey"`
	Password     string
	Gmail        *sql.NullString
	Scope        *sql.NullString
	AccessToken  *sql.NullString
	RefreshToken *sql.NullString
}
