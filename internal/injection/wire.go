//go:build wireinject
// +build wireinject

package injection

import (
	"context"
	"passenger_service_backend/internal/config"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/handler"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProvideDB provides database connection
func ProvideDB() *gorm.DB {
	return db.DB
}

// ProvideJWTService provides JWT service instance with config
func ProvideJWTService(config *config.Config) *utils.JWTService {
	return utils.NewJWTService(config)
}

// Repository Providers
var repositorySet = wire.NewSet(
	ProvideDB,

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
)

// Service Providers
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

// Utils Providers
var utilsSet = wire.NewSet(
	ProvideJWTService,
)

// Handler Providers
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

// InitializeAllHandler initializes all handler with config
func InitializeAllHandler(config *config.Config, ctx context.Context) *Handler {
	wire.Build(
		repositorySet,
		serviceSet,
		utilsSet,
		handlerSet,
		NewHandler,
	)
	return &Handler{}
}

// Handler struct contains all handler and services
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

// NewHandler creates new Handler instance
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
