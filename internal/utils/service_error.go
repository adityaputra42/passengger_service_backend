package utils

import (
	"errors"
	"net/http"
)

func WriteServiceError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	status := statusForService(err)
	WriteError(w, status, err.Error(), err)
}

func statusForService(err error) int {
	switch {
	case errors.Is(err, ErrInvalidSegmentCount),
		errors.Is(err, ErrDuplicatePassengerSeat),
		errors.Is(err, ErrDuplicateSeatSelection),
		errors.Is(err, ErrInvalidPassengerIndex),
		errors.Is(err, ErrWrongPassword):
		return http.StatusBadRequest

	case errors.Is(err, ErrInvalidCredentials),
		errors.Is(err, ErrTokenInvalid),
		errors.Is(err, ErrTokenExpired),
		errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized

	case errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrRoleNotFound),
		errors.Is(err, ErrAirportNotFound),
		errors.Is(err, ErrAircraftNotFound),
		errors.Is(err, ErrScheduleNotFound),
		errors.Is(err, ErrFlightNotFound),
		errors.Is(err, ErrFlightSeatNotFound),
		errors.Is(err, ErrPNRNotFound),
		errors.Is(err, ErrPassengerNotFound),
		errors.Is(err, ErrSeatLockNotFound),
		errors.Is(err, ErrPaymentNotFound),
		errors.Is(err, ErrTicketNotFound),
		errors.Is(err, ErrBoardingPassNotFound),
		errors.Is(err, ErrBaggageNotFound),
		errors.Is(err, ErrSSRTypeNotFound),
		errors.Is(err, ErrMealNotFound):
		return http.StatusNotFound

	// 409 Conflict
	case errors.Is(err, ErrEmailAlreadyExists),
		errors.Is(err, ErrRoleNameDuplicate),
		errors.Is(err, ErrAirportCodeDuplicate),
		errors.Is(err, ErrFlightNumberDuplicate),
		errors.Is(err, ErrSeatAlreadyBooked),
		errors.Is(err, ErrSeatAlreadyLocked),
		errors.Is(err, ErrAlreadyCheckedIn),
		errors.Is(err, ErrBoardingPassExists),
		errors.Is(err, ErrTicketAlreadyIssued):
		return http.StatusConflict

	case errors.Is(err, ErrPNRHoldExpired):
		return http.StatusGone

	case errors.Is(err, ErrSegmentChronologyInvalid),
		errors.Is(err, ErrRoundTripAirportMismatch),
		errors.Is(err, ErrSeatFlightMismatch),
		errors.Is(err, ErrFlightAlreadyDeparted),
		errors.Is(err, ErrPNRAlreadyCancelled),
		errors.Is(err, ErrPNRAlreadyTicketed),
		errors.Is(err, ErrPaymentNotPending),
		errors.Is(err, ErrPaymentNotSuccess),
		errors.Is(err, ErrSystemRoleDelete),
		errors.Is(err, ErrCheckinTooEarly),
		errors.Is(err, ErrCheckinClosed),
		errors.Is(err, ErrTicketRequiredCheckin),
		errors.Is(err, ErrCheckinRequiredBoarding),
		errors.Is(err, ErrCheckinRequiredBaggage):
		return http.StatusUnprocessableEntity

	default:
		return http.StatusInternalServerError
	}
}
