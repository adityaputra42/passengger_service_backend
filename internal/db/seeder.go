package db

import (
	"fmt"
	"log"
	"math/rand"
	"passenger_service_backend/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() error {
	var count int64
	DB.Model(&models.Permission{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Starting database seeding...")

	steps := []struct {
		name string
		fn   func() error
	}{
		{"permissions", seedPermissions},
		{"roles", seedRoles},
		{"users", seedUsers},
		{"seat_classes", seedSeatClasses},
		{"ssr_types", seedSSRTypes},
		{"meals", seedMeals},
		{"airports", seedAirports},
		{"aircrafts", seedAircrafts},
		{"flight_schedules", seedFlightSchedules},
		{"flights", seedFlights},
		{"pnr_demo", seedDemoPNR},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			return fmt.Errorf("seed %s: %w", step.name, err)
		}
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
		{Name: "roles.create", Resource: "roles", Action: "create", Description: "Create roles"},
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
		// PNR
		{Name: "pnr.create", Resource: "pnr", Action: "create", Description: "Create PNR / booking"},
		{Name: "pnr.read", Resource: "pnr", Action: "read", Description: "View all PNRs"},
		{Name: "pnr.read_own", Resource: "pnr", Action: "read_own", Description: "View own PNRs"},
		{Name: "pnr.update", Resource: "pnr", Action: "update", Description: "Update PNR"},
		{Name: "pnr.cancel", Resource: "pnr", Action: "cancel", Description: "Cancel PNR"},
		// Seats
		{Name: "seats.read", Resource: "seats", Action: "read", Description: "View seat map"},
		{Name: "seats.assign", Resource: "seats", Action: "assign", Description: "Assign seats"},
		{Name: "seats.manage", Resource: "seats", Action: "manage", Description: "Manage seat inventory"},
		// Check-in
		{Name: "checkin.perform", Resource: "checkin", Action: "perform", Description: "Perform check-in"},
		{Name: "checkin.read", Resource: "checkin", Action: "read", Description: "View check-in records"},
		// Boarding
		{Name: "boarding.read", Resource: "boarding", Action: "read", Description: "View boarding passes"},
		{Name: "boarding.issue", Resource: "boarding", Action: "issue", Description: "Issue boarding passes"},
		// Payments
		{Name: "payments.create", Resource: "payments", Action: "create", Description: "Create payments"},
		{Name: "payments.read", Resource: "payments", Action: "read", Description: "View all payments"},
		{Name: "payments.read_own", Resource: "payments", Action: "read_own", Description: "View own payments"},
		{Name: "payments.refund", Resource: "payments", Action: "refund", Description: "Process refunds"},
		// Dashboard
		{Name: "dashboard.read", Resource: "dashboard", Action: "read", Description: "View dashboard"},
		{Name: "reports.read", Resource: "reports", Action: "read", Description: "View reports"},
	}

	for _, p := range permissions {
		var existing models.Permission
		if err := DB.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&p).Error; err != nil {
				return fmt.Errorf("create permission %s: %w", p.Name, err)
			}
			log.Printf("  Created permission: %s", p.Name)
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
			"permissions.read", "profile.read", "profile.update",
			"airports.create", "airports.read", "airports.update", "airports.delete",
			"aircrafts.create", "aircrafts.read", "aircrafts.update", "aircrafts.delete",
			"flights.create", "flights.read", "flights.update", "flights.delete",
			"pnr.create", "pnr.read", "pnr.update", "pnr.cancel",
			"seats.read", "seats.assign", "seats.manage",
			"checkin.perform", "checkin.read",
			"boarding.read", "boarding.issue",
			"payments.create", "payments.read", "payments.refund",
			"dashboard.read", "reports.read",
		},
		"admin": {
			"users.read", "users.update",
			"roles.read", "permissions.read",
			"profile.read", "profile.update",
			"airports.read", "aircrafts.read",
			"flights.create", "flights.read", "flights.update",
			"pnr.read", "pnr.update", "pnr.cancel",
			"seats.read", "seats.assign", "seats.manage",
			"checkin.perform", "checkin.read",
			"boarding.read", "boarding.issue",
			"payments.read", "payments.refund",
			"dashboard.read", "reports.read",
		},
		"agent": {
			"profile.read", "profile.update",
			"airports.read", "aircrafts.read", "flights.read",
			"pnr.create", "pnr.read", "pnr.update", "pnr.cancel",
			"seats.read", "seats.assign",
			"checkin.perform", "checkin.read",
			"boarding.read", "boarding.issue",
			"payments.create", "payments.read",
		},
		"customer": {
			"profile.read", "profile.update",
			"airports.read", "flights.read",
			"pnr.create", "pnr.read_own", "pnr.cancel",
			"seats.read", "seats.assign",
			"payments.create", "payments.read_own",
		},
	}

	roles := []models.Role{
		{Name: "super_admin", Description: "Full system access", Level: 4, IsSystemRole: true},
		{Name: "admin", Description: "Administrative access", Level: 3, IsSystemRole: true},
		{Name: "agent", Description: "Booking and check-in agent", Level: 2, IsSystemRole: true},
		{Name: "customer", Description: "Passenger self-service access", Level: 1, IsSystemRole: true},
	}

	for _, role := range roles {
		var existing models.Role
		if err := DB.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&role).Error; err != nil {
				return fmt.Errorf("create role %s: %w", role.Name, err)
			}
			var perms []*models.Permission
			if err := DB.Where("name IN ?", rolePermissions[role.Name]).Find(&perms).Error; err != nil {
				return fmt.Errorf("fetch permissions for %s: %w", role.Name, err)
			}
			if err := DB.Model(&role).Association("Permissions").Replace(perms); err != nil {
				return fmt.Errorf("assign permissions to %s: %w", role.Name, err)
			}
			log.Printf("  Created role: %s (level=%d, permissions=%d)", role.Name, role.Level, len(perms))
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Users
// ─────────────────────────────────────────────

