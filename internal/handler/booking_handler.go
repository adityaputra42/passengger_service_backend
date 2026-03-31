package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"

	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	svc services.BookingService
}

func NewBookingHandler(svc services.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

// POST /bookings  [AuthRequired]
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	pnr, err := h.svc.CreateBooking(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, "success", dto.ToPNRResponse(pnr))

}

// GET /bookings  [admin]
func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := utils.PageLimit(r)
	pnrs, total, err := h.svc.GetAllPNRs(r.Context(), page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPNRListResponse(pnrs, total, page, limit))

}

// GET /bookings/{id}  [AuthRequired]
func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	pnr, err := h.svc.GetPNRByID(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPNRResponse(pnr))
}

// GET /bookings/locator/{locator}  [AuthRequired]
func (h *BookingHandler) GetByLocator(w http.ResponseWriter, r *http.Request) {
	locator := chi.URLParam(r, "locator")
	pnr, err := h.svc.GetPNR(r.Context(), locator)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPNRResponse(pnr))
}

// PUT /bookings/{id}/contact  [AuthRequired]
func (h *BookingHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	var req dto.UpdateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	contact, err := h.svc.UpdateContact(r.Context(), id, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPNRContactResponse(contact))

}

// DELETE /bookings/{id}  [AuthRequired]
func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	err := h.svc.CancelPNR(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", nil)

}

// POST /bookings/passengers/{passengerID}/ssr  [AuthRequired]
func (h *BookingHandler) AddSSR(w http.ResponseWriter, r *http.Request) {
	passengerID, ok := utils.UUIDParam(w, r, "passengerID")
	if !ok {
		return
	}
	var req dto.AddSSRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	ssr, err := h.svc.AddSSR(r.Context(), passengerID, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPassengerSSRResponse(ssr))

}

// DELETE /bookings/ssr/{ssrID}  [AuthRequired]
func (h *BookingHandler) RemoveSSR(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "ssrID")
	if !ok {
		return
	}
	err := h.svc.RemoveSSR(r.Context(), id)

	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", nil)
}

// POST /bookings/passengers/{passengerID}/meal  [AuthRequired]
func (h *BookingHandler) AddMeal(w http.ResponseWriter, r *http.Request) {
	passengerID, ok := utils.UUIDParam(w, r, "passengerID")
	if !ok {
		return
	}
	var req dto.AddMealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	meal, err := h.svc.AddMeal(r.Context(), passengerID, req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPassengerMealResponse(meal))

}

// DELETE /bookings/meal/{mealID}  [AuthRequired]
func (h *BookingHandler) RemoveMeal(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "mealID")
	if !ok {
		return
	}
	err := h.svc.RemoveMeal(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", nil)
}
