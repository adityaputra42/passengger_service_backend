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
		PrepareStmt:                              true,
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
// Urutan wajib dijaga: tabel tanpa FK → tabel dengan FK → many2many
// ─────────────────────────────────────────────

func Migrate() error {
	err := DB.AutoMigrate(
		// ── Tier 0: tidak ada FK sama sekali ──────────────────
		&models.Permission{}, // tabel "permissions", no FK
		&models.Airport{},    // tabel "Airport", no FK
		&models.Aircraft{},   // tabel "Aircraft", no FK
		&models.Passenger{},  // tabel "Passengger", no FK

		// ── Tier 1: FK ke Tier 0 ──────────────────────────────
		&models.Role{},       // tabel "roles", many2many → permissions (role_permissions junction)
		&models.Seat{},       // FK → Aircraft

		// ── Tier 2: FK ke Tier 1 ──────────────────────────────
		&models.User{},       // FK → roles (RoleID)
		&models.Flight{},     // FK → Aircraft, Airport x2

		// ── Tier 3: FK ke Tier 2 ──────────────────────────────
		&models.FlightSeat{}, // FK → Flight, Seat
		&models.Booking{},    // FK → User, Flight

		// ── Tier 4: FK ke Tier 3 ──────────────────────────────
		&models.BookingPassenger{}, // FK → Booking, Passenger, FlightSeat
		&models.Payment{},          // FK → Booking
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
		// ── Permission ────────────────────────────────────────
		{
			"idx_permission_name",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_permission_name ON permissions(name)`,
		},
		{
			"idx_permission_resource_action",
			// Composite: cari semua action dalam satu resource (e.g., "flights.*")
			`CREATE INDEX IF NOT EXISTS idx_permission_resource_action ON permissions(resource, action)`,
		},

		// ── Role ──────────────────────────────────────────────
		{
			"idx_role_name",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_role_name ON roles(name)`,
		},
		{
			"idx_role_level",
			// Filter/sort berdasarkan hierarki role
			`CREATE INDEX IF NOT EXISTS idx_role_level ON roles(level)`,
		},
		{
			"idx_role_system",
			`CREATE INDEX IF NOT EXISTS idx_role_system ON roles(is_system_role)`,
		},

		// ── role_permissions (junction table) ─────────────────
		{
			"idx_roleperm_role",
			`CREATE INDEX IF NOT EXISTS idx_roleperm_role ON role_permissions(role_id)`,
		},
		{
			"idx_roleperm_permission",
			`CREATE INDEX IF NOT EXISTS idx_roleperm_permission ON role_permissions(permission_id)`,
		},

		// ── User ──────────────────────────────────────────────
		{
			"idx_user_email",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_user_email ON "User"(email)`,
		},
		{
			"idx_user_role_id",
			// role_id menggantikan kolom role (string) — FK lookup
			`CREATE INDEX IF NOT EXISTS idx_user_role_id ON "User"(role_id)`,
		},
		{
			"idx_user_role_created",
			// Paginasi user per role, urut terbaru
			`CREATE INDEX IF NOT EXISTS idx_user_role_created ON "User"(role_id, created_at)`,
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
			// Cari semua kursi economy/business/first di satu pesawat
			`CREATE INDEX IF NOT EXISTS idx_seat_aircraft_class ON "Seat"(aircraft_id, seat_class)`,
		},
		{
			"idx_seat_aircraft_number",
			// Cegah nomor kursi duplikat di pesawat yang sama
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_seat_aircraft_number ON "Seat"(aircraft_id, seat_number)`,
		},

		// ── Flight ────────────────────────────────────────────
		{
			"idx_flight_number",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_flight_number ON "Flight"(flight_number)`,
		},
		{
			"idx_flight_dep_arr_time",
			// Query utama: cari penerbangan CGK → DPS tanggal tertentu
			`CREATE INDEX IF NOT EXISTS idx_flight_dep_arr_time ON "Flight"(departure_airport_id, arrival_airport_id, departure_time)`,
		},
		{
			"idx_flight_status_departure",
			// Filter penerbangan aktif/scheduled yang akan datang
			`CREATE INDEX IF NOT EXISTS idx_flight_status_departure ON "Flight"(status, departure_time)`,
		},
		{
			"idx_flight_aircraft_status",
			`CREATE INDEX IF NOT EXISTS idx_flight_aircraft_status ON "Flight"(aircraft_id, status)`,
		},

		// ── FlightSeat ────────────────────────────────────────
		{
			"idx_flightseat_flight_status",
			// Query terpenting: ambil kursi tersedia untuk satu penerbangan
			`CREATE INDEX IF NOT EXISTS idx_flightseat_flight_status ON "FlightSeat"(flight_id, status)`,
		},
		{
			"idx_flightseat_flight_seat",
			// Cegah satu seat muncul dua kali di flight yang sama
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_flightseat_flight_seat ON "FlightSeat"(flight_id, seat_id)`,
		},

		// ── Booking ───────────────────────────────────────────
		{
			"idx_booking_code",
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_booking_code ON "booking"(booking_code)`,
		},
		{
			"idx_booking_user_status",
			// Riwayat booking seorang user berdasarkan status
			`CREATE INDEX IF NOT EXISTS idx_booking_user_status ON "booking"(user_id, status)`,
		},
		{
			"idx_booking_flight_status",
			// Berapa tiket terjual per penerbangan
			`CREATE INDEX IF NOT EXISTS idx_booking_flight_status ON "booking"(flight_id, status)`,
		},
		{
			"idx_booking_status_created",
			// Admin: list semua booking pending urut terbaru
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
			// Satu flight seat hanya boleh ada di satu booking — mencegah double-booking di level DB
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