func seedUsers() error {
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
		{"superadmin@airline.com", "System Administrator", "SuperAdmin@123", "super_admin"},
		{"admin@airline.com", "Operations Admin", "Admin@12345", "admin"},
		{"agent@airline.com", "Booking Agent", "Agent@12345", "agent"},
		{"budi.santoso@gmail.com", "Budi Santoso", "Password123", "customer"},
		{"siti.rahayu@gmail.com", "Siti Rahayu", "Password123", "customer"},
		{"agus.wijaya@gmail.com", "Agus Wijaya", "Password123", "customer"},
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
				return fmt.Errorf("hash password: %w", err)
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
			log.Printf("  Created user: %s (%s)", u.fullName, u.roleName)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// SeatClasses
// ─────────────────────────────────────────────

func seedSeatClasses() error {
	classes := []models.SeatClass{
		{Code: "F", Name: "First Class"},
		{Code: "C", Name: "Business Class"},
		{Code: "Y", Name: "Economy Class"},
	}
	for _, sc := range classes {
		var existing models.SeatClass
		if err := DB.Where("code = ?", sc.Code).First(&existing).Error; err != nil {
			if err := DB.Create(&sc).Error; err != nil {
				return fmt.Errorf("create seat class %s: %w", sc.Code, err)
			}
			log.Printf("  Created seat class: %s (%s)", sc.Name, sc.Code)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// SSR Types
// ─────────────────────────────────────────────

func seedSSRTypes() error {
	ssrTypes := []models.SSRType{
		{Code: "WCHR", Name: "Wheelchair - Ramp"},
		{Code: "WCHS", Name: "Wheelchair - Steps"},
		{Code: "WCHC", Name: "Wheelchair - Cabin Seat"},
		{Code: "BLND", Name: "Blind Passenger"},
		{Code: "DEAF", Name: "Deaf Passenger"},
		{Code: "UMNR", Name: "Unaccompanied Minor"},
		{Code: "VGML", Name: "Vegetarian Meal"},
		{Code: "KSML", Name: "Kosher Meal"},
		{Code: "MOML", Name: "Muslim Meal"},
		{Code: "PETC", Name: "Pet in Cabin"},
	}
	for _, s := range ssrTypes {
		var existing models.SSRType
		if err := DB.Where("code = ?", s.Code).First(&existing).Error; err != nil {
			if err := DB.Create(&s).Error; err != nil {
				return fmt.Errorf("create SSR type %s: %w", s.Code, err)
			}
			log.Printf("  Created SSR type: %s", s.Code)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Meals
// ─────────────────────────────────────────────

func seedMeals() error {
	meals := []models.Meal{
		{Code: "BBML", Name: "Baby Meal"},
		{Code: "BLML", Name: "Bland Meal"},
		{Code: "CHML", Name: "Child Meal"},
		{Code: "DBML", Name: "Diabetic Meal"},
		{Code: "FPML", Name: "Fruit Platter"},
		{Code: "GFML", Name: "Gluten-Free Meal"},
		{Code: "HNML", Name: "Hindu Meal"},
		{Code: "KSML", Name: "Kosher Meal"},
		{Code: "LCML", Name: "Low Calorie Meal"},
		{Code: "MOML", Name: "Muslim Meal"},
		{Code: "NLML", Name: "Low Lactose Meal"},
		{Code: "RVML", Name: "Raw Vegetarian Meal"},
		{Code: "SFML", Name: "Seafood Meal"},
		{Code: "VGML", Name: "Vegan Meal"},
		{Code: "VJML", Name: "Vegetarian Jain Meal"},
		{Code: "VLML", Name: "Vegetarian Lacto-Ovo Meal"},
	}
	for _, m := range meals {
		var existing models.Meal
		if err := DB.Where("code = ?", m.Code).First(&existing).Error; err != nil {
			if err := DB.Create(&m).Error; err != nil {
				return fmt.Errorf("create meal %s: %w", m.Code, err)
			}
			log.Printf("  Created meal: %s", m.Code)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Airports
// ─────────────────────────────────────────────

func seedAirports() error {
	airports := []models.Airport{
		{Code: "CGK", Name: "Soekarno-Hatta International Airport", City: "Tangerang", Country: "Indonesia", Timezone: "Asia/Jakarta"},
		{Code: "DPS", Name: "Ngurah Rai International Airport", City: "Denpasar", Country: "Indonesia", Timezone: "Asia/Makassar"},
		{Code: "SUB", Name: "Juanda International Airport", City: "Surabaya", Country: "Indonesia", Timezone: "Asia/Jakarta"},
		{Code: "UPG", Name: "Sultan Hasanuddin International Airport", City: "Makassar", Country: "Indonesia", Timezone: "Asia/Makassar"},
		{Code: "KNO", Name: "Kualanamu International Airport", City: "Medan", Country: "Indonesia", Timezone: "Asia/Jakarta"},
		{Code: "BPN", Name: "Sultan Aji Muhammad Sulaiman Airport", City: "Balikpapan", Country: "Indonesia", Timezone: "Asia/Makassar"},
		{Code: "LOP", Name: "Lombok International Airport", City: "Lombok", Country: "Indonesia", Timezone: "Asia/Makassar"},
		{Code: "PLM", Name: "Sultan Mahmud Badaruddin II Airport", City: "Palembang", Country: "Indonesia", Timezone: "Asia/Jakarta"},
		{Code: "SIN", Name: "Changi Airport", City: "Singapore", Country: "Singapore", Timezone: "Asia/Singapore"},
		{Code: "KUL", Name: "Kuala Lumpur International Airport", City: "Kuala Lumpur", Country: "Malaysia", Timezone: "Asia/Kuala_Lumpur"},
	}
	for _, a := range airports {
		var existing models.Airport
		if err := DB.Where("code = ?", a.Code).First(&existing).Error; err != nil {
			if err := DB.Create(&a).Error; err != nil {
				return fmt.Errorf("create airport %s: %w", a.Code, err)
			}
			log.Printf("  Created airport: %s (%s)", a.Name, a.Code)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Aircrafts + AircraftSeats
// ─────────────────────────────────────────────

type acSeed struct {
	model        string
	manufacturer string
	firstRows    int // cols A-B
	bizRows      int // cols A-D
	econRows     int // cols A-F
}

func seedAircrafts() error {
	// Fetch seat class IDs
	classMap := map[string]models.SeatClass{}
	var classes []models.SeatClass
	DB.Find(&classes)
	for _, c := range classes {
		classMap[c.Code] = c
	}

	seeds := []acSeed{
		{model: "737-800", manufacturer: "Boeing", firstRows: 0, bizRows: 4, econRows: 26},
		{model: "A320", manufacturer: "Airbus", firstRows: 0, bizRows: 4, econRows: 24},
		{model: "777-300ER", manufacturer: "Boeing", firstRows: 2, bizRows: 8, econRows: 38},
		{model: "A330-300", manufacturer: "Airbus", firstRows: 2, bizRows: 6, econRows: 35},
	}

	for _, seed := range seeds {
		var existing models.Aircraft
		if err := DB.Where("model = ? AND manufacturer = ?", seed.model, seed.manufacturer).First(&existing).Error; err != nil {
			first := seed.firstRows * 2
			biz := seed.bizRows * 4
			econ := seed.econRows * 6
			total := first + biz + econ

			aircraft := models.Aircraft{
				Model:        seed.model,
				Manufacturer: seed.manufacturer,
				TotalSeats:   total,
			}
			if err := DB.Create(&aircraft).Error; err != nil {
				return fmt.Errorf("create aircraft %s: %w", seed.model, err)
			}

			row := 1
			// First class
			if fc, ok := classMap["F"]; ok {
				classID := fc.ID
				for r := 0; r < seed.firstRows; r++ {
					for _, letter := range []string{"A", "B"} {
						y := (r * 2) + map[string]int{"A": 0, "B": 1}[letter]
						seat := models.AircraftSeat{
							AircraftID:  aircraft.ID,
							SeatNumber:  fmt.Sprintf("%d%s", row, letter),
							RowNumber:   row,
							SeatLetter:  letter,
							XPosition:   row,
							YPosition:   y,
							SeatClassID: &classID,
							SeatType:    "window",
						}
						if letter == "B" {
							seat.SeatType = "aisle"
						}
						if err := DB.Create(&seat).Error; err != nil {
							return fmt.Errorf("create first seat: %w", err)
						}
					}
					row++
				}
			}
			// Business class
			if bc, ok := classMap["C"]; ok {
				classID := bc.ID
				bizLetters := []string{"A", "B", "C", "D"}
				for r := 0; r < seed.bizRows; r++ {
					for i, letter := range bizLetters {
						seatType := "middle"
						switch letter {
						case "A", "D":
							seatType = "window"
						case "B", "C":
							seatType = "aisle"
						}
						seat := models.AircraftSeat{
							AircraftID:  aircraft.ID,
							SeatNumber:  fmt.Sprintf("%d%s", row, letter),
							RowNumber:   row,
							SeatLetter:  letter,
							XPosition:   row,
							YPosition:   i,
							SeatClassID: &classID,
							SeatType:    seatType,
						}
						if err := DB.Create(&seat).Error; err != nil {
							return fmt.Errorf("create business seat: %w", err)
						}
					}
					row++
				}
			}
			// Economy class
			if ec, ok := classMap["Y"]; ok {
				classID := ec.ID
				econLetters := []string{"A", "B", "C", "D", "E", "F"}
				for r := 0; r < seed.econRows; r++ {
					isExit := (r == seed.econRows/2)
					for i, letter := range econLetters {
						seatType := "middle"
						switch letter {
						case "A", "F":
							seatType = "window"
						case "C", "D":
							seatType = "aisle"
						}
						seat := models.AircraftSeat{
							AircraftID:  aircraft.ID,
							SeatNumber:  fmt.Sprintf("%d%s", row, letter),
							RowNumber:   row,
							SeatLetter:  letter,
							XPosition:   row,
							YPosition:   i,
							SeatClassID: &classID,
							SeatType:    seatType,
							IsExitRow:   isExit,
						}
						if err := DB.Create(&seat).Error; err != nil {
							return fmt.Errorf("create economy seat: %w", err)
						}
					}
					row++
				}
			}
			log.Printf("  Created aircraft: %s %s (%d seats)", seed.manufacturer, seed.model, total)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// FlightSchedules
// ─────────────────────────────────────────────

func seedFlightSchedules() error {
	apMap := map[string]models.Airport{}
	var airports []models.Airport
	DB.Find(&airports)
	for _, a := range airports {
		apMap[a.Code] = a
	}

	type scheduleSeed struct {
		flightNo string
		dep, arr string
		depTime  string
		arrTime  string
		days     string
	}

	schedules := []scheduleSeed{
		{"GA-401", "CGK", "DPS", "06:00:00", "07:55:00", "1,2,3,4,5,6,7"},
		{"GA-402", "DPS", "CGK", "09:00:00", "10:55:00", "1,2,3,4,5,6,7"},
		{"GA-101", "CGK", "SUB", "07:00:00", "08:20:00", "1,2,3,4,5,6,7"},
		{"GA-102", "SUB", "CGK", "10:00:00", "11:20:00", "1,2,3,4,5,6,7"},
		{"GA-601", "CGK", "UPG", "08:00:00", "11:00:00", "1,2,3,4,5"},
		{"GA-201", "CGK", "KNO", "06:30:00", "10:00:00", "1,2,3,4,5,6,7"},
		{"GA-701", "CGK", "SIN", "09:00:00", "12:10:00", "1,2,3,4,5,6,7"},
		{"GA-702", "SIN", "CGK", "14:00:00", "17:10:00", "1,2,3,4,5,6,7"},
		{"GA-801", "CGK", "KUL", "10:00:00", "13:00:00", "1,2,3,4,5,6"},
		{"GA-901", "DPS", "SIN", "12:00:00", "14:30:00", "1,3,5,7"},
	}

	for _, s := range schedules {
		var existing models.FlightSchedule
		if err := DB.Where("flight_number = ?", s.flightNo).First(&existing).Error; err != nil {
			dep := apMap[s.dep]
			arr := apMap[s.arr]
			sched := models.FlightSchedule{
				FlightNumber:       s.flightNo,
				DepartureAirportID: dep.ID,
				ArrivalAirportID:   arr.ID,
				DepartureTime:      s.depTime,
				ArrivalTime:        s.arrTime,
				OperatingDays:      s.days,
			}
			if err := DB.Create(&sched).Error; err != nil {
				return fmt.Errorf("create schedule %s: %w", s.flightNo, err)
			}
			log.Printf("  Created schedule: %s (%s → %s)", s.flightNo, s.dep, s.arr)
		}
	}
	return nil
}

// ─────────────────────────────────────────────
// Flights + FlightSeats
// ─────────────────────────────────────────────

func seedFlights() error {
	var schedules []models.FlightSchedule
	if err := DB.Preload("DepartureAirport").Preload("ArrivalAirport").Find(&schedules).Error; err != nil {
		return fmt.Errorf("fetch schedules: %w", err)
	}

	var aircrafts []models.Aircraft
	DB.Find(&aircrafts)
	if len(aircrafts) == 0 {
		return fmt.Errorf("no aircrafts found")
	}

	now := time.Now()

	for i, sched := range schedules {
		ac := aircrafts[i%len(aircrafts)]

		// Parse schedule time to build departure datetime
		depHour, depMin := parseTime(sched.DepartureTime)
		arrHour, arrMin := parseTime(sched.ArrivalTime)

		for dayOffset := 0; dayOffset < 7; dayOffset++ {
			date := now.AddDate(0, 0, dayOffset).Truncate(24 * time.Hour)
			depTime := date.Add(time.Duration(depHour)*time.Hour + time.Duration(depMin)*time.Minute)
			arrTime := date.Add(time.Duration(arrHour)*time.Hour + time.Duration(arrMin)*time.Minute)
			if arrTime.Before(depTime) {
				arrTime = arrTime.AddDate(0, 0, 1) // overnight flight
			}

			var existing models.Flight
			schedID := sched.ID
			acID := ac.ID
			if err := DB.Where("schedule_id = ? AND departure_time = ?", schedID, depTime).First(&existing).Error; err != nil {
				flight := models.Flight{
					ScheduleID:    &schedID,
					AircraftID:    &acID,
					DepartureTime: &depTime,
					ArrivalTime:   &arrTime,
					Status:        models.FlightStatusScheduled,
				}
				if err := DB.Create(&flight).Error; err != nil {
					return fmt.Errorf("create flight for %s: %w", sched.FlightNumber, err)
				}

				// Create FlightSeats from aircraft seats
				var acSeats []models.AircraftSeat
				DB.Where("aircraft_id = ?", ac.ID).Preload("SeatClass").Find(&acSeats)

				for _, acSeat := range acSeats {
					basePrice := basePriceForClass(acSeat.SeatClass)
					fs := models.FlightSeat{
						FlightID:       &flight.ID,
						AircraftSeatID: &acSeat.ID,
						Price:          basePrice,
						Status:         models.FlightSeatAvailable,
					}
					if err := DB.Create(&fs).Error; err != nil {
						return fmt.Errorf("create flight seat: %w", err)
					}
				}
			}
		}
		log.Printf("  Created 7-day flights for schedule: %s", sched.FlightNumber)
	}
	return nil
}

// ─────────────────────────────────────────────
// Demo PNR (full booking flow)
// ─────────────────────────────────────────────

func seedDemoPNR() error {
	var existing models.PNR
	if err := DB.Where("record_locator = ?", "DEMO01").First(&existing).Error; err == nil {
		log.Println("  Demo PNR already exists")
		return nil
	}

	// Pick first available flight with economy seats
	var flightSeat models.FlightSeat
	if err := DB.
		Joins("JOIN aircraft_seats ON aircraft_seats.id = flight_seats.aircraft_seat_id").
		Joins("JOIN seat_classes ON seat_classes.id = aircraft_seats.seat_class_id").
		Where("flight_seats.status = ? AND seat_classes.code = ?", models.FlightSeatAvailable, "Y").
		First(&flightSeat).Error; err != nil {
		log.Println("  [SKIP] No available economy seat for demo PNR")
		return nil
	}

	var flight models.Flight
	DB.First(&flight, "id = ?", *flightSeat.FlightID)

	// Create PNR
	ttl := time.Now().Add(24 * time.Hour)
	pnr := models.PNR{
		RecordLocator: "DEMO01",
		Status:        models.PNRStatusConfirmed,
		TTL:           &ttl,
	}
	if err := DB.Create(&pnr).Error; err != nil {
		return fmt.Errorf("create PNR: %w", err)
	}

	// Contact
	contact := models.PNRContact{
		PNRID: &pnr.ID,
		Name:  "Budi Santoso",
		Email: "budi.santoso@gmail.com",
		Phone: "+6281234567890",
	}
	if err := DB.Create(&contact).Error; err != nil {
		return fmt.Errorf("create PNR contact: %w", err)
	}

	// Passenger
	dob := mustParseDate("1990-05-15")
	passenger := models.PNRPassenger{
		PNRID:          &pnr.ID,
		FirstName:      "Budi",
		LastName:       "Santoso",
		PassengerType:  models.PassengerADT,
		BirthDate:      &dob,
		PassportNumber: "A12345678",
	}
	if err := DB.Create(&passenger).Error; err != nil {
		return fmt.Errorf("create PNR passenger: %w", err)
	}

	// Segment
	flightID := *flightSeat.FlightID
	segment := models.PNRSegment{
		PNRID:        &pnr.ID,
		FlightID:     &flightID,
		SegmentOrder: 1,
	}
	if err := DB.Create(&segment).Error; err != nil {
		return fmt.Errorf("create PNR segment: %w", err)
	}

	// Lock seat
	now := time.Now()
	expires := now.Add(15 * time.Minute)
	seatLock := models.SeatLock{
		FlightSeatID: &flightSeat.ID,
		PNRID:        &pnr.ID,
		LockedAt:     &now,
		ExpiresAt:    &expires,
	}
	if err := DB.Create(&seatLock).Error; err != nil {
		return fmt.Errorf("create seat lock: %w", err)
	}

	// Assign seat
	assignedAt := time.Now()
	assignment := models.SeatAssignment{
		PassengerID:  &passenger.ID,
		SegmentID:    &segment.ID,
		FlightSeatID: &flightSeat.ID,
		AssignedAt:   &assignedAt,
	}
	if err := DB.Create(&assignment).Error; err != nil {
		return fmt.Errorf("create seat assignment: %w", err)
	}

	// Mark seat as booked
	DB.Model(&flightSeat).Update("status", models.FlightSeatBooked)

	// Payment
	paidAt := time.Now()
	payment := models.Payment{
		PNRID:  &pnr.ID,
		Amount: flightSeat.Price,
		Method: models.PaymentMethodBankTransfer,
		Status: models.PaymentStatusSuccess,
		PaidAt: &paidAt,
	}
	if err := DB.Create(&payment).Error; err != nil {
		return fmt.Errorf("create payment: %w", err)
	}

	// Ticket
	issuedAt := time.Now()
	ticketNum := fmt.Sprintf("GA-%010d", rand.Intn(9999999999))
	ticket := models.Ticket{
		PassengerID:  &passenger.ID,
		TicketNumber: ticketNum,
		IssuedAt:     &issuedAt,
	}
	if err := DB.Create(&ticket).Error; err != nil {
		return fmt.Errorf("create ticket: %w", err)
	}

	// TicketSegment
	ts := models.TicketSegment{TicketID: &ticket.ID, SegmentID: &segment.ID}
	if err := DB.Create(&ts).Error; err != nil {
		return fmt.Errorf("create ticket segment: %w", err)
	}

	// Update PNR status to ticketed
	DB.Model(&pnr).Update("status", models.PNRStatusTicketed)

	log.Printf("  Created demo PNR: DEMO01 (ticket: %s)", ticketNum)
	return nil
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func mustParseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(fmt.Sprintf("mustParseDate %q: %v", s, err))
	}
	return t
}

// parseTime parses "HH:MM:SS" or "HH:MM" → hour, minute
func parseTime(s string) (int, int) {
	var h, m, sec int
	fmt.Sscanf(s, "%d:%d:%d", &h, &m, &sec)
	return h, m
}

func basePriceForClass(sc *models.SeatClass) float64 {
	if sc == nil {
		return 500000
	}
	switch sc.Code {
	case "F":
		return 5000000
	case "C":
		return 2500000
	default:
		return 850000
	}
}

// ForceSeedDatabase truncates all tables and re-runs seeder. Dev only.
func ForceSeedDatabase() error {
	tables := []string{
		"boarding_passes", "checkins", "baggage",
		"passenger_meals", "passenger_ssr",
		"ticket_segments", "tickets",
		"seat_assignments", "seat_locks",
		"pnr_segments", "pnr_passengers", "pnr_contacts",
		"payments", "pnrs",
		"flight_seats", "flights", "flight_schedules",
		"aircraft_seats", "aircrafts",
		"airports", "seat_classes",
		"ssr_types", "meals",
		"users", "role_permissions", "roles", "permissions",
	}
	for _, t := range tables {
		if err := DB.Exec(fmt.Sprintf(`TRUNCATE TABLE "%s" CASCADE`, t)).Error; err != nil {
			return fmt.Errorf("truncate %s: %w", t, err)
		}
	}
	log.Println("All tables truncated. Re-seeding...")
	return SeedDatabase()
}
