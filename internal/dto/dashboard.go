package dto

type DashboardSummaryResponse struct {
	TotalBookings   int64   `json:"total_bookings"`
	TotalPassengers int64   `json:"total_passengers"`
	TodayFlights    int64   `json:"today_flights"`
	TotalRevenue    float64 `json:"total_revenue"`
}

type RevenueTrendResponse struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
}

type BookingStatusResponse struct {
	Status string `json:"status"`
	Value  int64  `json:"value"`
}

type TodayFlightResponse struct {
	ID             uint   `json:"id"`
	FlightNumber   string `json:"flight_number"`
	Origin         string `json:"origin"`
	Destination    string `json:"destination"`
	DepartureTime  string `json:"departure_time"`
	Aircraft       string `json:"aircraft"`
	PassengerCount int64  `json:"passenger_count"`
	Status         string `json:"status"`
}

type RecentBookingResponse struct {
	ID            uint   `json:"id"`
	BookingCode   string `json:"booking_code"`
	PassengerName string `json:"passenger_name"`
	Route         string `json:"route"`
	PaymentStatus string `json:"payment_status"`
}

type OperationalAlertResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
}
