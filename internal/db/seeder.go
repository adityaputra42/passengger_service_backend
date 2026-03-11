package db

import (
	"fmt"
	"log"
	"passenger_service_backend/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SeedDatabase is the main entry point. Skips if already seeded.
func SeedDatabase() error {
	var count int64
	DB.Model(&models.Permission{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Starting database seeding...")

	if err := seedPermissions(); err != nil {
		return fmt.Errorf("seedPermissions: %w", err)
	}
	if err := seedRoles(); err != nil {
		return fmt.Errorf("seedRoles: %w", err)
	}
	if err := seedAirports(); err != nil {
		return fmt.Errorf("seedAirports: %w", err)
	}
	if err := seedAircrafts(); err != nil {
		return fmt.Errorf("seedAircrafts: %w", err)
	}
	if err := seedUsers(); err != nil {
		return fmt.Errorf("seedUsers: %w", err)
	}
	if err := seedPassengers(); err != nil {
		return fmt.Errorf("seedPassengers: %w", err)
	}
	if err := seedFlights(); err != nil {
		return fmt.Errorf("seedFlights: %w", err)
	}
	if err := seedBookings(); err != nil {
		return fmt.Errorf("seedBookings: %w", err)
	}

	log.Println("Database seeding completed successfully.")
	return nil
}

// ─────────────────────────────────────────────
// Permissions
// ─────────────────────────────────────────────

func seedPermissions() error {
	permissions := []models.Permission{
		// Users
		{Name: "users.create", Resource: "users", Action: "create", Description: "Create new users"},
		{Name: "users.read", Resource: "users", Action: "read", Description: "View all users"},
		{Name: "users.update", Resource: "users", Action: "update", Description: "Update user data"},
		{Name: "users.delete", Resource: "users", Action: "delete", Description: "Delete users"},
		// Roles
		{Name: "roles.create", Resource: "roles", Action: "create", Description: "Create new roles"},
		{Name: "roles.read", Resource: "roles", Action: "read", Description: "View roles"},
		{Name: "roles.update", Resource: "roles", Action: "update", Description: "Update roles"},
		{Name: "roles.delete", Resource: "roles", Action: "delete", Description: "Delete roles"},
		// Permissions
		{Name: "permissions.read", Resource: "permissions", Action: "read", Description: "View permissions"},
		// Profile
		{Name: "profile.read", Resource: "profile", Action: "read", Description: "View own profile"},
		{Name: "profile.update", Resource: "profile", Action: "update", Description: "Update own profile"},
		// Airports
		{Name: "airports.create", Resource: "airports", Action: "create", Description: "Create airports"},
		{Name: "airports.read", Resource: "airports", Action: "read", Description: "View airports"},
		{Name: "airports.update", Resource: "airports", Action: "update", Description: "Update airports"},
		{Name: "airports.delete", Resource: "airports", Action: "delete", Description: "Delete airports"},
		// Aircrafts
		{Name: "aircrafts.create", Resource: "aircrafts", Action: "create", Description: "Create aircrafts"},
		{Name: "aircrafts.read", Resource: "aircrafts", Action: "read", Description: "View aircrafts"},
		{Name: "aircrafts.update", Resource: "aircrafts", Action: "update", Description: "Update aircrafts"},
		{Name: "aircrafts.delete", Resource: "aircrafts", Action: "delete", Description: "Delete aircrafts"},
		// Flights
		{Name: "flights.create", Resource: "flights", Action: "create", Description: "Create flights"},
		{Name: "flights.read", Resource: "flights", Action: "read", Description: "View all flights"},
		{Name: "flights.update", Resource: "flights", Action: "update", Description: "Update flights"},
		{Name: "flights.delete", Resource: "flights", Action: "delete", Description: "Delete flights"},
		// Seats
		{Name: "seats.create", Resource: "seats", Action: "create", Description: "Create seats"},
		{Name: "seats.read", Resource: "seats", Action: "read", Description: "View seats"},
		{Name: "seats.update", Resource: "seats", Action: "update", Description: "Update seats"},
		// Bookings
		{Name: "bookings.create", Resource: "bookings", Action: "create", Description: "Create bookings"},
		{Name: "bookings.read", Resource: "bookings", Action: "read", Description: "View all bookings"},
		{Name: "bookings.read_own", Resource: "bookings", Action: "read_own", Description: "View own bookings"},
		{Name: "bookings.update", Resource: "bookings", Action: "update", Description: "Update bookings"},
		{Name: "bookings.delete", Resource: "bookings", Action: "delete", Description: "Cancel/delete bookings"},
		// Passengers
		{Name: "passengers.create", Resource: "passengers", Action: "create", Description: "Register passengers"},
		{Name: "passengers.read", Resource: "passengers", Action: "read", Description: "View all passengers"},
		{Name: "passengers.update", Resource: "passengers", Action: "update", Description: "Update passenger data"},
		{Name: "passengers.delete", Resource: "passengers", Action: "delete", Description: "Delete passengers"},
		// Payments
		{Name: "payments.create", Resource: "payments", Action: "create", Description: "Create payments"},
		{Name: "payments.read", Resource: "payments", Action: "read", Description: "View all payments"},
		{Name: "payments.read_own", Resource: "payments", Action: "read_own", Description: "View own payments"},
		{Name: "payments.update", Resource: "payments", Action: "update", Description: "Update payment status"},
		// Dashboard
		{Name: "dashboard.read", Resource: "dashboard", Action: "read", Description: "View dashboard"},
		{Name: "reports.read", Resource: "reports", Action: "read", Description: "View reports & analytics"},
	}

	for _, p := range permissions {
		var existing models.Permission
		if err := DB.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&p).Error; err != nil {
				return fmt.Errorf("create permission %s: %w", p.Name, err)
			}
			log.Printf("  Created permission: %s", p.Name)
		} else {
			log.Printf("  Permission already exists: %s", p.Name)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Roles
// ─────────────────────────────────────────────

func seedRoles() error {
	rolePermissions := map[string][]string{
		"super_admin": {
			"users.create", "users.read", "users.update", "users.delete",
			"roles.create", "roles.read", "roles.update", "roles.delete",
			"permissions.read",
			"profile.read", "profile.update",
			"airports.create", "airports.read", "airports.update", "airports.delete",
			"aircrafts.create", "aircrafts.read", "aircrafts.update", "aircrafts.delete",
			"flights.create", "flights.read", "flights.update", "flights.delete",
			"seats.create", "seats.read", "seats.update",
			"bookings.create", "bookings.read", "bookings.update", "bookings.delete",
			"passengers.create", "passengers.read", "passengers.update", "passengers.delete",
			"payments.create", "payments.read", "payments.update",
			"dashboard.read", "reports.read",
		},
		"admin": {
			"users.read", "users.update",
			"roles.read", "permissions.read",
			"profile.read", "profile.update",
			"airports.read",
			"aircrafts.read",
			"flights.create", "flights.read", "flights.update",
			"seats.read", "seats.update",
			"bookings.read", "bookings.update",
			"passengers.read", "passengers.update",
			"payments.read", "payments.update",
			"dashboard.read", "reports.read",
		},
		"agent": {
			"profile.read", "profile.update",
			"airports.read",
			"aircrafts.read",
			"flights.read",
			"seats.read",
			"bookings.create", "bookings.read", "bookings.update",
			"passengers.create", "passengers.read", "passengers.update",
			"payments.create", "payments.read",
		},
		"customer": {
			"profile.read", "profile.update",
			"airports.read",
			"flights.read",
			"seats.read",
			"bookings.create", "bookings.read_own",
			"passengers.create",
			"payments.create", "payments.read_own",
		},
	}

	roles := []models.Role{
		{Name: "super_admin", Description: "Full system access with all permissions", Level: 4, IsSystemRole: true},
		{Name: "admin", Description: "Administrative access to manage flights and bookings", Level: 3, IsSystemRole: true},
		{Name: "agent", Description: "Booking agent with operational access", Level: 2, IsSystemRole: true},
		{Name: "customer", Description: "Passenger access to search and book flights", Level: 1, IsSystemRole: true},
	}

	for _, role := range roles {
		var existing models.Role
		if err := DB.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&role).Error; err != nil {
				return fmt.Errorf("create role %s: %w", role.Name, err)
			}

			var perms []*models.Permission
			if err := DB.Where("name IN ?", rolePermissions[role.Name]).Find(&perms).Error; err != nil {
				return fmt.Errorf("fetch permissions for role %s: %w", role.Name, err)
			}
			if err := DB.Model(&role).Association("Permissions").Replace(perms); err != nil {
				return fmt.Errorf("assign permissions to role %s: %w", role.Name, err)
			}

			log.Printf("  Created role: %s (level=%d) with %d permissions", role.Name, role.Level, len(perms))
		} else {
			log.Printf("  Role already exists: %s", role.Name)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Airports
// ─────────────────────────────────────────────

func seedAirports() error {
	airports := []models.Airport{
		{Code: "CGK", Name: "Soekarno-Hatta International Airport", City: "Tangerang", Country: "Indonesia"},
		{Code: "DPS", Name: "Ngurah Rai International Airport", City: "Denpasar", Country: "Indonesia"},
		{Code: "SUB", Name: "Juanda International Airport", City: "Surabaya", Country: "Indonesia"},
		{Code: "UPG", Name: "Sultan Hasanuddin International Airport", City: "Makassar", Country: "Indonesia"},
		{Code: "KNO", Name: "Kualanamu International Airport", City: "Medan", Country: "Indonesia"},
		{Code: "BPN", Name: "Sultan Aji Muhammad Sulaiman Airport", City: "Balikpapan", Country: "Indonesia"},
		{Code: "LOP", Name: "Lombok International Airport", City: "Lombok", Country: "Indonesia"},
		{Code: "SIN", Name: "Changi Airport", City: "Singapore", Country: "Singapore"},
		{Code: "KUL", Name: "Kuala Lumpur International Airport", City: "Kuala Lumpur", Country: "Malaysia"},
		{Code: "BKK", Name: "Suvarnabhumi Airport", City: "Bangkok", Country: "Thailand"},
	}

	for _, airport := range airports {
		var existing models.Airport
		if err := DB.Where("code = ?", airport.Code).First(&existing).Error; err != nil {
			if err := DB.Create(&airport).Error; err != nil {
				return fmt.Errorf("create airport %s: %w", airport.Code, err)
			}
			log.Printf("  Created airport: %s (%s)", airport.Name, airport.Code)
		} else {
			log.Printf("  Airport already exists: %s", airport.Code)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Aircrafts + Seats
// ─────────────────────────────────────────────

type aircraftSeed struct {
	model          string
	registration   string
	economyRows    int
	businessRows   int
	firstClassRows int
}

func seedAircrafts() error {
	aircrafts := []aircraftSeed{
		{model: "Boeing 737-800", registration: "PK-GFX", economyRows: 26, businessRows: 4, firstClassRows: 0},
		{model: "Airbus A320", registration: "PK-LAX", economyRows: 24, businessRows: 4, firstClassRows: 0},
		{model: "Boeing 777-300ER", registration: "PK-GII", economyRows: 38, businessRows: 8, firstClassRows: 2},
		{model: "Airbus A330-300", registration: "PK-GPH", economyRows: 35, businessRows: 6, firstClassRows: 2},
	}

	for _, seed := range aircrafts {
		var existing models.Aircraft
		if err := DB.Where("registration_no = ?", seed.registration).First(&existing).Error; err != nil {
			econSeats := seed.economyRows * 6
			bizSeats := seed.businessRows * 4
			firstSeats := seed.firstClassRows * 2
			total := int64(econSeats + bizSeats + firstSeats)

			aircraft := models.Aircraft{
				Model:          seed.model,
				RegistrationNo: seed.registration,
				TotalSeats:     total,
			}
			if err := DB.Create(&aircraft).Error; err != nil {
				return fmt.Errorf("create aircraft %s: %w", seed.registration, err)
			}

			// First Class
			for row := 1; row <= seed.firstClassRows; row++ {
				for _, col := range []string{"A", "B"} {
					seat := models.Seat{AircraftID: aircraft.ID, SeatNumber: fmt.Sprintf("%d%s", row, col), SeatClass: models.SeatFirst}
					if err := DB.Create(&seat).Error; err != nil {
						return fmt.Errorf("create first class seat: %w", err)
					}
				}
			}
			// Business
			bizStart := seed.firstClassRows + 1
			for row := bizStart; row < bizStart+seed.businessRows; row++ {
				for _, col := range []string{"A", "B", "C", "D"} {
					seat := models.Seat{AircraftID: aircraft.ID, SeatNumber: fmt.Sprintf("%d%s", row, col), SeatClass: models.SeatBusiness}
					if err := DB.Create(&seat).Error; err != nil {
						return fmt.Errorf("create business seat: %w", err)
					}
				}
			}
			// Economy
			econStart := seed.firstClassRows + seed.businessRows + 1
			for row := econStart; row < econStart+seed.economyRows; row++ {
				for _, col := range []string{"A", "B", "C", "D", "E", "F"} {
					seat := models.Seat{AircraftID: aircraft.ID, SeatNumber: fmt.Sprintf("%d%s", row, col), SeatClass: models.SeatEconomy}
					if err := DB.Create(&seat).Error; err != nil {
						return fmt.Errorf("create economy seat: %w", err)
					}
				}
			}

			log.Printf("  Created aircraft: %s (%s) with %d seats", aircraft.Model, aircraft.RegistrationNo, total)
		} else {
			log.Printf("  Aircraft already exists: %s", seed.registration)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Users
// ─────────────────────────────────────────────

func seedUsers() error {
	// Build role name → ID map
	roleMap := map[string]uint{}
	var roles []models.Role
	if err := DB.Find(&roles).Error; err != nil {
		return fmt.Errorf("fetch roles: %w", err)
	}
	for _, r := range roles {
		roleMap[r.Name] = r.ID
	}

	users := []struct {
		email    string
		fullName string
		password string
		roleName string
	}{
		{"superadmin@flightapp.com", "System Administrator", "SuperAdmin@123", "super_admin"},
		{"admin@flightapp.com", "Flight Admin", "Admin@12345", "admin"},
		{"agent@flightapp.com", "Booking Agent", "Agent@12345", "agent"},
		{"budi.santoso@gmail.com", "Budi Santoso", "Password123", "customer"},
		{"siti.rahayu@gmail.com", "Siti Rahayu", "Password123", "customer"},
		{"agus.wijaya@gmail.com", "Agus Wijaya", "Password123", "customer"},
		{"dewi.lestari@gmail.com", "Dewi Lestari", "Password123", "customer"},
	}

	for _, u := range users {
		var existing models.User
		if err := DB.Where("email = ?", u.email).First(&existing).Error; err != nil {
			roleID, ok := roleMap[u.roleName]
			if !ok {
				return fmt.Errorf("role %q not found for user %s", u.roleName, u.email)
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("hash password for %s: %w", u.email, err)
			}
			user := models.User{
				Email:        u.email,
				FullName:     u.fullName,
				PasswordHash: string(hash),
				RoleID:       roleID,
			}
			if err := DB.Create(&user).Error; err != nil {
				return fmt.Errorf("create user %s: %w", u.email, err)
			}
			log.Printf("  Created user: %s (role=%s)", u.fullName, u.roleName)
		} else {
			log.Printf("  User already exists: %s", u.email)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Passengers
// ─────────────────────────────────────────────

func seedPassengers() error {
	passengers := []models.Passenger{
		{FirstName: "Budi", LastName: "Santoso", DateOfBirth: mustParseDate("1990-05-15"), PassportNumber: "A12345678"},
		{FirstName: "Siti", LastName: "Rahayu", DateOfBirth: mustParseDate("1992-08-22"), PassportNumber: "B87654321"},
		{FirstName: "Agus", LastName: "Wijaya", DateOfBirth: mustParseDate("1988-03-10"), PassportNumber: "C11223344"},
		{FirstName: "Dewi", LastName: "Lestari", DateOfBirth: mustParseDate("1995-11-30"), PassportNumber: "D44332211"},
		{FirstName: "Rudi", LastName: "Hermawan", DateOfBirth: mustParseDate("1985-07-04"), PassportNumber: "E55667788"},
	}

	for _, p := range passengers {
		var existing models.Passenger
		if err := DB.Where("passport_number = ?", p.PassportNumber).First(&existing).Error; err != nil {
			if err := DB.Create(&p).Error; err != nil {
				return fmt.Errorf("create passenger %s %s: %w", p.FirstName, p.LastName, err)
			}
			log.Printf("  Created passenger: %s %s", p.FirstName, p.LastName)
		} else {
			log.Printf("  Passenger already exists: %s", p.PassportNumber)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Flights + FlightSeats
// ─────────────────────────────────────────────

func seedFlights() error {
	airports := map[string]models.Airport{}
	for _, code := range []string{"CGK", "DPS", "SUB", "UPG", "KNO", "SIN", "KUL"} {
		var ap models.Airport
		if err := DB.Where("code = ?", code).First(&ap).Error; err != nil {
			return fmt.Errorf("airport %s not found: %w", code, err)
		}
		airports[code] = ap
	}

	var aircraft []models.Aircraft
	if err := DB.Find(&aircraft).Error; err != nil || len(aircraft) == 0 {
		return fmt.Errorf("no aircraft found")
	}

	now := time.Now()
	flightDefs := []struct {
		number    string
		dep, arr  string
		depOffset time.Duration
		duration  time.Duration
		basePrice float64
		acIdx     int
	}{
		{"GA-401", "CGK", "DPS", 2, 2 * time.Hour, 850000, 0},
		{"GA-402", "DPS", "CGK", 6, 2 * time.Hour, 850000, 0},
		{"GA-101", "CGK", "SUB", 3, 80 * time.Minute, 650000, 1},
		{"GA-102", "SUB", "CGK", 7, 80 * time.Minute, 650000, 1},
		{"GA-601", "CGK", "UPG", 4, 3 * time.Hour, 1100000, 2},
		{"GA-201", "CGK", "KNO", 5, 210 * time.Minute, 1250000, 2},
		{"GA-701", "CGK", "SIN", 8, 130 * time.Minute, 2500000, 3},
		{"GA-702", "SIN", "CGK", 14, 130 * time.Minute, 2500000, 3},
		{"GA-801", "CGK", "KUL", 10, 2 * time.Hour, 1800000, 0},
		{"GA-901", "DPS", "SIN", 12, 150 * time.Minute, 2200000, 1},
	}

	for _, fd := range flightDefs {
		var existing models.Flight
		if err := DB.Where("flight_number = ?", fd.number).First(&existing).Error; err != nil {
			depTime := now.Add(fd.depOffset * time.Hour).Truncate(time.Minute)
			arrTime := depTime.Add(fd.duration)
			ac := aircraft[fd.acIdx%len(aircraft)]

			flight := models.Flight{
				FlightNumber:       fd.number,
				AircraftID:         ac.ID,
				DepartureAirportID: airports[fd.dep].ID,
				ArrivalAirportID:   airports[fd.arr].ID,
				DepartureTime:      depTime,
				ArrivalTime:        arrTime,
				BasePrice:          fd.basePrice,
				Status:             models.FlightScheduled,
			}
			if err := DB.Create(&flight).Error; err != nil {
				return fmt.Errorf("create flight %s: %w", fd.number, err)
			}

			var seats []models.Seat
			DB.Where("aircraft_id = ?", ac.ID).Find(&seats)

			for _, seat := range seats {
				price := fd.basePrice
				switch seat.SeatClass {
				case models.SeatBusiness:
					price = fd.basePrice * 2.5
				case models.SeatFirst:
					price = fd.basePrice * 5.0
				}
				fs := models.FlightSeat{FlightID: flight.ID, SeatID: seat.ID, Price: price, Status: true}
				if err := DB.Create(&fs).Error; err != nil {
					return fmt.Errorf("create flight seat: %w", err)
				}
			}
			log.Printf("  Created flight: %s (%s → %s) with %d seats", fd.number, fd.dep, fd.arr, len(seats))
		} else {
			log.Printf("  Flight already exists: %s", fd.number)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Bookings + Payments
// ─────────────────────────────────────────────

func seedBookings() error {
	var customer models.User
	if err := DB.Joins("Role").Where("roles.name = ?", "customer").First(&customer).Error; err != nil {
		return fmt.Errorf("no customer user found: %w", err)
	}

	var flight models.Flight
	if err := DB.Where("flight_number = ?", "GA-401").First(&flight).Error; err != nil {
		return fmt.Errorf("flight GA-401 not found: %w", err)
	}

	var flightSeat models.FlightSeat
	if err := DB.
		Joins(`JOIN "Seat" ON "Seat".id = "FlightSeat".seat_id`).
		Where(`"FlightSeat".flight_id = ? AND "FlightSeat".status = true AND "Seat".seat_class = ?`, flight.ID, models.SeatEconomy).
		First(&flightSeat).Error; err != nil {
		return fmt.Errorf("no available economy seat for GA-401: %w", err)
	}

	var passenger models.Passenger
	if err := DB.Where("passport_number = ?", "A12345678").First(&passenger).Error; err != nil {
		return fmt.Errorf("passenger not found: %w", err)
	}

	var existingBooking models.Booking
	if err := DB.Where("booking_code = ?", "BK-SEED-001").First(&existingBooking).Error; err != nil {
		booking := models.Booking{
			UserID:      customer.UID,
			FlightID:    flight.ID,
			BookingCode: "BK-SEED-001",
			TotalAmount: flightSeat.Price,
			Status:      models.BookingConfirmed,
		}
		if err := DB.Create(&booking).Error; err != nil {
			return fmt.Errorf("create booking: %w", err)
		}

		bp := models.BookingPassenger{
			BookingID:    booking.ID,
			PassengerID:  passenger.ID,
			FlightSeatID: flightSeat.ID,
		}
		if err := DB.Create(&bp).Error; err != nil {
			return fmt.Errorf("create booking passenger: %w", err)
		}

		DB.Model(&flightSeat).Update("status", false)

		payment := models.Payment{
			BookingID: booking.ID,
			Amount:    booking.TotalAmount,
			Status:    models.PaymentSuccess,
			Method:    models.MethodBankTransfer,
			PaidAt:    time.Now(),
		}
		if err := DB.Create(&payment).Error; err != nil {
			return fmt.Errorf("create payment: %w", err)
		}
		log.Printf("  Created booking: %s (flight %s, passenger %s %s)", booking.BookingCode, flight.FlightNumber, passenger.FirstName, passenger.LastName)
	} else {
		log.Println("  Booking BK-SEED-001 already exists")
	}
	return nil
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func mustParseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(fmt.Sprintf("mustParseDate: invalid date %q: %v", s, err))
	}
	return t
}

// ForceSeedDatabase truncates all tables and re-runs seeder. Dev/staging only.
func ForceSeedDatabase() error {
	tables := []string{
		`"payment"`,
		`"Booking Passengger"`,
		`"booking"`,
		`"FlightSeat"`,
		`"Flight"`,
		`"Seat"`,
		`"Aircraft"`,
		`"Airport"`,
		`"Passengger"`,
		`"User"`,
		`role_permissions`,
		`roles`,
		`permissions`,
	}
	for _, t := range tables {
		if err := DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", t)).Error; err != nil {
			return fmt.Errorf("truncate %s: %w", t, err)
		}
	}
	log.Println("All tables truncated. Re-seeding...")
	return SeedDatabase()
}
