package routes

import (
	"net/http"
	"passenger_service_backend/internal/config"
	"passenger_service_backend/internal/injection"
	"passenger_service_backend/internal/middleware"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type Dependencies struct {
	RBACService services.RBACService
	UserService services.UserService
	JWTService  *utils.JWTService
}

func SetupRoutes(handler *injection.Handler, logger *zap.Logger, cfg config.CORSConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		chimiddleware.RequestID,
		chimiddleware.RealIP,
		middleware.Logger(logger),
		middleware.Recovery(logger),
	)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chimiddleware.AllowContentType("application/json", "multipart/form-data"))

	buildDependencies(handler)

	return r
}

func buildDependencies(handler *injection.Handler) Dependencies {
	return Dependencies{
		RBACService: handler.RBACService,
		UserService: handler.UserService,
		JWTService:  handler.JWTService,
	}
}
