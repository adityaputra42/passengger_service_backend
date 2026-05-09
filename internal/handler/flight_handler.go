package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
	"time"
)

type FlightHandler struct {
	svc services.FlightService
}

func NewFlightHandler(svc services.FlightService) *FlightHandler {
	return &FlightHandler{svc: svc}
}

// Search godoc
// @Summary      Cari penerbangan tersedia
// @Description  Mencari penerbangan yang tersedia berdasarkan rute dan tanggal. Endpoint ini bersifat publik.
// @Tags         Flight
// @Produce      json
// @Param        dep   query     string  true  "Kode IATA keberangkatan"  example(CGK)
// @Param        arr   query     string  true  "Kode IATA tujuan"          example(DPS)
// @Param        date  query     string  true  "Tanggal (YYYY-MM-DD)"      example(2026-05-15)
// @Success      200   {object}  utils.Response{data=dto.FlightListResponse}
// @Failure      400   {object}  utils.Response
// @Failure      404   {object}  utils.Response  "Airport tidak ditemukan"
// @Router       /flights/search [get]
func (h *FlightHandler) Search(w http.ResponseWriter, r *http.Request) {
	dep := r.URL.Query().Get("dep")
	arr := r.URL.Query().Get("arr")
	dateStr := r.URL.Query().Get("date")

	if dep == "" || arr == "" || dateStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "parameter dep, arr, dan date wajib diisi", fmt.Errorf("parameter dep, arr, dan date wajib diisi"))
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "format date tidak valid, gunakan YYYY-MM-DD", fmt.Errorf("format date tidak valid, gunakan YYYY-MM-DD"))

		return
	}

	results, err := h.svc.Search(r.Context(), dto.SearchFlightRequest{
		DepartureCode: dep,
		ArrivalCode:   arr,
		Date:          date,
	})
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightSearchResponseList(results))

}

// Get godoc
// @Summary      Detail penerbangan
// @Description  Mengambil detail satu penerbangan beserta jadwal dan aircraft. Endpoint ini bersifat publik.
// @Tags         Flight
// @Produce      json
// @Param        id   path      string  true  "Flight UUID"
// @Success      200  {object}  utils.Response{data=dto.FlightResponse}
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /flights/{id} [get]
func (h *FlightHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	flight, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightResponse(flight))
}

// SeatMap godoc
// @Summary      Peta kursi penerbangan
// @Description  Mengambil seluruh peta kursi penerbangan beserta status ketersediaan dan harga. Membutuhkan login.
// @Tags         Flight
// @Produce      json
// @Param        id   path      string  true  "Flight UUID"
// @Success      200  {object}  utils.Response{data=[]dto.FlightSeatResult}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /flights/{id}/seat-map [get]
func (h *FlightHandler) SeatMap(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	seats, err := h.svc.GetSeatMap(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", seats)
}

// Generate godoc
// @Summary      Generate penerbangan dari jadwal
// @Description  Membuat instance penerbangan dari sebuah jadwal untuk rentang tanggal tertentu. Hanya admin.
// @Tags         Flight
// @Accept       json
// @Produce      json
// @Param        request  body      object{schedule_id=string,from=string,to=string}  true  "Parameter generate"
// @Success      200      {object}  utils.Response{data=object{generated=int}}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /flights/generate [post]
func (h *FlightHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ScheduleID string `json:"schedule_id"`
		From       string `json:"from"` // YYYY-MM-DD
		To         string `json:"to"`   // YYYY-MM-DD
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	schedID, err := utils.ParseUUIDStr(w, body.ScheduleID, "schedule_id")
	if err != nil {
		return
	}
	from, err := time.Parse("2006-01-02", body.From)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "format from tidak valid", fmt.Errorf("format from tidak valid"))
		return
	}
	to, err := time.Parse("2006-01-02", body.To)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "format from tidak valid", fmt.Errorf("format from tidak valid"))
		return
	}
	count, err := h.svc.GenerateFromSchedule(r.Context(), schedID, from, to)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", map[string]int{"generated": count})

}

// UpdateStatus godoc
// @Summary      Update status penerbangan
// @Description  Mengubah status penerbangan (scheduled, boarding, departed, arrived, cancelled, delayed). Hanya admin.
// @Tags         Flight
// @Accept       json
// @Produce      json
// @Param        id       path      string                  true  "Flight UUID"
// @Param        request  body      object{status=string}   true  "Status baru"
// @Success      200      {object}  utils.Response
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /flights/{id}/status [patch]
func (h *FlightHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	err := h.svc.UpdateStatus(r.Context(), id, models.FlightStatus(body.Status))
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", nil)

}
