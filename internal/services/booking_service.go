package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
	"passenger_service_backend/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (*models.PNR, error)

	GetPNR(ctx context.Context, locator string) (*models.PNR, error)
	GetPNRByID(ctx context.Context, id uuid.UUID) (*models.PNR, error)
	GetAllPNRs(ctx context.Context, page, limit int) ([]models.PNR, int64, error)
	UpdateContact(ctx context.Context, pnrID uuid.UUID, req dto.UpdateContactRequest) (*models.PNRContact, error)
	AddSSR(ctx context.Context, passengerID uuid.UUID, req dto.AddSSRRequest) (*models.PassengerSSR, error)
	RemoveSSR(ctx context.Context, ssrID uuid.UUID) error
	AddMeal(ctx context.Context, passengerID uuid.UUID, req dto.AddMealRequest) (*models.PassengerMeal, error)
	RemoveMeal(ctx context.Context, mealID uuid.UUID) error
	CancelPNR(ctx context.Context, pnrID uuid.UUID) error
}

const (
	locatorChars    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	locatorLen      = 6
	seatLockTTL     = 30 * time.Minute
	maxLocatorRetry = 5
)

type bookingService struct {
	db                *gorm.DB
	pnrRepo           repository.PNRRepository
	contactRepo       repository.PNRContactRepository
	passengerRepo     repository.PNRPassengerRepository
	segmentRepo       repository.PNRSegmentRepository
	flightRepo        repository.FlightRepository
	flightSeatRepo    repository.FlightSeatRepository
	seatLockRepo      repository.SeatLockRepository
	seatAssignRepo    repository.SeatAssignmentRepository
	ssrTypeRepo       repository.SSRTypeRepository
	passengerSSRRepo  repository.PassengerSSRRepository
	mealRepo          repository.MealRepository
	passengerMealRepo repository.PassengerMealRepository
}

