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

// POST /schedules  [admin]
func (h *FlightScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFlightScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	sched, err := h.svc.Create(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightScheduleResponse(sched))

}

// GET /schedules  — ?dep=CGK&arr=DPS untuk filter rute
func (h *FlightScheduleHandler) List(w http.ResponseWriter, r *http.Request) {
	dep := r.URL.Query().Get("dep")
	arr := r.URL.Query().Get("arr")

	if dep != "" && arr != "" {
		scheds, err := h.svc.GetByRoute(r.Context(), dep, arr)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
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

// GET /schedules/{id}
func (h *FlightScheduleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	sched, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightScheduleResponse(sched))

}

// PUT /schedules/{id}  [admin]
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
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToFlightScheduleResponse(sched))

}

// DELETE /schedules/{id}  [admin]
func (h *FlightScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
