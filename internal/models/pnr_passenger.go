package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type PassengerType string

const (
	PassengerADT PassengerType = "ADT" // Adult
	PassengerCHD PassengerType = "CHD" // Child
	PassengerINF PassengerType = "INF" // Infant
)

type PNRPassenger struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PNRID          *uuid.UUID    `gorm:"type:uuid;index"                                  json:"pnr_id"`
	FirstName      string        `gorm:"type:varchar(100)"                                json:"first_name"      validate:"required,max=100"`
	LastName       string        `gorm:"type:varchar(100)"                                json:"last_name"       validate:"required,max=100"`
	PassengerType  PassengerType `gorm:"type:varchar(10)"                                 json:"passenger_type"  validate:"oneof=ADT CHD INF"`
	BirthDate      *time.Time    `gorm:"type:date"                                        json:"birth_date"`
	PassportNumber string        `gorm:"type:varchar(50)"                                 json:"passport_number" validate:"max=50"`
	PNR            *PNR          `gorm:"foreignKey:PNRID"                                 json:"pnr,omitempty"`
	SeatAssignment *SeatAssignment `gorm:"foreignKey:PassengerID"                         json:"seat_assignment,omitempty"`
	Ticket         *Ticket       `gorm:"foreignKey:PassengerID"                           json:"ticket,omitempty"`
	SSRs           []PassengerSSR  `gorm:"foreignKey:PassengerID"                         json:"ssrs,omitempty"`
	Meals          []PassengerMeal `gorm:"foreignKey:PassengerID"                         json:"meals,omitempty"`
	Baggage        []Baggage     `gorm:"foreignKey:PassengerID"                           json:"baggage,omitempty"`
	Checkins       []Checkin     `gorm:"foreignKey:PassengerID"                           json:"checkins,omitempty"`
	BoardingPasses []BoardingPass `gorm:"foreignKey:PassengerID"                          json:"boarding_passes,omitempty"`
}

func (p *PNRPassenger) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
