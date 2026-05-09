package cache

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	TTLRBACRole       = 30 * 60 // 30 minutes (seconds for reference)
	TTLRBACPermission = 30 * 60 // 30 minutes

	TTLFlightSearch  = 5 * 60  // 5 minutes  — availability/price changes quickly
	TTLFlightDetail  = 10 * 60 // 10 minutes — schedule detail changes rarely
	TTLFlightSeatMap = 2 * 60  // 2 minutes  — seat status is highly volatile

	TTLAirport        = 60 * 60   // 1 hour
	TTLAircraft       = 60 * 60   // 1 hour
	TTLFlightSchedule = 30 * 60   // 30 minutes
	TTLSeatClass      = 24 * 3600 // 24 hours — near-static lookup

	TTLBookingPNR = 5 * 60 // 5 minutes

	TTLUserProfile = 15 * 60 // 15 minutes
)

func KeyUserRole(uid uuid.UUID) string {
	return fmt.Sprintf("rbac:user_role:%s", uid)
}

func KeyUserPermission(uid uuid.UUID, resource, action string) string {
	return fmt.Sprintf("rbac:permission:%s:%s:%s", uid, resource, action)
}

func PatternUserRBAC(uid uuid.UUID) string {
	return fmt.Sprintf("rbac:*:%s:*", uid)
}

func KeyUserByUID(uid uuid.UUID) string {
	return fmt.Sprintf("user:uid:%s", uid)
}

func KeyAirportByID(id uuid.UUID) string {
	return fmt.Sprintf("airport:id:%s", id)
}

func KeyAirportByCode(code string) string {
	return fmt.Sprintf("airport:code:%s", code)
}

func KeyAirportList() string {
	return "airport:list"
}

func KeyAircraftByID(id uuid.UUID) string {
	return fmt.Sprintf("aircraft:id:%s", id)
}

func KeyAircraftWithSeats(id uuid.UUID) string {
	return fmt.Sprintf("aircraft:seats:%s", id)
}

func KeyAircraftList() string {
	return "aircraft:list"
}

func KeyFlightScheduleByID(id uuid.UUID) string {
	return fmt.Sprintf("schedule:id:%s", id)
}

func KeyFlightScheduleList() string {
	return "schedule:list"
}

func KeyFlightScheduleByRoute(depCode, arrCode string) string {
	return fmt.Sprintf("schedule:route:%s:%s", depCode, arrCode)
}

func KeyFlightByID(id uuid.UUID) string {
	return fmt.Sprintf("flight:id:%s", id)
}

func KeyFlightSearch(depCode, arrCode, date string) string {
	return fmt.Sprintf("flight:search:%s:%s:%s", depCode, arrCode, date)
}

func KeyFlightSeatMap(flightID uuid.UUID) string {
	return fmt.Sprintf("flight:seatmap:%s", flightID)
}

func PatternFlight(flightID uuid.UUID) string {
	return fmt.Sprintf("flight:*:%s*", flightID)
}

func KeyPNRByID(id uuid.UUID) string {
	return fmt.Sprintf("pnr:id:%s", id)
}

func KeyPNRByLocator(locator string) string {
	return fmt.Sprintf("pnr:locator:%s", locator)
}
