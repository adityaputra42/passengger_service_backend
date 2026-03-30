package utils

import (
	"errors"
	"net/http"
)

var (
	// Auth
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrTokenExpired       = errors.New("token sudah expired")
	ErrTokenInvalid       = errors.New("token tidak valid")
	ErrUnauthorized       = errors.New("tidak memiliki akses")

	// User
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
	ErrWrongPassword      = errors.New("password lama tidak sesuai")

	// Role
	ErrRoleNotFound      = errors.New("role tidak ditemukan")
	ErrRoleNameDuplicate = errors.New("nama role sudah digunakan")
	ErrSystemRoleDelete  = errors.New("system role tidak bisa dihapus")

	// Airport
	ErrAirportNotFound      = errors.New("airport tidak ditemukan")
	ErrAirportCodeDuplicate = errors.New("kode IATA sudah terdaftar")

	// Aircraft
	ErrAircraftNotFound   = errors.New("aircraft tidak ditemukan")
	ErrAircraftHasFlights = errors.New("aircraft masih memiliki penerbangan aktif")

	// FlightSchedule
	ErrScheduleNotFound      = errors.New("jadwal penerbangan tidak ditemukan")
	ErrFlightNumberDuplicate = errors.New("nomor penerbangan sudah terdaftar")

	// Flight
	ErrFlightNotFound        = errors.New("penerbangan tidak ditemukan")
	ErrFlightNotScheduled    = errors.New("status penerbangan bukan scheduled")
	ErrFlightAlreadyDeparted = errors.New("penerbangan sudah berangkat")
	ErrNoFlightsGenerated    = errors.New("tidak ada penerbangan yang bisa di-generate")

	// FlightSeat
	ErrFlightSeatNotFound = errors.New("kursi tidak ditemukan")
	ErrSeatNotAvailable   = errors.New("kursi tidak tersedia")
	ErrSeatAlreadyLocked  = errors.New("kursi sedang dipesan oleh pengguna lain")
	ErrSeatAlreadyBooked  = errors.New("kursi sudah dipesan")

	// PNR / Booking
	ErrPNRNotFound         = errors.New("PNR tidak ditemukan")
	ErrPNRAlreadyCancelled = errors.New("PNR sudah dibatalkan")
	ErrPNRAlreadyTicketed  = errors.New("PNR sudah diterbitkan tiketnya")
	ErrPNRHoldExpired      = errors.New("waktu hold PNR sudah habis")
	ErrLocatorGenFailed    = errors.New("gagal generate record locator unik")

	// Passenger
	ErrPassengerNotFound = errors.New("penumpang tidak ditemukan")

	// Seat Lock
	ErrSeatLockNotFound = errors.New("seat lock tidak ditemukan")
	ErrSeatLockExpired  = errors.New("seat lock sudah expired")

	// Payment
	ErrPaymentNotFound    = errors.New("pembayaran tidak ditemukan")
	ErrPaymentNotPending  = errors.New("pembayaran tidak dalam status pending")
	ErrPaymentAlreadyPaid = errors.New("pembayaran sudah berhasil diproses")
	ErrPaymentNotSuccess  = errors.New("pembayaran belum berhasil")
	ErrPNRPaymentMismatch = errors.New("pembayaran tidak terkait dengan PNR ini")

	// Ticket
	ErrTicketNotFound      = errors.New("tiket tidak ditemukan")
	ErrTicketAlreadyIssued = errors.New("tiket sudah diterbitkan untuk penumpang ini")

	// Check-in
	ErrAlreadyCheckedIn      = errors.New("penumpang sudah check-in")
	ErrCheckinTooEarly       = errors.New("check-in belum dibuka")
	ErrCheckinClosed         = errors.New("check-in sudah ditutup")
	ErrTicketRequiredCheckin = errors.New("tiket harus diterbitkan sebelum check-in")

	// Boarding Pass
	ErrBoardingPassNotFound    = errors.New("boarding pass tidak ditemukan")
	ErrBoardingPassExists      = errors.New("boarding pass sudah diterbitkan")
	ErrCheckinRequiredBoarding = errors.New("harus check-in sebelum boarding pass diterbitkan")

	// Baggage
	ErrBaggageNotFound        = errors.New("bagasi tidak ditemukan")
	ErrCheckinRequiredBaggage = errors.New("harus check-in sebelum mendaftarkan bagasi")

	// SSR / Meal
	ErrSSRTypeNotFound = errors.New("tipe SSR tidak ditemukan")
	ErrMealNotFound    = errors.New("pilihan meal tidak ditemukan")
	ErrSSRNotFound     = errors.New("SSR penumpang tidak ditemukan")
	ErrMealSSRNotFound = errors.New("meal penumpang tidak ditemukan")

	// Multi-segment / Trip validation
	ErrInvalidSegmentCount      = errors.New("jumlah segment tidak sesuai trip_type")
	ErrSegmentChronologyInvalid = errors.New("urutan waktu segment tidak valid")
	ErrRoundTripAirportMismatch = errors.New("airport return tidak cocok dengan arrival outbound")
	ErrSeatFlightMismatch       = errors.New("seat tidak milik penerbangan yang dipilih")
	ErrDuplicatePassengerSeat   = errors.New("penumpang memilih kursi dua kali di segment yang sama")
	ErrDuplicateSeatSelection   = errors.New("kursi yang sama dipilih dua penumpang berbeda")
	ErrInvalidPassengerIndex    = errors.New("passenger_index tidak valid")
)

