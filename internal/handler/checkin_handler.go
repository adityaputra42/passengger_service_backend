package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
)

type CheckinHandler struct {
	svc services.CheckinService
}

func NewCheckinHandler(svc services.CheckinService) *CheckinHandler {
	return &CheckinHandler{svc: svc}
}

// Checkin godoc
// @Summary      Lakukan check-in
// @Description  Melakukan check-in penumpang untuk segment penerbangan. Check-in dibuka 24 jam sebelum dan ditutup 45 menit sebelum keberangkatan. Tiket harus sudah diterbitkan.
// @Tags         CheckIn
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CheckinRequest  true  "Data check-in"
// @Success      200      {object}  utils.Response{data=dto.CheckinResultResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Failure      409      {object}  utils.Response  "Sudah check-in"
// @Failure      422      {object}  utils.Response  "Di luar window check-in atau tiket belum ada"
// @Security     BearerAuth
// @Router       /checkin [post]
func (h *CheckinHandler) Checkin(w http.ResponseWriter, r *http.Request) {
	var req dto.CheckinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	result, err := h.svc.Checkin(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToCheckinResultResponse(result))

}

// GetByPassenger godoc
// @Summary      Riwayat check-in penumpang
// @Description  Mengambil seluruh riwayat check-in seorang penumpang.
// @Tags         CheckIn
// @Produce      json
// @Param        passengerID  path      string  true  "Passenger UUID"
// @Success      200          {object}  utils.Response{data=[]dto.CheckinResponse}
// @Failure      400          {object}  utils.Response
// @Failure      401          {object}  utils.Response
// @Security     BearerAuth
// @Router       /checkin/passenger/{passengerID} [get]
func (h *CheckinHandler) GetByPassenger(w http.ResponseWriter, r *http.Request) {
	passengerID, ok := utils.UUIDParam(w, r, "passengerID")
	if !ok {
		return
	}
	checkins, err := h.svc.GetByPassenger(r.Context(), passengerID)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	out := make([]dto.CheckinResponse, 0, len(checkins))
	for i := range checkins {
		if resp := dto.ToCheckinResponse(&checkins[i]); resp != nil {
			out = append(out, *resp)
		}
	}

	utils.WriteJSON(w, http.StatusOK, "success", out)
}
