package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
)

type BaggageHandler struct {
	svc services.BaggageService
}

func NewBaggageHandler(svc services.BaggageService) *BaggageHandler {
	return &BaggageHandler{svc: svc}
}

// POST /baggage  [AuthRequired]
func (h *BaggageHandler) Add(w http.ResponseWriter, r *http.Request) {
	var req dto.AddBaggageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	baggage, err := h.svc.Add(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBaggageResponse(baggage))

}

// PATCH /baggage/{id}/status  [agent/admin]
func (h *BaggageHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
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
	baggage, err := h.svc.UpdateStatus(r.Context(), id, models.BaggageStatus(body.Status))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBaggageResponse(baggage))

}

// GET /baggage/passenger/{passengerID}  [AuthRequired]
func (h *BaggageHandler) GetByPassenger(w http.ResponseWriter, r *http.Request) {
	passengerID, ok := utils.UUIDParam(w, r, "passengerID")
	if !ok {
		return
	}
	bags, err := h.svc.GetByPassenger(r.Context(), passengerID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBaggageResponseList(bags))

}