func NewBookingService(
	db *gorm.DB,
	pnrRepo repository.PNRRepository,
	contactRepo repository.PNRContactRepository,
	passengerRepo repository.PNRPassengerRepository,
	segmentRepo repository.PNRSegmentRepository,
	flightRepo repository.FlightRepository,
	flightSeatRepo repository.FlightSeatRepository,
	seatLockRepo repository.SeatLockRepository,
	seatAssignRepo repository.SeatAssignmentRepository,
	ssrTypeRepo repository.SSRTypeRepository,
	passengerSSRRepo repository.PassengerSSRRepository,
	mealRepo repository.MealRepository,
	passengerMealRepo repository.PassengerMealRepository,
) BookingService {
	return &bookingService{
		db:                db,
		pnrRepo:           pnrRepo,
		contactRepo:       contactRepo,
		passengerRepo:     passengerRepo,
		segmentRepo:       segmentRepo,
		flightRepo:        flightRepo,
		flightSeatRepo:    flightSeatRepo,
		seatLockRepo:      seatLockRepo,
		seatAssignRepo:    seatAssignRepo,
		ssrTypeRepo:       ssrTypeRepo,
		passengerSSRRepo:  passengerSSRRepo,
		mealRepo:          mealRepo,
		passengerMealRepo: passengerMealRepo,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (*models.PNR, error) {

	if err := s.validateTripType(req); err != nil {
		return nil, err
	}

	flights, err := s.loadAndValidateFlights(ctx, req.Segments)
	if err != nil {
		return nil, err
	}

	if err := s.validateSegmentChronology(flights); err != nil {
		return nil, err
	}

	if req.TripType == dto.TripRoundTrip {
		if err := s.validateRoundTrip(ctx, flights); err != nil {
			return nil, err
		}
	}

	seatMap, err := s.validateAndLoadSeats(ctx, req)
	if err != nil {
		return nil, err
	}

	locator, err := s.generateLocator(ctx)
	if err != nil {
		return nil, err
	}

	ttl := time.Now().Add(seatLockTTL)

	// ── Step 2: Atomic transaction ─────────────

	var pnr *models.PNR

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		lockExpires := now.Add(seatLockTTL)

		lockedSeats := collectAllSeats(seatMap)
		for _, seat := range lockedSeats {
			result := tx.Model(&models.FlightSeat{}).
				Where("id = ? AND status = ?", seat.ID, models.FlightSeatAvailable).
				Update("status", models.FlightSeatLocked)
			if result.Error != nil {
				return fmt.Errorf("lock seat %s: %w", seat.ID, result.Error)
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("%w: seat %s", utils.ErrSeatAlreadyLocked, SeatNumber(seat))
			}
		}

		// 2b. Create PNR
		p := &models.PNR{
			RecordLocator: locator,
			Status:        models.PNRStatusHold,
			TTL:           &ttl,
		}
		if err := tx.Create(p).Error; err != nil {
			return fmt.Errorf("create PNR: %w", err)
		}

		// 2c. Create SeatLock records (satu per seat)
		for _, seat := range lockedSeats {
			seatID := seat.ID
			pnrID := p.ID
			lock := &models.SeatLock{
				FlightSeatID: &seatID,
				PNRID:        &pnrID,
				LockedAt:     &now,
				ExpiresAt:    &lockExpires,
			}
			if err := tx.Create(lock).Error; err != nil {
				return fmt.Errorf("create seat lock: %w", err)
			}
		}

		// 2d. Create Contact
		pnrID := p.ID
		contact := &models.PNRContact{
			PNRID: &pnrID,
			Name:  req.Contact.Name,
			Email: req.Contact.Email,
			Phone: req.Contact.Phone,
		}
		if err := tx.Create(contact).Error; err != nil {
			return fmt.Errorf("create contact: %w", err)
		}

		// 2e. Create Passengers
		passengers := make([]models.PNRPassenger, 0, len(req.Passengers))
		for _, pr := range req.Passengers {
			pass := models.PNRPassenger{
				PNRID:          &pnrID,
				FirstName:      pr.FirstName,
				LastName:       pr.LastName,
				PassengerType:  pr.PassengerType,
				BirthDate:      pr.BirthDate,
				PassportNumber: pr.PassportNumber,
			}
			if err := tx.Create(&pass).Error; err != nil {
				return fmt.Errorf("create passenger %s %s: %w", pr.FirstName, pr.LastName, err)
			}
			passengers = append(passengers, pass)
		}

		// 2f. Create Segments + 2g. SeatAssignments per segment
		for segIdx, segReq := range req.Segments {
			flightID := flights[segIdx].ID
			segment := &models.PNRSegment{
				PNRID:        &pnrID,
				FlightID:     &flightID,
				SegmentOrder: segIdx + 1,
			}
			if err := tx.Create(segment).Error; err != nil {
				return fmt.Errorf("create segment %d: %w", segIdx+1, err)
			}

			// Create SeatAssignment untuk setiap seat selection di segment ini
			segID := segment.ID
			for _, sel := range segReq.SeatSelections {
				passenger := passengers[sel.PassengerIndex]
				seat := seatMap[segIdx][sel.PassengerIndex]

				passengerID := passenger.ID
				seatID := seat.ID
				assignedAt := now

				assignment := &models.SeatAssignment{
					PassengerID:  &passengerID,
					SegmentID:    &segID,
					FlightSeatID: &seatID,
					AssignedAt:   &assignedAt,
				}
				if err := tx.Create(assignment).Error; err != nil {
					return fmt.Errorf("create seat assignment seg %d pax %d: %w",
						segIdx+1, sel.PassengerIndex, err)
				}
			}
		}

		pnr = p
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return s.pnrRepo.FindWithFull(ctx, pnr.ID)
}

func (s *bookingService) validateTripType(req dto.CreateBookingRequest) error {
	n := len(req.Segments)
	switch req.TripType {
	case dto.TripOneWay:
		if n != 1 {
			return fmt.Errorf("%w: one_way harus memiliki tepat 1 segment, diberikan %d", utils.ErrInvalidSegmentCount, n)
		}
	case dto.TripRoundTrip:
		if n != 2 {
			return fmt.Errorf("%w: round_trip harus memiliki tepat 2 segment, diberikan %d", utils.ErrInvalidSegmentCount, n)
		}
	case dto.TripMultiCity:
		if n < 2 {
			return fmt.Errorf("%w: multi_city minimal 2 segment, diberikan %d", utils.ErrInvalidSegmentCount, n)
		}
	}
	return nil
}

// loadAndValidateFlights memuat semua flight dan memvalidasi statusnya.
func (s *bookingService) loadAndValidateFlights(ctx context.Context, segs []dto.SegmentRequest) ([]*models.Flight, error) {
	flights := make([]*models.Flight, 0, len(segs))
	for i, seg := range segs {
		f, err := s.flightRepo.FindWithDetails(ctx, seg.FlightID)
		if err != nil {
			return nil, fmt.Errorf("segment %d: %w", i+1, utils.ErrFlightNotFound)
		}
		switch f.Status {
		case models.FlightStatusCancelled:
			return nil, fmt.Errorf("segment %d: penerbangan %s dibatalkan", i+1, f.ID)
		case models.FlightStatusDeparted, models.FlightStatusArrived:
			return nil, fmt.Errorf("segment %d: %w", i+1, utils.ErrFlightAlreadyDeparted)
		}
		flights = append(flights, f)
	}
	return flights, nil
}

// validateSegmentChronology memastikan departure tiap segment
// terjadi setelah arrival segment sebelumnya.
// Memberikan toleransi minimum 30 menit koneksi antar flight.
func (s *bookingService) validateSegmentChronology(flights []*models.Flight) error {
	const minConnectionMinutes = 30
	for i := 1; i < len(flights); i++ {
		prev := flights[i-1]
		curr := flights[i]
		if prev.ArrivalTime == nil || curr.DepartureTime == nil {
			continue
		}
		minDep := prev.ArrivalTime.Add(time.Duration(minConnectionMinutes) * time.Minute)
		if curr.DepartureTime.Before(minDep) {
			return fmt.Errorf(
				"%w: segment %d berangkat %s, terlalu dekat dengan arrival segment %d (%s). Minimum koneksi %d menit",
				utils.ErrSegmentChronologyInvalid,
				i+1, curr.DepartureTime.Format("02 Jan 15:04"),
				i, prev.ArrivalTime.Format("02 Jan 15:04"),
				minConnectionMinutes,
			)
		}
	}
	return nil
}

// validateRoundTrip memastikan departure airport segment ke-2
// sama dengan arrival airport segment ke-1.
func (s *bookingService) validateRoundTrip(ctx context.Context, flights []*models.Flight) error {
	outbound := flights[0]
	ret := flights[1]

	// Ambil schedule untuk mendapat airport info
	if outbound.Schedule == nil || ret.Schedule == nil {
		return nil // tidak bisa validasi tanpa schedule, lanjut
	}
	if outbound.Schedule.ArrivalAirportID != ret.Schedule.DepartureAirportID {
		return fmt.Errorf(
			"%w: departure airport return (%s) harus sama dengan arrival airport outbound (%s)",
			utils.ErrRoundTripAirportMismatch,
			ret.Schedule.DepartureAirportID,
			outbound.Schedule.ArrivalAirportID,
		)
	}
	return nil
}

// validateAndLoadSeats memvalidasi semua seat selection untuk semua segment.
// Return: seatMap[segmentIdx][passengerIdx] = *models.FlightSeat
func (s *bookingService) validateAndLoadSeats(
	ctx context.Context,
	req dto.CreateBookingRequest,
) (map[int]map[int]*models.FlightSeat, error) {

	seatMap := make(map[int]map[int]*models.FlightSeat)

	for segIdx, seg := range req.Segments {
		if len(seg.SeatSelections) == 0 {
			seatMap[segIdx] = map[int]*models.FlightSeat{}
			continue
		}

		// Cek duplikat passenger_index dalam segment ini
		seenPax := map[int]bool{}
		// Cek duplikat flight_seat_id dalam segment ini
		seenSeat := map[uuid.UUID]bool{}

		seatMap[segIdx] = make(map[int]*models.FlightSeat)

		for _, sel := range seg.SeatSelections {
			// Validasi passenger_index dalam range
			if sel.PassengerIndex < 0 || sel.PassengerIndex >= len(req.Passengers) {
				return nil, fmt.Errorf(
					"%w: segment %d passenger_index %d tidak valid (total penumpang: %d)",
					utils.ErrInvalidPassengerIndex,
					segIdx+1, sel.PassengerIndex, len(req.Passengers),
				)
			}

			if seenPax[sel.PassengerIndex] {
				return nil, fmt.Errorf(
					"%w: segment %d passenger_index %d muncul lebih dari satu kali",
					utils.ErrDuplicatePassengerSeat,
					segIdx+1, sel.PassengerIndex,
				)
			}
			if seenSeat[sel.FlightSeatID] {
				return nil, fmt.Errorf(
					"%w: segment %d flight_seat_id %s dipilih lebih dari satu penumpang",
					utils.ErrDuplicateSeatSelection,
					segIdx+1, sel.FlightSeatID,
				)
			}

			// Load seat dengan detail
			seat, err := s.flightSeatRepo.FindWithSeatDetail(ctx, sel.FlightSeatID)
			if err != nil {
				return nil, fmt.Errorf("segment %d: %w (seat_id: %s)", segIdx+1, utils.ErrFlightSeatNotFound, sel.FlightSeatID)
			}

			// Validasi seat milik flight ini
			if seat.FlightID == nil || *seat.FlightID != seg.FlightID {
				return nil, fmt.Errorf(
					"%w: segment %d seat %s bukan milik flight %s",
					utils.ErrSeatFlightMismatch,
					segIdx+1, sel.FlightSeatID, seg.FlightID,
				)
			}

			// Validasi status seat
			switch seat.Status {
			case models.FlightSeatBooked:
				return nil, fmt.Errorf("segment %d: %w (seat: %s)", segIdx+1, utils.ErrSeatAlreadyBooked, SeatNumber(seat))
			case models.FlightSeatLocked:
				return nil, fmt.Errorf("segment %d: %w (seat: %s)", segIdx+1, utils.ErrSeatAlreadyLocked, SeatNumber(seat))
			}

			seenPax[sel.PassengerIndex] = true
			seenSeat[sel.FlightSeatID] = true
			seatMap[segIdx][sel.PassengerIndex] = seat
		}
	}
	return seatMap, nil
}

// collectAllSeats mengumpulkan semua FlightSeat unik dari seatMap.
func collectAllSeats(seatMap map[int]map[int]*models.FlightSeat) []*models.FlightSeat {
	seen := map[uuid.UUID]bool{}
	var result []*models.FlightSeat
	for _, passengerSeats := range seatMap {
		for _, seat := range passengerSeats {
			if !seen[seat.ID] {
				seen[seat.ID] = true
				result = append(result, seat)
			}
		}
	}
	return result
}

// SeatNumber helper — returns seat number string or ID fallback.
func SeatNumber(fs *models.FlightSeat) string {
	if fs.AircraftSeat != nil {
		return fs.AircraftSeat.SeatNumber
	}
	return fs.ID.String()[:8]
}

// ─────────────────────────────────────────────
// GetPNR
// ─────────────────────────────────────────────

func (s *bookingService) GetPNR(ctx context.Context, locator string) (*models.PNR, error) {
	pnr, err := s.pnrRepo.FindByLocator(ctx, strings.ToUpper(locator))
	if err != nil {
		return nil, utils.ErrPNRNotFound
	}
	return s.pnrRepo.FindWithFull(ctx, pnr.ID)
}

func (s *bookingService) GetPNRByID(ctx context.Context, id uuid.UUID) (*models.PNR, error) {
	pnr, err := s.pnrRepo.FindWithFull(ctx, id)
	if err != nil {
		return nil, utils.ErrPNRNotFound
	}
	return pnr, nil
}

func (s *bookingService) GetAllPNRs(ctx context.Context, page, limit int) ([]models.PNR, int64, error) {
	return s.pnrRepo.FindAll(ctx, page, limit)
}

// ─────────────────────────────────────────────
// UpdateContact
// ─────────────────────────────────────────────

func (s *bookingService) UpdateContact(ctx context.Context, pnrID uuid.UUID, req dto.UpdateContactRequest) (*models.PNRContact, error) {
	pnr, err := s.pnrRepo.FindByID(ctx, pnrID)
	if err != nil {
		return nil, utils.ErrPNRNotFound
	}
	if pnr.Status == models.PNRStatusCancelled {
		return nil, utils.ErrPNRAlreadyCancelled
	}

	contact, err := s.contactRepo.FindByPNRID(ctx, pnrID)
	if err != nil {
		return nil, fmt.Errorf("contact not found: %w", err)
	}

	if req.Name != "" {
		contact.Name = req.Name
	}
	if req.Email != "" {
		contact.Email = req.Email
	}
	if req.Phone != "" {
		contact.Phone = req.Phone
	}

	if err := s.contactRepo.Update(ctx, contact); err != nil {
		return nil, err
	}
	return contact, nil
}

// ─────────────────────────────────────────────
// SSR
// ─────────────────────────────────────────────

func (s *bookingService) AddSSR(ctx context.Context, passengerID uuid.UUID, req dto.AddSSRRequest) (*models.PassengerSSR, error) {
	if _, err := s.passengerRepo.FindByID(ctx, passengerID); err != nil {
		return nil, utils.ErrPassengerNotFound
	}
	ssrType, err := s.ssrTypeRepo.FindByCode(ctx, req.SSRCode)
	if err != nil {
		return nil, utils.ErrSSRTypeNotFound
	}
	ssr := &models.PassengerSSR{
		PassengerID: &passengerID,
		SegmentID:   &req.SegmentID,
		SSRTypeID:   &ssrType.ID,
	}
	if err := s.passengerSSRRepo.Create(ctx, ssr); err != nil {
		return nil, err
	}
	return ssr, nil
}

func (s *bookingService) RemoveSSR(ctx context.Context, ssrID uuid.UUID) error {
	return s.passengerSSRRepo.Delete(ctx, ssrID)
}

// ─────────────────────────────────────────────
// Meal
// ─────────────────────────────────────────────

func (s *bookingService) AddMeal(ctx context.Context, passengerID uuid.UUID, req dto.AddMealRequest) (*models.PassengerMeal, error) {
	if _, err := s.passengerRepo.FindByID(ctx, passengerID); err != nil {
		return nil, utils.ErrPassengerNotFound
	}
	meal, err := s.mealRepo.FindByCode(ctx, req.MealCode)
	if err != nil {
		return nil, utils.ErrMealNotFound
	}
	pm := &models.PassengerMeal{
		PassengerID: &passengerID,
		SegmentID:   &req.SegmentID,
		MealID:      &meal.ID,
	}
	if err := s.passengerMealRepo.Create(ctx, pm); err != nil {
		return nil, err
	}
	return pm, nil
}

func (s *bookingService) RemoveMeal(ctx context.Context, mealID uuid.UUID) error {
	return s.passengerMealRepo.Delete(ctx, mealID)
}

// ─────────────────────────────────────────────
// CancelPNR — reset semua seat di semua segment
// ─────────────────────────────────────────────

func (s *bookingService) CancelPNR(ctx context.Context, pnrID uuid.UUID) error {
	pnr, err := s.pnrRepo.FindByID(ctx, pnrID)
	if err != nil {
		return utils.ErrPNRNotFound
	}
	if pnr.Status == models.PNRStatusCancelled {
		return utils.ErrPNRAlreadyCancelled
	}
	if pnr.Status == models.PNRStatusTicketed {
		return utils.ErrPNRAlreadyTicketed
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Reset semua flight_seat yang di-assign di PNR ini kembali ke available
		// Join: seat_assignments → pnr_segments → pnrs
		tx.Exec(`
			UPDATE flight_seats SET status = 'available'
			WHERE id IN (
				SELECT sa.flight_seat_id
				FROM seat_assignments sa
				JOIN pnr_segments ps ON ps.id = sa.segment_id
				WHERE ps.pnr_id = ?
			)
		`, pnrID)

		// Hapus semua seat locks milik PNR ini
		tx.Where("pnr_id = ?", pnrID).Delete(&models.SeatLock{})

		// Update PNR status
		return tx.Model(&models.PNR{}).
			Where("id = ?", pnrID).
			Update("status", models.PNRStatusCancelled).Error
	})
}

// ─────────────────────────────────────────────
// generateLocator — crypto-random 6-char alphanumeric
// ─────────────────────────────────────────────

func (s *bookingService) generateLocator(ctx context.Context) (string, error) {
	for attempt := 0; attempt < maxLocatorRetry; attempt++ {
		locator, err := randomString(locatorLen, locatorChars)
		if err != nil {
			return "", fmt.Errorf("generate locator: %w", err)
		}
		if _, err := s.pnrRepo.FindByLocator(ctx, locator); err != nil {
			return locator, nil
		}
	}
	return "", utils.ErrLocatorGenFailed
}

func randomString(n int, charset string) (string, error) {
	b := make([]byte, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[idx.Int64()]
	}
	return string(b), nil
}
