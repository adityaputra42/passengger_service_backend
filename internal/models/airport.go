package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)



type Airport struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code     string    `gorm:"type:varchar(3);uniqueIndex;not null"             json:"code"      validate:"required,len=3"`
	Name     string    `gorm:"type:varchar(255)"                                json:"name"      validate:"max=255"`
	City     string    `gorm:"type:varchar(255)"                                json:"city"      validate:"max=255"`
	Country  string    `gorm:"type:varchar(255)"                                json:"country"   validate:"max=255"`
	Timezone string    `gorm:"type:varchar(100)"                                json:"timezone"  validate:"max=100"`
}

func (a *Airport) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
