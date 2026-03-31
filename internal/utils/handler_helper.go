package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func WriteError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}
func ChiParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
func ParseUUIDStr(w http.ResponseWriter, raw, field string) (uuid.UUID, error) {
	if raw == "" {
		err := fmt.Errorf("%s wajib diisi", field)
		WriteError(w, http.StatusBadRequest, err.Error(), err)

		return uuid.Nil, err
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		msg := fmt.Sprintf("%s harus berupa UUID yang valid", field)
		WriteError(w, http.StatusBadRequest, msg, fmt.Errorf("%s", msg))

		return uuid.Nil, fmt.Errorf("%s", msg)
	}
	return id, nil
}
