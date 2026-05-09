package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
)

type AirportHandler struct {
	svc services.AirportService
}

func NewAirportHandler(svc services.AirportService) *AirportHandler {
	return &AirportHandler{svc: svc}
}

// Create godoc
// @Summary      Tambah airport baru
// @Description  Membuat data airport baru. Hanya admin dan super_admin.
// @Tags         Airport
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateAirportRequest  true  "Data airport"
// @Success      200      {object}  utils.Response{data=dto.AirportResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      409      {object}  utils.Response  "Kode IATA sudah terdaftar"
// @Security     BearerAuth
// @Router       /airport [post]
func (h *AirportHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAirportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	airport, err := h.svc.Create(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponse(airport))
}

// List godoc
// @Summary      Daftar semua airport
// @Description  Mengambil seluruh daftar airport. Bisa filter dengan query param ?q=jakarta untuk pencarian.
// @Tags         Airport
// @Produce      json
// @Param        q    query     string  false  "Kata kunci pencarian (nama, kota, kode, negara)"
// @Success      200  {object}  utils.Response{data=[]dto.AirportResponse}
// @Failure      401  {object}  utils.Response
// @Security     BearerAuth
// @Router       /airport [get]
func (h *AirportHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q != "" {
		airports, err := h.svc.Search(r.Context(), q)
		if err != nil {
			utils.WriteServiceError(w, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponseList(airports))
		return
	}
	airports, err := h.svc.GetAll(r.Context())
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponseList(airports))

}

// Get godoc
// @Summary      Detail airport
// @Description  Mengambil detail satu airport berdasarkan ID.
// @Tags         Airport
// @Produce      json
// @Param        id   path      string  true  "Airport UUID"
// @Success      200  {object}  utils.Response{data=dto.AirportResponse}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /airport/{id} [get]
func (h *AirportHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	airport, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponse(airport))
}

// GetByCode godoc
// @Summary      Cari airport berdasarkan kode IATA
// @Description  Mengambil detail airport menggunakan kode IATA 3 huruf (contoh: CGK, DPS).
// @Tags         Airport
// @Produce      json
// @Param        code  path      string  true  "Kode IATA (3 huruf)"  example(CGK)
// @Success      200   {object}  utils.Response{data=dto.AirportResponse}
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Security     BearerAuth
// @Router       /airport/code/{code} [get]
func (h *AirportHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := utils.ChiParam(r, "code")
	airport, err := h.svc.GetByCode(r.Context(), code)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponse(airport))
}

// Update godoc
// @Summary      Update data airport
// @Description  Memperbarui informasi airport. Hanya admin dan super_admin.
// @Tags         Airport
// @Accept       json
// @Produce      json
// @Param        id       path      string                    true  "Airport UUID"
// @Param        request  body      dto.UpdateAirportRequest  true  "Data update"
// @Success      200      {object}  utils.Response{data=dto.AirportResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /airport/{id} [put]
func (h *AirportHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	var req dto.UpdateAirportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	airport, err := h.svc.Update(r.Context(), id, req)

	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponse(airport))
}

// Delete godoc
// @Summary      Hapus airport
// @Description  Menghapus data airport. Hanya admin dan super_admin.
// @Tags         Airport
// @Produce      json
// @Param        id   path      string  true  "Airport UUID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /airport/{id} [delete]
func (h *AirportHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
