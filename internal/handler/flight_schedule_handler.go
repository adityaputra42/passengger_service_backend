package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
)

type FlightScheduleHandler struct {
	svc services.FlightScheduleService
}

func NewFlightScheduleHandler(svc services.FlightScheduleService) *FlightScheduleHandler {
	return &FlightScheduleHandler{svc: svc}
}

// Create godoc
// @Summary      Buat jadwal penerbangan
// @Description  Membuat jadwal penerbangan baru (recurring schedule). Hanya admin.
// @Tags         FlightSchedule
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateFlightScheduleRequest  true  "Data jadwal"
// @Success      200      {object}  utils.Response{data=dto.FlightScheduleResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      409      {object}  utils.Response  "Nomor penerbangan sudah ada"
// @Security     BearerAuth
// @Router       /schedules [post]
func (h *FlightScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFlightScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	sched, err := h.svc.Create(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightScheduleResponse(sched))

}

// List godoc
// @Summary      Daftar jadwal penerbangan
// @Description  Mengambil seluruh jadwal penerbangan. Bisa filter rute dengan ?dep=CGK&arr=DPS. Endpoint publik.
// @Tags         FlightSchedule
// @Produce      json
// @Param        dep  query     string  false  "Kode IATA keberangkatan"
// @Param        arr  query     string  false  "Kode IATA tujuan"
// @Success      200  {object}  utils.Response{data=[]dto.FlightScheduleResponse}
// @Router       /schedules [get]
func (h *FlightScheduleHandler) List(w http.ResponseWriter, r *http.Request) {
	dep := r.URL.Query().Get("dep")
	arr := r.URL.Query().Get("arr")

	if dep != "" && arr != "" {
		scheds, err := h.svc.GetByRoute(r.Context(), dep, arr)
		if err != nil {
			utils.WriteServiceError(w, err)
			return
		}
		out := make([]dto.FlightScheduleResponse, 0, len(scheds))
		for i := range scheds {
			if resp := dto.ToFlightScheduleResponse(&scheds[i]); resp != nil {
				out = append(out, *resp)
			}
		}

		utils.WriteJSON(w, http.StatusOK, "success", out)

		return
	}

	scheds, err := h.svc.GetAll(r.Context())
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	out := make([]dto.FlightScheduleResponse, 0, len(scheds))
	for i := range scheds {
		if resp := dto.ToFlightScheduleResponse(&scheds[i]); resp != nil {
			out = append(out, *resp)
		}
	}
	utils.WriteJSON(w, http.StatusOK, "success", out)
}

// Get godoc
// @Summary      Detail jadwal penerbangan
// @Description  Mengambil detail satu jadwal penerbangan. Endpoint publik.
// @Tags         FlightSchedule
// @Produce      json
// @Param        id   path      string  true  "Schedule UUID"
// @Success      200  {object}  utils.Response{data=dto.FlightScheduleResponse}
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /schedules/{id} [get]
func (h *FlightScheduleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	sched, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightScheduleResponse(sched))

}

// Update godoc
// @Summary      Update jadwal penerbangan
// @Description  Memperbarui jadwal penerbangan (waktu, operating days). Hanya admin.
// @Tags         FlightSchedule
// @Accept       json
// @Produce      json
// @Param        id       path      string                          true  "Schedule UUID"
// @Param        request  body      dto.UpdateFlightScheduleRequest  true  "Data update"
// @Success      200      {object}  utils.Response{data=dto.FlightScheduleResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /schedules/{id} [put]
func (h *FlightScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	var req dto.UpdateFlightScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	sched, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightScheduleResponse(sched))

}

// Delete godoc
// @Summary      Hapus jadwal penerbangan
// @Description  Menghapus jadwal penerbangan. Hanya admin.
// @Tags         FlightSchedule
// @Produce      json
// @Param        id   path      string  true  "Schedule UUID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /schedules/{id} [delete]
func (h *FlightScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
