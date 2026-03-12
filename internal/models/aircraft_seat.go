package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AircraftSeat struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	AircraftID   uuid.UUID  `gorm:"type:uuid;not null;index"                         json:"aircraft_id"    validate:"required"`
	SeatNumber   string     `gorm:"type:varchar(5);not null"                         json:"seat_number"    validate:"required,max=5"`
	RowNumber    int        `gorm:"not null"                                         json:"row_number"     validate:"required,min=1"`
	SeatLetter   string     `gorm:"type:char(1)"                                     json:"seat_letter"    validate:"max=1"`
	XPosition    int        `gorm:"default:0"                                        json:"x_position"`
	YPosition    int        `gorm:"default:0"                                        json:"y_position"`
	SeatClassID  *uuid.UUID `gorm:"type:uuid"                                        json:"seat_class_id"`
	SeatType     string     `gorm:"type:varchar(20)"                                 json:"seat_type"      validate:"max=20"`
	IsExitRow    bool       `gorm:"default:false"                                    json:"is_exit_row"`
	Aircraft     Aircraft   `gorm:"foreignKey:AircraftID"                            json:"aircraft,omitempty"`
	SeatClass    *SeatClass `gorm:"foreignKey:SeatClassID"                           json:"seat_class,omitempty"`
}

func (a *AircraftSeat) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
