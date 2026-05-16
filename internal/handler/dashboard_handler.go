package handler

import (
	"net/http"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
	"strconv"
)

type DashboardHandler struct {
	svc services.DashboardService
}

func NewDashboardHandler(
	svc services.DashboardService,
) *DashboardHandler {
	return &DashboardHandler{
		svc: svc,
	}
}

// Summary godoc
// @Summary      Dashboard summary
// @Description  Mendapatkan statistik utama dashboard PSS
// @Tags         Dashboard
// @Produce      json
// @Success      200  {object}  utils.Response{data=dto.DashboardSummaryResponse}
// @Failure      500  {object}  utils.Response
// @Security     BearerAuth
// @Router       /dashboard/summary [get]
func (h *DashboardHandler) Summary(
	w http.ResponseWriter,
	r *http.Request,
) {

	result, err := h.svc.GetSummary(r.Context())
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		"success",
		result,
	)
}

// RevenueTrend godoc
// @Summary      Revenue trend
// @Description  Mendapatkan trend revenue harian
// @Tags         Dashboard
// @Produce      json
// @Param        days  query     int  false  "Jumlah hari" default(7)
// @Success      200   {object}  utils.Response{data=[]dto.RevenueTrendResponse}
// @Failure      500   {object}  utils.Response
// @Security     BearerAuth
// @Router       /dashboard/revenue-trend [get]
func (h *DashboardHandler) RevenueTrend(
	w http.ResponseWriter,
	r *http.Request,
) {

	days := 7

	if q := r.URL.Query().Get("days"); q != "" {
		if parsed, err := strconv.Atoi(q); err == nil {
			days = parsed
		}
	}

	result, err := h.svc.GetRevenueTrend(
		r.Context(),
		days,
	)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		"success",
		result,
	)
}

// BookingStatus godoc
// @Summary      Booking status distribution
// @Description  Mendapatkan distribusi status booking
// @Tags         Dashboard
// @Produce      json
// @Success      200  {object}  utils.Response{data=[]dto.BookingStatusResponse}
// @Failure      500  {object}  utils.Response
// @Security     BearerAuth
// @Router       /dashboard/booking-status [get]
func (h *DashboardHandler) BookingStatus(
	w http.ResponseWriter,
	r *http.Request,
) {

	result, err := h.svc.GetBookingStatus(
		r.Context(),
	)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		"success",
		result,
	)
}

// TodayFlights godoc
// @Summary      Today's flights
// @Description  Mendapatkan daftar penerbangan hari ini
// @Tags         Dashboard
// @Produce      json
// @Success      200  {object}  utils.Response{data=[]dto.TodayFlightResponse}
// @Failure      500  {object}  utils.Response
// @Security     BearerAuth
// @Router       /dashboard/today-flights [get]
func (h *DashboardHandler) TodayFlights(
	w http.ResponseWriter,
	r *http.Request,
) {

	result, err := h.svc.GetTodayFlights(
		r.Context(),
	)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		"success",
		result,
	)
}

// RecentBookings godoc
// @Summary      Recent bookings
// @Description  Mendapatkan booking terbaru
// @Tags         Dashboard
// @Produce      json
// @Param        limit  query     int  false  "Limit data" default(10)
// @Success      200    {object}  utils.Response{data=[]dto.RecentBookingResponse}
// @Failure      500    {object}  utils.Response
// @Security     BearerAuth
// @Router       /dashboard/recent-bookings [get]
func (h *DashboardHandler) RecentBookings(
	w http.ResponseWriter,
	r *http.Request,
) {

	limit := 10

	if q := r.URL.Query().Get("limit"); q != "" {
		if parsed, err := strconv.Atoi(q); err == nil {
			limit = parsed
		}
	}

	result, err := h.svc.GetRecentBookings(
		r.Context(),
		limit,
	)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		"success",
		result,
	)
}
