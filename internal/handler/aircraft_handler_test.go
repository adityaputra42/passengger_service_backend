package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/mocks"
	"passenger_service_backend/internal/models"
)

func TestAircraftHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAircraftService(ctrl)
	h := handler.NewAircraftHandler(mockSvc)

	reqDto := dto.CreateAircraftRequest{
		Model:        "B737",
		Manufacturer: "Boeing",
	}

	body, _ := json.Marshal(reqDto)
	req := httptest.NewRequest(http.MethodPost, "/aircraft", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.Aircraft{Model: "B737"}, nil)

	h.Create(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "B737")
}

func TestAircraftHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAircraftService(ctrl)
	h := handler.NewAircraftHandler(mockSvc)

	id := uuid.New()
	targetAircraft := &models.Aircraft{ID: id, Model: "B737"}

	req := httptest.NewRequest(http.MethodGet, "/aircraft/"+id.String(), nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().GetByID(gomock.Any(), id).Return(targetAircraft, nil)

	h.Get(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "B737")
}

func TestAircraftHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAircraftService(ctrl)
	h := handler.NewAircraftHandler(mockSvc)

	aircrafts := []models.Aircraft{{Model: "B737"}, {Model: "A320"}}

	req := httptest.NewRequest(http.MethodGet, "/aircraft", nil)
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().GetAll(gomock.Any()).Return(aircrafts, nil)

	h.List(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "B737")
	assert.Contains(t, rec.Body.String(), "A320")
}
