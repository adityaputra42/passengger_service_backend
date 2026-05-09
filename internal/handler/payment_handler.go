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

// Initiate godoc
// @Summary      Inisiasi pembayaran
// @Description  Memulai proses pembayaran untuk sebuah PNR. Total dihitung dari harga kursi yang di-assign. Status awal: pending.
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        request  body      dto.InitiatePaymentRequest  true  "Data pembayaran"
// @Success      200      {object}  utils.Response{data=dto.PaymentResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Failure      422      {object}  utils.Response  "PNR sudah dibatalkan atau sudah ticketed"
// @Security     BearerAuth
// @Router       /payments [post]
func (h *PaymentHandler) Initiate(w http.ResponseWriter, r *http.Request) {
	var req dto.InitiatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	payment, err := h.svc.Initiate(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPaymentResponse(payment))

}

// Confirm godoc
// @Summary      Konfirmasi pembayaran (webhook)
// @Description  Endpoint untuk payment gateway callback. Jika success=true, kursi diubah ke booked dan tiket diterbitkan. Hanya admin.
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        id       path      string                    true  "Payment UUID"
// @Param        request  body      object{success=boolean}   true  "Hasil pembayaran"
// @Success      200      {object}  utils.Response{data=dto.PaymentResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Failure      422      {object}  utils.Response  "Pembayaran tidak dalam status pending"
// @Security     BearerAuth
// @Router       /payments/{id}/confirm [post]
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
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPaymentResponse(payment))

}

// Refund godoc
// @Summary      Proses refund
// @Description  Memproses pengembalian uang untuk pembayaran yang sudah sukses. Hanya admin.
// @Tags         Payment
// @Produce      json
// @Param        id   path      string  true  "Payment UUID"
// @Success      200  {object}  utils.Response{data=dto.PaymentResponse}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      422  {object}  utils.Response  "Pembayaran belum success"
// @Security     BearerAuth
// @Router       /payments/{id}/refund [post]
func (h *PaymentHandler) Refund(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	payment, err := h.svc.Refund(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPaymentResponse(payment))
}

// ListByPNR godoc
// @Summary      Daftar pembayaran per PNR
// @Description  Mengambil semua riwayat pembayaran untuk sebuah PNR.
// @Tags         Payment
// @Produce      json
// @Param        pnr_id  query     string  true  "PNR UUID"
// @Success      200     {object}  utils.Response{data=[]dto.PaymentResponse}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Security     BearerAuth
// @Router       /payments [get]
func (h *PaymentHandler) ListByPNR(w http.ResponseWriter, r *http.Request) {
	pnrID, err := utils.ParseUUIDStr(w, r.URL.Query().Get("pnr_id"), "pnr_id")
	if err != nil {
		return
	}
	payments, err := h.svc.GetByPNR(r.Context(), pnrID)

	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPaymentResponseList(payments))

}
