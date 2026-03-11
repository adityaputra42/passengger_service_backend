package db

import (
	"fmt"
	"log"
	"passenger_service_backend/internal/config"
	"passenger_service_backend/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: false,
		PrepareStmt:                              true, // cache prepared statements → faster repeated queries
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Connection pool tuning
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	log.Println("Database connected successfully!")
	return nil
}

// ─────────────────────────────────────────────
// AutoMigrate
// ─────────────────────────────────────────────

func Migrate() error {
	// Order matters: tables with no FK dependencies go first
	err := DB.AutoMigrate(
		&models.User{},        // no FK
		&models.Aircraft{},    // no FK
		&models.Airport{},     // no FK
		&models.Passenger{},   // no FK (maps to "Passengger" table via TableName())
		&models.Seat{},        // FK → Aircraft
		&models.Flight{},      // FK → Aircraft, Airport x2
		&models.FlightSeat{},  // FK → Flight, Seat
		&models.Booking{},     // FK → User, Flight
		&models.BookingPassenger{}, // FK → Booking, Passenger, FlightSeat
		&models.Payment{},     // FK → Booking
	)
	if err != nil {
		return fmt.Errorf("failed to run AutoMigrate: %w", err)
	}
	log.Println("Database migration completed successfully")

	if err := createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// ─────────────────────────────────────────────
// Indexes
// ─────────────────────────────────────────────

func createIndexes() error {
	indexes := []struct {
		name string
		sql  string
	}{
		// ── User ──────────────────────────────────────────────
		{
			"idx_user_email",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_user_email ON "User"(email)`,
		},
		{
			"idx_user_role_created",
			`CREATE INDEX IF NOT EXISTS idx_user_role_created ON "User"(role, created_at)`,
		},

		// ── Airport ───────────────────────────────────────────
		{
			"idx_airport_code",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_airport_code ON "Airport"(code)`,
		},
		{
			"idx_airport_country_city",
			`CREATE INDEX IF NOT EXISTS idx_airport_country_city ON "Airport"(country, city)`,
		},

		// ── Aircraft ──────────────────────────────────────────
		{
			"idx_aircraft_registration",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_aircraft_registration ON "Aircraft"(registration_no)`,
		},

		// ── Seat ──────────────────────────────────────────────
		{
			"idx_seat_aircraft_class",
			`CREATE INDEX IF NOT EXISTS idx_seat_aircraft_class ON "Seat"(aircraft_id, seat_class)`,
		},
		{
			"idx_seat_aircraft_number",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_seat_aircraft_number ON "Seat"(aircraft_id, seat_number)`,
		},

		// ── Flight ────────────────────────────────────────────
		{
			"idx_flight_number",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_flight_number ON "Flight"(flight_number)`,
		},
		{
			"idx_flight_dep_arr_time",
			`CREATE INDEX IF NOT EXISTS idx_flight_dep_arr_time ON "Flight"(departure_airport_id, arrival_airport_id, departure_time)`,
		},
		{
			"idx_flight_status_departure",
			`CREATE INDEX IF NOT EXISTS idx_flight_status_departure ON "Flight"(status, departure_time)`,
		},
		{
			"idx_flight_aircraft_status",
			`CREATE INDEX IF NOT EXISTS idx_flight_aircraft_status ON "Flight"(aircraft_id, status)`,
		},

		// ── FlightSeat ────────────────────────────────────────
		{
			"idx_flightseat_flight_status",
			// Most critical: searching available seats for a flight
			`CREATE INDEX IF NOT EXISTS idx_flightseat_flight_status ON "FlightSeat"(flight_id, status)`,
		},
		{
			"idx_flightseat_flight_seat",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_flightseat_flight_seat ON "FlightSeat"(flight_id, seat_id)`,
		},

		// ── Booking ───────────────────────────────────────────
		{
			"idx_booking_code",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_booking_code ON "booking"(booking_code)`,
		},
		{
			"idx_booking_user_status",
			`CREATE INDEX IF NOT EXISTS idx_booking_user_status ON "booking"(user_id, status)`,
		},
		{
			"idx_booking_flight_status",
			`CREATE INDEX IF NOT EXISTS idx_booking_flight_status ON "booking"(flight_id, status)`,
		},
		{
			"idx_booking_status_created",
			`CREATE INDEX IF NOT EXISTS idx_booking_status_created ON "booking"(status, created_at)`,
		},

		// ── BookingPassenger ──────────────────────────────────
		{
			"idx_bookingpax_booking",
			`CREATE INDEX IF NOT EXISTS idx_bookingpax_booking ON "Booking Passengger"(booking_id)`,
		},
		{
			"idx_bookingpax_passenger",
			`CREATE INDEX IF NOT EXISTS idx_bookingpax_passenger ON "Booking Passengger"(passenger_id)`,
		},
		{
			"idx_bookingpax_flightseat",
			// Prevent double-booking the same seat in the same booking
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_bookingpax_flightseat ON "Booking Passengger"(flight_seat_id)`,
		},

		// ── Passenger ─────────────────────────────────────────
		{
			"idx_passenger_passport",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_passenger_passport ON "Passengger"(passport_number)`,
		},

		// ── Payment ───────────────────────────────────────────
		{
			"idx_payment_booking_status",
			`CREATE INDEX IF NOT EXISTS idx_payment_booking_status ON "payment"(booking_id, status)`,
		},
		{
			"idx_payment_status_paid",
			`CREATE INDEX IF NOT EXISTS idx_payment_status_paid ON "payment"(status, paid_at)`,
		},
		{
			"idx_payment_method_status",
			`CREATE INDEX IF NOT EXISTS idx_payment_method_status ON "payment"(method, status)`,
		},
	}

	errorCount := 0
	for _, idx := range indexes {
		if err := DB.Exec(idx.sql).Error; err != nil {
			log.Printf("  [WARN] index %q: %v", idx.name, err)
			errorCount++
		} else {
			log.Printf("  [OK]   index: %s", idx.name)
		}
	}

	if errorCount > 0 {
		log.Printf("Indexes done with %d warning(s)", errorCount)
	} else {
		log.Println("All indexes created/verified successfully")
	}
	return nil
}

// ─────────────────────────────────────────────
// Utilities
// ─────────────────────────────────────────────

func GetDBStats() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting DB stats: %v", err)
		return
	}
	s := sqlDB.Stats()
	log.Printf(
		"DB Stats — Open: %d | In Use: %d | Idle: %d | WaitCount: %d | WaitDuration: %v | MaxIdleClosed: %d | MaxLifetimeClosed: %d",
		s.OpenConnections, s.InUse, s.Idle,
		s.WaitCount, s.WaitDuration,
		s.MaxIdleClosed, s.MaxLifetimeClosed,
	)
}

func GetDB() *gorm.DB {
	return DB
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error closing DB: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing DB connection: %v", err)
	}
}
