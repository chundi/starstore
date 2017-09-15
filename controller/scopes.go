package controller

import "github.com/jinzhu/gorm"

func BaseAvailable(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_date IS NULL")
}

func BasePublished(db *gorm.DB) *gorm.DB {
	return db.Where("publish_date IS NOT NULL")
}
