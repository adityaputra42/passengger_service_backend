package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
)

// UUIDParam reads a chi URL param as uuid.UUID. Returns false and writes 400 if invalid.
func UUIDParam(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	raw := chi.URLParam(r, key)
	id, err := uuid.Parse(raw)
	if err != nil {
		WriteError(w, http.StatusBadRequest, key+" harus berupa UUID yang valid", fmt.Errorf("%s", key+" harus berupa UUID yang valid"))
		return uuid.Nil, false
	}
	return id, true
}

func PageLimit(r *http.Request) (page, limit int) {

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	p, _ := strconv.Atoi(pageStr)
	if p < 1 {
		p = defaultPage
	}

	l, _ := strconv.Atoi(limitStr)
	if l < 1 {
		l = defaultLimit
	}

	return p, l
}
