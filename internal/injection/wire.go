//go:build wireinject
// +build wireinject

package injection

// ─────────────────────────────────────────────────────────────
// This file shows the CHANGES needed to wire.go.
// It is NOT a complete replacement — merge with your existing wire.go.
// ─────────────────────────────────────────────────────────────
//
// STEP 1: Add ProvideCache and ProvideRedisConfig to wire.go

import (
	"context"
	"passenger_service_backend/internal/cache"
	"passenger_service_backend/internal/config"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func ProvideDB() *gorm.DB {
	return db.DB
}

func ProvideJWTService(cfg *config.Config) *utils.JWTService {
	return utils.NewJWTService(cfg)
}

func ProvideCache(cfg *config.Config) (*cache.Client, error) {
	return cache.New(cache.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}

var repositorySet = wire.NewSet(
	ProvideDB,

	// ── Base repos ──
	repository.NewAircraftRepository,
	repository.NewAirportRepository,
	repository.NewAircraftSeatRepository,
	repository.NewBaggageRepository,
	repository.NewBoardingPassRepository,
	repository.NewCheckinRepository,
	repository.NewPaymentRepository,
	repository.NewPermissionRepository,
	repository.NewRBACRepository,
	repository.NewFlightRepository,
	repository.NewFlightSeatRepository,
	repository.NewFlightScheduleRepository,
	repository.NewMealRepository,
	repository.NewUserReposiory,
	repository.NewPassengerMealRepository,
	repository.NewPassengerSSRRepository,
	repository.NewRoleRepository,
	repository.NewPNRContactRepository,
	repository.NewPNRPassengerRepository,
	repository.NewPNRRepository,
	repository.NewPNRSegmentRepository,
	repository.NewSeatAssignmentRepository,
	repository.NewSeatClassRepository,
	repository.NewSeatLockRepository,
	repository.NewSSRTypeRepository,
	repository.NewTicketRepository,
	repository.NewTicketSegmentRepository,

	// ── Cached wrappers ──
	repository.NewCachedAircraftRepository,
	repository.NewCachedAirportRepository,
	repository.NewCachedRBACRepository,
	repository.NewCachedFlightRepository,
	repository.NewCachedFlightSeatRepository,
	repository.NewCachedFlightScheduleRepository,
	repository.NewCachedPNRRepository,

	// ── Bind HANYA untuk repo yang punya cached wrapper ──
	wire.Bind(new(repository.AircraftRepository), new(*repository.CachedAircraftRepository)),
	wire.Bind(new(repository.AirportRepository), new(*repository.CachedAirportRepository)),
	wire.Bind(new(repository.RBACRepository), new(*repository.CachedRBACRepository)),
	wire.Bind(new(repository.FlightRepository), new(*repository.CachedFlightRepository)),
	wire.Bind(new(repository.FlightSeatRepository), new(*repository.CachedFlightSeatRepository)),
	wire.Bind(new(repository.FlightScheduleRepository), new(*repository.CachedFlightScheduleRepository)),
	wire.Bind(new(repository.PNRRepository), new(*repository.CachedPNRRepository)),
)

var utilsSet = wire.NewSet(
	ProvideJWTService,
	ProvideCache, // ← ADD THIS
)

var serviceSet = wire.NewSet(
	services.NewAircraftService,
	services.NewAirportService,
	services.NewAuthService,
	services.NewBaggageService,
	services.NewBoardingPassService,
	services.NewBookingService,
	services.NewCheckinService,
	services.NewFlightScheduleService,
	services.NewFlightService,
	services.NewPaymentService,
	services.NewRBACService,
	services.NewRoleService,
	services.NewSeatLockService,
	services.NewTicketService,
	services.NewUserService,
)

var handlerSet = wire.NewSet(
	handler.NewAircraftHandler,
	handler.NewAirportHandler,
	handler.NewAuthHandler,
	handler.NewBaggageHandler,
	handler.NewBoardingPassHandler,
	handler.NewBookingHandler,
	handler.NewCheckinHandler,
	handler.NewFlightHandler,
	handler.NewFlightScheduleHandler,
	handler.NewPaymentHandler,
	handler.NewUserHandler,
	handler.NewRoleHandler,
)

func InitializeAllHandler(cfg *config.Config, ctx context.Context) (*Handler, error) {
	wire.Build(
		repositorySet,
		serviceSet,
		utilsSet,
		handlerSet,
		NewHandler,
	)
	return &Handler{}, nil
}

type Handler struct {
	AircraftHandler       *handler.AircraftHandler
	AirportHandler        *handler.AirportHandler
	AuthHandler           *handler.AuthHandler
	BaggageHandler        *handler.BaggageHandler
	BoardingPassHandler   *handler.BoardingPassHandler
	BookingHandler        *handler.BookingHandler
	CheckinHandler        *handler.CheckinHandler
	FlightHandler         *handler.FlightHandler
	FlightScheduleHandler *handler.FlightScheduleHandler
	PaymentHandler        *handler.PaymentHandler
	RoleHandler           *handler.RoleHandler
	UserHandler           *handler.UserHandler

	RBACService services.RBACService
	UserService services.UserService
	JWTService  *utils.JWTService
}

func NewHandler(
	aircraftHandler *handler.AircraftHandler,
	airportHandler *handler.AirportHandler,
	authHandler *handler.AuthHandler,
	baggageHandler *handler.BaggageHandler,
	boardingPassHandler *handler.BoardingPassHandler,
	bookingHandler *handler.BookingHandler,
	checkinHandler *handler.CheckinHandler,
	flightHandler *handler.FlightHandler,
	flightScheduleHandler *handler.FlightScheduleHandler,
	paymentHandler *handler.PaymentHandler,
	roleHandler *handler.RoleHandler,
	userHandler *handler.UserHandler,
	rbacService services.RBACService,
	userService services.UserService,
	jwtService *utils.JWTService,
) *Handler {
	return &Handler{
		AircraftHandler:       aircraftHandler,
		AirportHandler:        airportHandler,
		AuthHandler:           authHandler,
		BaggageHandler:        baggageHandler,
		BoardingPassHandler:   boardingPassHandler,
		BookingHandler:        bookingHandler,
		CheckinHandler:        checkinHandler,
		FlightHandler:         flightHandler,
		FlightScheduleHandler: flightScheduleHandler,
		PaymentHandler:        paymentHandler,
		UserHandler:           userHandler,
		RoleHandler:           roleHandler,
		RBACService:           rbacService,
		UserService:           userService,
		JWTService:            jwtService,
	}
}
