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

// POST /airports  [admin]
func (h *AirportHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAirportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	airport, err := h.svc.Create(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponse(airport))
}

// GET /airports  — ?q=jakarta untuk search
func (h *AirportHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q != "" {
		airports, err := h.svc.Search(r.Context(), q)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponseList(airports))
		return
	}
	airports, err := h.svc.GetAll(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponseList(airports))

}

// GET /airports/{id}
func (h *AirportHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	airport, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponse(airport))
}

// GET /airports/code/{code}
func (h *AirportHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := utils.ChiParam(r, "code")
	airport, err := h.svc.GetByCode(r.Context(), code)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponse(airport))
}

// PUT /airports/{id}  [admin]
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToAirportResponse(airport))
}

// DELETE /airports/{id}  [admin]
func (h *AirportHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}

	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", nil)
}
