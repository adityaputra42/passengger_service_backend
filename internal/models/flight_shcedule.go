package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlightSchedule struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	FlightNumber       string    `gorm:"type:varchar(10)"                                 json:"flight_number"         validate:"max=10"`
	DepartureAirportID uuid.UUID `gorm:"type:uuid;index"                                  json:"departure_airport_id"`
	ArrivalAirportID   uuid.UUID `gorm:"type:uuid"                                        json:"arrival_airport_id"`
	DepartureTime    string    `gorm:"type:time"                                        json:"departure_time"`
	ArrivalTime      string    `gorm:"type:time"                                        json:"arrival_time"`
	OperatingDays    string    `gorm:"type:varchar(20)"                                 json:"operating_days"`
	DepartureAirport Airport   `gorm:"foreignKey:DepartureAirportID"                    json:"departure_airport,omitempty"`
	ArrivalAirport   Airport   `gorm:"foreignKey:ArrivalAirportID"                      json:"arrival_airport,omitempty"`
	Flights          []Flight  `gorm:"foreignKey:ScheduleID"                            json:"flights,omitempty"`
}

func (f *FlightSchedule) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
