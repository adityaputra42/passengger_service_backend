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


// Add godoc
// @Summary      Daftarkan bagasi penumpang
// @Description  Menambahkan bagasi baru untuk penumpang. Penumpang harus sudah check-in terlebih dahulu.
// @Tags         Baggage
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AddBaggageRequest  true  "Data bagasi"
// @Success      200      {object}  utils.Response{data=dto.BaggageResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      404      {object}  utils.Response  "Penumpang tidak ditemukan"
// @Failure      422      {object}  utils.Response  "Penumpang belum check-in"
// @Security     BearerAuth
// @Router       /baggage [post]
func (h *BaggageHandler) Add(w http.ResponseWriter, r *http.Request) {
	var req dto.AddBaggageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	baggage, err := h.svc.Add(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBaggageResponse(baggage))

}

// UpdateStatus godoc
// @Summary      Update status bagasi
// @Description  Memperbarui status bagasi (checked_in, loaded, delivered, lost). Hanya agent ke atas.
// @Tags         Baggage
// @Accept       json
// @Produce      json
// @Param        id       path      string          true  "Baggage UUID"
// @Param        request  body      object{status=string}  true  "Status baru"
// @Success      200      {object}  utils.Response{data=dto.BaggageResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /baggage/{id}/status [put]
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
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBaggageResponse(baggage))

}

// GetByPassenger godoc
// @Summary      Daftar bagasi penumpang
// @Description  Mengambil semua bagasi milik seorang penumpang.
// @Tags         Baggage
// @Produce      json
// @Param        passengerID  path      string  true  "Passenger UUID"
// @Success      200          {object}  utils.Response{data=[]dto.BaggageResponse}
// @Failure      400          {object}  utils.Response
// @Failure      401          {object}  utils.Response
// @Security     BearerAuth
// @Router       /baggage/passenger/{passengerID} [get]
func (h *BaggageHandler) GetByPassenger(w http.ResponseWriter, r *http.Request) {
	passengerID, ok := utils.UUIDParam(w, r, "passengerID")
	if !ok {
		return
	}
	bags, err := h.svc.GetByPassenger(r.Context(), passengerID)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBaggageResponseList(bags))

}
