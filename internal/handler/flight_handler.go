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

// GET /flights/search?dep=CGK&arr=DPS&date=2026-04-01
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightSearchResponseList(results))

}

// GET /flights/{id}
func (h *FlightHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	flight, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightResponse(flight))
}

// GET /flights/{id}/seat-map
func (h *FlightHandler) SeatMap(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	seats, err := h.svc.GetSeatMap(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", seats)
}

// POST /flights/generate  [admin]  body: { schedule_id, from, to }
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", map[string]int{"generated": count})

}

// PATCH /flights/{id}/status  [admin]
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", nil)

}
