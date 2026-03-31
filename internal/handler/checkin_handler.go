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

// POST /checkin  [AuthRequired]
func (h *CheckinHandler) Checkin(w http.ResponseWriter, r *http.Request) {
	var req dto.CheckinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	result, err := h.svc.Checkin(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToCheckinResultResponse(result))

}

// GET /checkin/passenger/{passengerID}  [AuthRequired]
func (h *CheckinHandler) GetByPassenger(w http.ResponseWriter, r *http.Request) {
	passengerID, ok := utils.UUIDParam(w, r, "passengerID")
	if !ok {
		return
	}
	checkins, err := h.svc.GetByPassenger(r.Context(), passengerID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
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
