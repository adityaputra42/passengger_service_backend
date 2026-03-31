package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
)

type AircraftHandler struct {
	svc services.AircraftService
}

func NewAircraftHandler(svc services.AircraftService) *AircraftHandler {
	return &AircraftHandler{svc: svc}
}

// POST /aircraft  [admin]
func (h *AircraftHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAircraftRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	aircraft, err := h.svc.Create(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAircraftResponse(aircraft))

}

// GET /aircraft
func (h *AircraftHandler) List(w http.ResponseWriter, r *http.Request) {
	aircrafts, err := h.svc.GetAll(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	out := make([]dto.AircraftResponse, 0, len(aircrafts))
	for i := range aircrafts {
		if resp := dto.ToAircraftResponse(&aircrafts[i]); resp != nil {
			out = append(out, *resp)
		}
	}
	utils.WriteJSON(w, http.StatusOK, "success", out)
}

// GET /aircraft/{id}
func (h *AircraftHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	aircraft, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAircraftResponse(aircraft))
}

// GET /aircraft/{id}/seats
func (h *AircraftHandler) GetWithSeats(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	aircraft, err := h.svc.GetWithSeats(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAircraftResponse(aircraft))
}

// PUT /aircraft/{id}  [admin]
func (h *AircraftHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	var req dto.UpdateAircraftRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	aircraft, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAircraftResponse(aircraft))
}

// DELETE /aircraft/{id}  [admin]
func (h *AircraftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", nil)
}

// POST /aircraft/{id}/seats/generate  [admin]
func (h *AircraftHandler) GenerateSeats(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	var req dto.GenerateSeatsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	seats, err := h.svc.GenerateSeats(r.Context(), id, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	out := make([]dto.AircraftSeatResponse, 0, len(seats))
	for i := range seats {
		if resp := dto.ToAircraftSeatResponse(&seats[i]); resp != nil {
			out = append(out, *resp)
		}
	}

	utils.WriteJSON(w, http.StatusOK, "success", out)
}
