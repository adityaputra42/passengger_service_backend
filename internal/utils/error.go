package utils

import "errors"

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
	ErrFlightNotFound      = errors.New("penerbangan tidak ditemukan")
	ErrFlightNotScheduled  = errors.New("status penerbangan bukan scheduled")
	ErrFlightAlreadyDeparted = errors.New("penerbangan sudah berangkat")
	ErrNoFlightsGenerated  = errors.New("tidak ada penerbangan yang bisa di-generate")

	// FlightSeat
	ErrFlightSeatNotFound   = errors.New("kursi tidak ditemukan")
	ErrSeatNotAvailable     = errors.New("kursi tidak tersedia")
	ErrSeatAlreadyLocked    = errors.New("kursi sedang dipesan oleh pengguna lain")
	ErrSeatAlreadyBooked    = errors.New("kursi sudah dipesan")

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
	ErrPaymentNotFound     = errors.New("pembayaran tidak ditemukan")
	ErrPaymentNotPending   = errors.New("pembayaran tidak dalam status pending")
	ErrPaymentAlreadyPaid  = errors.New("pembayaran sudah berhasil diproses")
	ErrPaymentNotSuccess   = errors.New("pembayaran belum berhasil")
	ErrPNRPaymentMismatch  = errors.New("pembayaran tidak terkait dengan PNR ini")

	// Ticket
	ErrTicketNotFound        = errors.New("tiket tidak ditemukan")
	ErrTicketAlreadyIssued   = errors.New("tiket sudah diterbitkan untuk penumpang ini")

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
	ErrBaggageNotFound       = errors.New("bagasi tidak ditemukan")
	ErrCheckinRequiredBaggage = errors.New("harus check-in sebelum mendaftarkan bagasi")

	// SSR / Meal
	ErrSSRTypeNotFound  = errors.New("tipe SSR tidak ditemukan")
	ErrMealNotFound     = errors.New("pilihan meal tidak ditemukan")
	ErrSSRNotFound      = errors.New("SSR penumpang tidak ditemukan")
	ErrMealSSRNotFound  = errors.New("meal penumpang tidak ditemukan")

	// Multi-segment / Trip validation
	ErrInvalidSegmentCount      = errors.New("jumlah segment tidak sesuai trip_type")
	ErrSegmentChronologyInvalid = errors.New("urutan waktu segment tidak valid")
	ErrRoundTripAirportMismatch = errors.New("airport return tidak cocok dengan arrival outbound")
	ErrSeatFlightMismatch       = errors.New("seat tidak milik penerbangan yang dipilih")
	ErrDuplicatePassengerSeat   = errors.New("penumpang memilih kursi dua kali di segment yang sama")
	ErrDuplicateSeatSelection   = errors.New("kursi yang sama dipilih dua penumpang berbeda")
	ErrInvalidPassengerIndex    = errors.New("passenger_index tidak valid")
)
