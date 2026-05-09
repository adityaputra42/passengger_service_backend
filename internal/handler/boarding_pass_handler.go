package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"
)

type BoardingPassHandler struct {
	svc services.BoardingPassService
}

func NewBoardingPassHandler(svc services.BoardingPassService) *BoardingPassHandler {
	return &BoardingPassHandler{svc: svc}
}

// Issue godoc
// @Summary      Terbitkan boarding pass
// @Description  Menerbitkan boarding pass untuk penumpang yang sudah check-in. Hanya agent ke atas.
// @Tags         BoardingPass
// @Accept       json
// @Produce      json
// @Param        request  body      dto.IssueBoardingPassRequest  true  "Data boarding pass"
// @Success      200      {object}  utils.Response{data=dto.BoardingPassResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      403      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Failure      409      {object}  utils.Response  "Boarding pass sudah diterbitkan"
// @Failure      422      {object}  utils.Response  "Penumpang belum check-in"
// @Security     BearerAuth
// @Router       /boarding_passes [post]
func (h *BoardingPassHandler) Issue(w http.ResponseWriter, r *http.Request) {
	var req dto.IssueBoardingPassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	bp, err := h.svc.Issue(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBoardingPassResponse(bp))

}

// Get godoc
// @Summary      Detail boarding pass
// @Description  Mengambil boarding pass penumpang untuk segment tertentu.
// @Tags         BoardingPass
// @Produce      json
// @Param        passengerID  path      string  true  "Passenger UUID"
// @Param        segmentID    path      string  true  "Segment UUID"
// @Success      200          {object}  utils.Response{data=dto.BoardingPassResponse}
// @Failure      400          {object}  utils.Response
// @Failure      401          {object}  utils.Response
// @Failure      404          {object}  utils.Response
// @Security     BearerAuth
// @Router       /boarding_passes/passenger/{passengerID}/segment/{segmentID} [get]
func (h *BoardingPassHandler) Get(w http.ResponseWriter, r *http.Request) {
	passengerID, ok := utils.UUIDParam(w, r, "passengerID")
	if !ok {
		return
	}
	segmentID, ok := utils.UUIDParam(w, r, "segmentID")
	if !ok {
		return
	}
	bp, err := h.svc.GetByPassengerAndSegment(r.Context(), passengerID, segmentID)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBoardingPassResponse(bp))
}

// GetBySegment godoc
// @Summary      Daftar boarding pass per segment
// @Description  Mengambil semua boarding pass untuk satu segment penerbangan. Hanya agent ke atas.
// @Tags         BoardingPass
// @Produce      json
// @Param        segmentID  path      string  true  "Segment UUID"
// @Success      200        {object}  utils.Response{data=[]dto.BoardingPassResponse}
// @Failure      400        {object}  utils.Response
// @Failure      401        {object}  utils.Response
// @Failure      403        {object}  utils.Response
// @Security     BearerAuth
// @Router       /boarding_passes/segment/{segmentID} [get]
func (h *BoardingPassHandler) GetBySegment(w http.ResponseWriter, r *http.Request) {
	segmentID, ok := utils.UUIDParam(w, r, "segmentID")
	if !ok {
		return
	}
	bps, err := h.svc.GetBySegment(r.Context(), segmentID)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBoardingPassResponseList(bps))

}
