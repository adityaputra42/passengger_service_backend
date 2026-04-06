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

func TestUserHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	h := handler.NewUserHandler(mockSvc)

	reqDto := dto.CreateUserRequest{
		Email:    "test@example.com",
		FullName: "Test User",
		Password: "password123",
		RoleID:   1,
	}

	body, _ := json.Marshal(reqDto)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.User{Email: "test@example.com"}, nil)

	h.Create(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "test@example.com")
}

func TestUserHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	h := handler.NewUserHandler(mockSvc)

	uid := uuid.New()
	user := &models.User{UID: uid, Email: "test@example.com"}

	req := httptest.NewRequest(http.MethodGet, "/users/"+uid.String(), nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("uid", uid.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	
	rec := httptest.NewRecorder()

	mockSvc.EXPECT().GetByUID(gomock.Any(), uid).Return(user, nil)

	h.Get(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "test@example.com")
}
