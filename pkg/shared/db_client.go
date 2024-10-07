package shared

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserSchema struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey"`
	Email    string    `gorm:"unique"`
	Password string
}

func (UserSchema) TableName() string {
	return "users"
}

func NewDBClient() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("db/test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&UserSchema{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
