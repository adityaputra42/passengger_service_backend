package db

import (
	"fmt"
	"log"
	"passenger_service_backend/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SeedDatabase is the main entry point. It skips if data already exists.
func SeedDatabase() error {
	var count int64
	DB.Model(&models.Airport{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Starting database seeding...")

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
	economyRows    int    // rows of economy seats (A-F per row)
	businessRows   int    // rows of business seats (A-D per row)
	firstClassRows int    // rows of first class seats (A-B per row)
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
			// Calculate total seats
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

			// Seed First Class seats (rows 1-N, cols A-B)
			for row := 1; row <= seed.firstClassRows; row++ {
				for _, col := range []string{"A", "B"} {
					seat := models.Seat{
						AircraftID: aircraft.ID,
						SeatNumber: fmt.Sprintf("%d%s", row, col),
						SeatClass:  models.SeatFirst,
					}
					if err := DB.Create(&seat).Error; err != nil {
						return fmt.Errorf("create seat %s: %w", seat.SeatNumber, err)
					}
				}
			}

			// Seed Business Class seats
			bizStartRow := seed.firstClassRows + 1
			for row := bizStartRow; row < bizStartRow+seed.businessRows; row++ {
				for _, col := range []string{"A", "B", "C", "D"} {
					seat := models.Seat{
						AircraftID: aircraft.ID,
						SeatNumber: fmt.Sprintf("%d%s", row, col),
						SeatClass:  models.SeatBusiness,
					}
					if err := DB.Create(&seat).Error; err != nil {
						return fmt.Errorf("create seat %s: %w", seat.SeatNumber, err)
					}
				}
			}

			// Seed Economy Class seats
			econStartRow := seed.firstClassRows + seed.businessRows + 1
			for row := econStartRow; row < econStartRow+seed.economyRows; row++ {
				for _, col := range []string{"A", "B", "C", "D", "E", "F"} {
					seat := models.Seat{
						AircraftID: aircraft.ID,
						SeatNumber: fmt.Sprintf("%d%s", row, col),
						SeatClass:  models.SeatEconomy,
					}
					if err := DB.Create(&seat).Error; err != nil {
						return fmt.Errorf("create seat %s: %w", seat.SeatNumber, err)
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
	users := []struct {
		email    string
		fullName string
		password string
		role     models.UserRole
	}{
		{"admin@flightapp.com", "System Administrator", "Admin@12345", models.RoleAdmin},
		{"budi.santoso@gmail.com", "Budi Santoso", "Password123", models.RoleCustomer},
		{"siti.rahayu@gmail.com", "Siti Rahayu", "Password123", models.RoleCustomer},
		{"agus.wijaya@gmail.com", "Agus Wijaya", "Password123", models.RoleCustomer},
		{"dewi.lestari@gmail.com", "Dewi Lestari", "Password123", models.RoleCustomer},
	}

	for _, u := range users {
		var existing models.User
		if err := DB.Where("email = ?", u.email).First(&existing).Error; err != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("hash password for %s: %w", u.email, err)
			}
			user := models.User{
				Email:        u.email,
				FullName:     u.fullName,
				PasswordHash: string(hash),
				Role:         u.role,
			}
			if err := DB.Create(&user).Error; err != nil {
				return fmt.Errorf("create user %s: %w", u.email, err)
			}
			log.Printf("  Created user: %s (%s)", u.fullName, u.role)
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
		{
			FirstName:      "Budi",
			LastName:       "Santoso",
			DateOfBirth:    mustParseDate("1990-05-15"),
			PassportNumber: "A12345678",
		},
		{
			FirstName:      "Siti",
			LastName:       "Rahayu",
			DateOfBirth:    mustParseDate("1992-08-22"),
			PassportNumber: "B87654321",
		},
		{
			FirstName:      "Agus",
			LastName:       "Wijaya",
			DateOfBirth:    mustParseDate("1988-03-10"),
			PassportNumber: "C11223344",
		},
		{
			FirstName:      "Dewi",
			LastName:       "Lestari",
			DateOfBirth:    mustParseDate("1995-11-30"),
			PassportNumber: "D44332211",
		},
		{
			FirstName:      "Rudi",
			LastName:       "Hermawan",
			DateOfBirth:    mustParseDate("1985-07-04"),
			PassportNumber: "E55667788",
		},
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
	// Fetch airports and first aircraft
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
		return fmt.Errorf("no aircraft found for seeding flights")
	}

	now := time.Now()

	flightDefs := []struct {
		number    string
		dep, arr  string
		depOffset time.Duration // hours from now
		duration  time.Duration
		basePrice float64
		acIdx     int
	}{
		{"GA-401", "CGK", "DPS", 2, 2 * time.Hour, 850000, 0},
		{"GA-402", "DPS", "CGK", 6, 2 * time.Hour, 850000, 0},
		{"GA-101", "CGK", "SUB", 3, 1*time.Hour + 20*time.Minute, 650000, 1},
		{"GA-102", "SUB", "CGK", 7, 1*time.Hour + 20*time.Minute, 650000, 1},
		{"GA-601", "CGK", "UPG", 4, 3 * time.Hour, 1100000, 2},
		{"GA-201", "CGK", "KNO", 5, 3*time.Hour + 30*time.Minute, 1250000, 2},
		{"GA-701", "CGK", "SIN", 8, 2*time.Hour + 10*time.Minute, 2500000, 3},
		{"GA-702", "SIN", "CGK", 14, 2*time.Hour + 10*time.Minute, 2500000, 3},
		{"GA-801", "CGK", "KUL", 10, 2 * time.Hour, 1800000, 0},
		{"GA-901", "DPS", "SIN", 12, 2*time.Hour + 30*time.Minute, 2200000, 1},
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

			// Create FlightSeats from the aircraft's seats
			var seats []models.Seat
			DB.Where("aircraft_id = ?", ac.ID).Find(&seats)

			for _, seat := range seats {
				// Price multiplier by class
				price := fd.basePrice
				switch seat.SeatClass {
				case models.SeatBusiness:
					price = fd.basePrice * 2.5
				case models.SeatFirst:
					price = fd.basePrice * 5.0
				}

				fs := models.FlightSeat{
					FlightID: flight.ID,
					SeatID:   seat.ID,
					Price:    price,
					Status:   true, // available
				}
				if err := DB.Create(&fs).Error; err != nil {
					return fmt.Errorf("create flight seat for flight %s: %w", fd.number, err)
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
// Bookings + BookingPassengers + Payments
// ─────────────────────────────────────────────

func seedBookings() error {
	// Fetch a customer user
	var customer models.User
	if err := DB.Where("role = ?", models.RoleCustomer).First(&customer).Error; err != nil {
		return fmt.Errorf("no customer user found: %w", err)
	}

	// Fetch a flight
	var flight models.Flight
	if err := DB.Where("flight_number = ?", "GA-401").First(&flight).Error; err != nil {
		return fmt.Errorf("flight GA-401 not found: %w", err)
	}

	// Fetch an available economy FlightSeat for this flight
	var flightSeat models.FlightSeat
	if err := DB.
		Joins("JOIN \"Seat\" ON \"Seat\".id = \"FlightSeat\".seat_id").
		Where("\"FlightSeat\".flight_id = ? AND \"FlightSeat\".status = true AND \"Seat\".seat_class = ?", flight.ID, models.SeatEconomy).
		First(&flightSeat).Error; err != nil {
		return fmt.Errorf("no available economy seat for GA-401: %w", err)
	}

	// Fetch a passenger
	var passenger models.Passenger
	if err := DB.Where("passport_number = ?", "A12345678").First(&passenger).Error; err != nil {
		return fmt.Errorf("passenger not found: %w", err)
	}

	// Check if booking already exists
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

		// BookingPassenger
		bp := models.BookingPassenger{
			BookingID:    booking.ID,
			PassengerID:  passenger.ID,
			FlightSeatID: flightSeat.ID,
		}
		if err := DB.Create(&bp).Error; err != nil {
			return fmt.Errorf("create booking passenger: %w", err)
		}

		// Mark seat as booked
		DB.Model(&flightSeat).Update("status", false)

		// Payment
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

// ForceSeedDatabase truncates all tables and re-runs the seeder.
// Use only in development/staging.
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
	}
	for _, t := range tables {
		if err := DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", t)).Error; err != nil {
			return fmt.Errorf("truncate %s: %w", t, err)
		}
	}
	log.Println("All tables truncated. Re-seeding...")
	return SeedDatabase()
}
