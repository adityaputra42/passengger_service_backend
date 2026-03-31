package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
)

type PaymentHandler struct {
	svc services.PaymentService
}

func NewPaymentHandler(svc services.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// POST /payments  [AuthRequired]
func (h *PaymentHandler) Initiate(w http.ResponseWriter, r *http.Request) {
	var req dto.InitiatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	payment, err := h.svc.Initiate(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPaymentResponse(payment))

}

// POST /payments/{id}/confirm  [internal/webhook]
func (h *PaymentHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	var body struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	payment, err := h.svc.Confirm(r.Context(), id, body.Success)

	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPaymentResponse(payment))

}

// POST /payments/{id}/refund  [admin]
func (h *PaymentHandler) Refund(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	payment, err := h.svc.Refund(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPaymentResponse(payment))
}

// GET /payments?pnr_id={uuid}  [AuthRequired]
func (h *PaymentHandler) ListByPNR(w http.ResponseWriter, r *http.Request) {
	pnrID, err := utils.ParseUUIDStr(w, r.URL.Query().Get("pnr_id"), "pnr_id")
	if err != nil {
		return
	}
	payments, err := h.svc.GetByPNR(r.Context(), pnrID)

	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPaymentResponseList(payments))

}