// // Errorf writes a structured error response from a string.
// func Errorf(w http.ResponseWriter, status int, msg string) {
// 	Error(w, status, errors.New(msg))
// }

// // ─────────────────────────────────────────────
// // Error → HTTP status mapping
// // ─────────────────────────────────────────────

// // ServiceError maps a domain error to the correct HTTP status code
// // and writes the response. Returns true if an error was written.
// func ServiceError(w http.ResponseWriter, err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	status := statusFor(err)
// 	Error(w, status, err)
// 	return true
// }

func statusFor(err error) int {
	switch {
	// 400
	case errors.Is(err, ErrInvalidSegmentCount),
		errors.Is(err, ErrDuplicatePassengerSeat),
		errors.Is(err, ErrDuplicateSeatSelection),
		errors.Is(err, ErrInvalidPassengerIndex),
		errors.Is(err, ErrWrongPassword):
		return http.StatusBadRequest

	// 401
	case errors.Is(err, ErrInvalidCredentials),
		errors.Is(err, ErrTokenInvalid),
		errors.Is(err, ErrTokenExpired),
		errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized

	// 404
	case errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrRoleNotFound),
		errors.Is(err, ErrAirportNotFound),
		errors.Is(err, ErrAircraftNotFound),
		errors.Is(err, ErrScheduleNotFound),
		errors.Is(err, ErrFlightNotFound),
		errors.Is(err, ErrFlightSeatNotFound),
		errors.Is(err, ErrPNRNotFound),
		errors.Is(err, ErrPassengerNotFound),
		errors.Is(err, ErrSeatLockNotFound),
		errors.Is(err, ErrPaymentNotFound),
		errors.Is(err, ErrTicketNotFound),
		errors.Is(err, ErrBoardingPassNotFound),
		errors.Is(err, ErrBaggageNotFound),
		errors.Is(err, ErrSSRTypeNotFound),
		errors.Is(err, ErrMealNotFound):
		// errors.Is(err, ErrWalletNotFound),
		// errors.Is(err, ErrTopUpOrderNotFound),
		// errors.Is(err, ErrReceiverNotFound):
		return http.StatusNotFound

	// 409
	case errors.Is(err, ErrEmailAlreadyExists),
		errors.Is(err, ErrRoleNameDuplicate),
		errors.Is(err, ErrAirportCodeDuplicate),
		errors.Is(err, ErrFlightNumberDuplicate),
		errors.Is(err, ErrSeatAlreadyBooked),
		errors.Is(err, ErrSeatAlreadyLocked),
		errors.Is(err, ErrAlreadyCheckedIn),
		errors.Is(err, ErrBoardingPassExists),
		errors.Is(err, ErrTicketAlreadyIssued):
		// errors.Is(err, ErrTopUpOrderFinalized):
		return http.StatusConflict

	// 410
	case errors.Is(err, ErrPNRHoldExpired):
		// errors.Is(err, ErrTopUpOrderExpired)
		return http.StatusGone

	// 422
	case errors.Is(err, ErrSegmentChronologyInvalid),
		errors.Is(err, ErrRoundTripAirportMismatch),
		errors.Is(err, ErrSeatFlightMismatch),
		errors.Is(err, ErrFlightAlreadyDeparted),
		errors.Is(err, ErrPNRAlreadyCancelled),
		errors.Is(err, ErrPNRAlreadyTicketed),
		errors.Is(err, ErrPaymentNotPending),
		errors.Is(err, ErrPaymentNotSuccess),
		errors.Is(err, ErrSystemRoleDelete),
		errors.Is(err, ErrCheckinTooEarly),
		errors.Is(err, ErrCheckinClosed),
		errors.Is(err, ErrTicketRequiredCheckin),
		errors.Is(err, ErrCheckinRequiredBoarding),
		errors.Is(err, ErrCheckinRequiredBaggage):
		// errors.Is(err, ErrTransferSelf),
		// errors.Is(err, ErrInsufficientBalance):
		return http.StatusUnprocessableEntity

	default:
		return http.StatusInternalServerError
	}
}

// errCode returns a short machine-readable error code.
func errCode(err error) string {
	switch {
	// case errors.Is(err, ErrInsufficientBalance):
	// 	return "INSUFFICIENT_BALANCE"
	case errors.Is(err, ErrSeatAlreadyLocked):
		return "SEAT_LOCKED"
	case errors.Is(err, ErrSeatAlreadyBooked):
		return "SEAT_BOOKED"
	case errors.Is(err, ErrPNRHoldExpired):
		return "PNR_EXPIRED"
	case errors.Is(err, ErrInvalidCredentials):
		return "INVALID_CREDENTIALS"
	case errors.Is(err, ErrTokenInvalid), errors.Is(err, ErrTokenExpired):
		return "TOKEN_INVALID"
	default:
		return ""
	}
}
