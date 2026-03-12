package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Aircraft struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Model        string         `gorm:"type:varchar(100)"                                json:"model"         validate:"max=100"`
	Manufacturer string         `gorm:"type:varchar(100)"                                json:"manufacturer"  validate:"max=100"`
	TotalSeats   int            `gorm:"default:0"                                        json:"total_seats"   validate:"min=0"`
	Seats        []AircraftSeat `gorm:"foreignKey:AircraftID"                            json:"seats,omitempty"`
}

func (a *Aircraft) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
