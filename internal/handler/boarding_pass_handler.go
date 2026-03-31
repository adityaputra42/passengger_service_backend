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

// POST /boarding-passes  [agent]
func (h *BoardingPassHandler) Issue(w http.ResponseWriter, r *http.Request) {
	var req dto.IssueBoardingPassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	bp, err := h.svc.Issue(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBoardingPassResponse(bp))

}

// GET /boarding-passes/passenger/{passengerID}/segment/{segmentID}  [AuthRequired]
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBoardingPassResponse(bp))
}

// GET /boarding-passes/segment/{segmentID}  [agent]
func (h *BoardingPassHandler) GetBySegment(w http.ResponseWriter, r *http.Request) {
	segmentID, ok := utils.UUIDParam(w, r, "segmentID")
	if !ok {
		return
	}
	bps, err := h.svc.GetBySegment(r.Context(), segmentID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToBoardingPassResponseList(bps))

}
