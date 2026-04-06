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

func TestAirportHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAirportService(ctrl)
	h := handler.NewAirportHandler(mockSvc)

	reqDto := dto.CreateAirportRequest{
		Code:     "CGK",
		Name:     "Soekarno-Hatta",
		City:     "Jakarta",
		Country:  "Indonesia",
		Timezone: "Asia/Jakarta",
	}

	body, _ := json.Marshal(reqDto)
	req := httptest.NewRequest(http.MethodPost, "/airports", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.Airport{Code: "CGK"}, nil)

	h.Create(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "CGK")
}

func TestAirportHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAirportService(ctrl)
	h := handler.NewAirportHandler(mockSvc)

	id := uuid.New()
	airport := &models.Airport{ID: id, Code: "CGK"}

	req := httptest.NewRequest(http.MethodGet, "/airports/"+id.String(), nil)
	
	// Add chi route context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().GetByID(gomock.Any(), id).Return(airport, nil)

	h.Get(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "CGK")
}

func TestAirportHandler_GetByCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAirportService(ctrl)
	h := handler.NewAirportHandler(mockSvc)

	airport := &models.Airport{Code: "CGK"}

	req := httptest.NewRequest(http.MethodGet, "/airports/code/CGK", nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("code", "CGK")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().GetByCode(gomock.Any(), "CGK").Return(airport, nil)

	h.GetByCode(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "CGK")
}

func TestAirportHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAirportService(ctrl)
	h := handler.NewAirportHandler(mockSvc)

	airports := []models.Airport{{Code: "CGK"}, {Code: "DPS"}}

	req := httptest.NewRequest(http.MethodGet, "/airports", nil)
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().GetAll(gomock.Any()).Return(airports, nil)

	h.List(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "CGK")
}

func TestAirportHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAirportService(ctrl)
	h := handler.NewAirportHandler(mockSvc)

	id := uuid.New()
	reqDto := dto.UpdateAirportRequest{
		Name: "New Name",
	}

	body, _ := json.Marshal(reqDto)
	req := httptest.NewRequest(http.MethodPut, "/airports/"+id.String(), bytes.NewBuffer(body))
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().Update(gomock.Any(), id, gomock.Any()).Return(&models.Airport{Name: "New Name"}, nil)

	h.Update(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "New Name")
}

func TestAirportHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAirportService(ctrl)
	h := handler.NewAirportHandler(mockSvc)

	id := uuid.New()

	req := httptest.NewRequest(http.MethodDelete, "/airports/"+id.String(), nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().Delete(gomock.Any(), id).Return(nil)

	h.Delete(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
