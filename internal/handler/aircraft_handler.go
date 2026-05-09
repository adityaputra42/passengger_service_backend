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

// Create godoc
// @Summary      Tambah aircraft baru
// @Description  Membuat data aircraft baru. Hanya admin dan super_admin yang bisa mengakses endpoint ini.
// @Tags         Aircraft
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateAircraftRequest  true  "Data aircraft"
// @Success      200      {object}  utils.Response{data=dto.AircraftResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Security     BearerAuth
// @Router       /aircraft [post]
func (h *AircraftHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAircraftRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	aircraft, err := h.svc.Create(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAircraftResponse(aircraft))
}

// List godoc
// @Summary      Daftar semua aircraft
// @Description  Mengambil seluruh daftar aircraft yang tersedia. Endpoint ini bersifat publik.
// @Tags         Aircraft
// @Produce      json
// @Success      200  {object}  utils.Response{data=[]dto.AircraftResponse}
// @Failure      500  {object}  utils.Response
// @Router       /aircraft [get]
func (h *AircraftHandler) List(w http.ResponseWriter, r *http.Request) {
	aircrafts, err := h.svc.GetAll(r.Context())
	if err != nil {
		utils.WriteServiceError(w, err)
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

// Get godoc
// @Summary      Detail aircraft
// @Description  Mengambil detail satu aircraft berdasarkan ID.
// @Tags         Aircraft
// @Produce      json
// @Param        id   path      string  true  "Aircraft UUID"
// @Success      200  {object}  utils.Response{data=dto.AircraftResponse}
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /aircraft/{id} [get]
func (h *AircraftHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	aircraft, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAircraftResponse(aircraft))
}

// GetWithSeats godoc
// @Summary      Detail aircraft beserta konfigurasi kursi
// @Description  Mengambil detail aircraft lengkap dengan daftar semua kursi dan kelasnya.
// @Tags         Aircraft
// @Produce      json
// @Param        id   path      string  true  "Aircraft UUID"
// @Success      200  {object}  utils.Response{data=dto.AircraftResponse}
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /aircraft/{id}/seats [get]
func (h *AircraftHandler) GetWithSeats(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	aircraft, err := h.svc.GetWithSeats(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAircraftResponse(aircraft))
}

// Update godoc
// @Summary      Update data aircraft
// @Description  Memperbarui informasi aircraft. Hanya admin dan super_admin.
// @Tags         Aircraft
// @Accept       json
// @Produce      json
// @Param        id       path      string                     true  "Aircraft UUID"
// @Param        request  body      dto.UpdateAircraftRequest  true  "Data update"
// @Success      200      {object}  utils.Response{data=dto.AircraftResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /aircraft/{id} [put]
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
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAircraftResponse(aircraft))
}

// Delete godoc
// @Summary      Hapus aircraft
// @Description  Menghapus data aircraft. Hanya admin dan super_admin.
// @Tags         Aircraft
// @Produce      json
// @Param        id   path      string  true  "Aircraft UUID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /aircraft/{id} [delete]
func (h *AircraftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", nil)
}

// GenerateSeats godoc
// @Summary      Generate kursi aircraft
// @Description  Membuat konfigurasi kursi aircraft berdasarkan layout yang diberikan. Hanya admin dan super_admin.
// @Tags         Aircraft
// @Accept       json
// @Produce      json
// @Param        id       path      string                    true  "Aircraft UUID"
// @Param        request  body      dto.GenerateSeatsRequest  true  "Konfigurasi layout kursi"
// @Success      200      {object}  utils.Response{data=[]dto.AircraftSeatResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /aircraft/{id}/seats/generate [post]
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
		utils.WriteServiceError(w, err)
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
