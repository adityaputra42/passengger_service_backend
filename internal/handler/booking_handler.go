package handler

import (
	"encoding/json"
	"net/http"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/services"
	"passenger_service_backend/internal/utils"

	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	svc services.BookingService
}

func NewBookingHandler(svc services.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

// Create godoc
// @Summary      Buat pemesanan (PNR)
// @Description  Membuat pemesanan tiket baru. Mendukung one_way, round_trip, dan multi_city. Seat akan di-lock selama 30 menit.
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateBookingRequest  true  "Data pemesanan"
// @Success      201      {object}  utils.Response{data=dto.PNRResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      404      {object}  utils.Response  "Penerbangan atau kursi tidak ditemukan"
// @Failure      409      {object}  utils.Response  "Kursi sudah dipesan"
// @Security     BearerAuth
// @Router       /bookings [post]
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	pnr, err := h.svc.CreateBooking(r.Context(), req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, "success", dto.ToPNRResponse(pnr))

}

// List godoc
// @Summary      Daftar semua PNR
// @Description  Mengambil seluruh daftar PNR. Hanya admin dan super_admin.
// @Tags         Booking
// @Produce      json
// @Param        page   query     int  false  "Halaman (default: 1)"
// @Param        limit  query     int  false  "Jumlah per halaman (default: 10)"
// @Success      200    {object}  utils.Response{data=dto.PNRListResponse}
// @Failure      401    {object}  utils.Response
// @Failure      403    {object}  utils.Response
// @Security     BearerAuth
// @Router       /bookings [get]
func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := utils.PageLimit(r)
	pnrs, total, err := h.svc.GetAllPNRs(r.Context(), page, limit)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPNRListResponse(pnrs, total, page, limit))

}

// GetByID godoc
// @Summary      Detail pemesanan berdasarkan ID
// @Description  Mengambil detail PNR lengkap termasuk penumpang, segmen, kursi, dan pembayaran.
// @Tags         Booking
// @Produce      json
// @Param        id   path      string  true  "PNR UUID"
// @Success      200  {object}  utils.Response{data=dto.PNRResponse}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     BearerAuth
// @Router       /bookings/{id} [get]
func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	pnr, err := h.svc.GetPNRByID(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPNRResponse(pnr))
}

// GetByLocator godoc
// @Summary      Cari pemesanan berdasarkan kode booking
// @Description  Mengambil detail PNR menggunakan 6-karakter record locator (contoh: ABC123).
// @Tags         Booking
// @Produce      json
// @Param        locator  path      string  true  "Record locator (6 karakter)"  example(ABC123)
// @Success      200      {object}  utils.Response{data=dto.PNRResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /bookings/locator/{locator} [get]
func (h *BookingHandler) GetByLocator(w http.ResponseWriter, r *http.Request) {
	locator := chi.URLParam(r, "locator")
	pnr, err := h.svc.GetPNR(r.Context(), locator)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPNRResponse(pnr))
}

// UpdateContact godoc
// @Summary      Update kontak pemesanan
// @Description  Memperbarui informasi kontak (nama, email, telepon) pada sebuah PNR.
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        id       path      string                    true  "PNR UUID"
// @Param        request  body      dto.UpdateContactRequest  true  "Data kontak baru"
// @Success      200      {object}  utils.Response{data=dto.PNRContactResponse}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Security     BearerAuth
// @Router       /bookings/{id}/contact [put]
func (h *BookingHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	var req dto.UpdateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	contact, err := h.svc.UpdateContact(r.Context(), id, req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPNRContactResponse(contact))

}

// Cancel godoc
// @Summary      Batalkan pemesanan
// @Description  Membatalkan PNR dan melepaskan kursi yang di-lock. PNR yang sudah ticketed tidak bisa dibatalkan.
// @Tags         Booking
// @Produce      json
// @Param        id   path      string  true  "PNR UUID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      422  {object}  utils.Response  "PNR sudah dibatalkan atau sudah ticketed"
// @Security     BearerAuth
// @Router       /bookings/{id} [delete]
func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "id")
	if !ok {
		return
	}
	err := h.svc.CancelPNR(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", nil)

}

// AddSSR godoc
// @Summary      Tambah SSR ke penumpang
// @Description  Menambahkan Special Service Request (kursi roda, meal khusus, dll) ke penumpang.
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        passengerID  path      string          true  "Passenger UUID"
// @Param        request      body      dto.AddSSRRequest  true  "Data SSR"
// @Success      200          {object}  utils.Response{data=dto.PassengerSSRResponse}
// @Failure      400          {object}  utils.Response
// @Failure      401          {object}  utils.Response
// @Failure      404          {object}  utils.Response
// @Security     BearerAuth
// @Router       /bookings/passengers/{passengerID}/ssr [post]
func (h *BookingHandler) AddSSR(w http.ResponseWriter, r *http.Request) {
	passengerID, ok := utils.UUIDParam(w, r, "passengerID")
	if !ok {
		return
	}
	var req dto.AddSSRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	ssr, err := h.svc.AddSSR(r.Context(), passengerID, req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPassengerSSRResponse(ssr))

}

// RemoveSSR godoc
// @Summary      Hapus SSR penumpang
// @Description  Menghapus SSR yang sudah ditambahkan ke penumpang.
// @Tags         Booking
// @Produce      json
// @Param        ssrID  path      string  true  "SSR UUID"
// @Success      200    {object}  utils.Response
// @Failure      400    {object}  utils.Response
// @Failure      401    {object}  utils.Response
// @Security     BearerAuth
// @Router       /bookings/ssr/{ssrID} [delete]
func (h *BookingHandler) RemoveSSR(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "ssrID")
	if !ok {
		return
	}
	err := h.svc.RemoveSSR(r.Context(), id)

	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", nil)
}

// AddMeal godoc
// @Summary      Tambah pilihan meal
// @Description  Menambahkan pilihan meal ke penumpang untuk segment tertentu.
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        passengerID  path      string            true  "Passenger UUID"
// @Param        request      body      dto.AddMealRequest  true  "Data meal"
// @Success      200          {object}  utils.Response{data=dto.PassengerMealResponse}
// @Failure      400          {object}  utils.Response
// @Failure      401          {object}  utils.Response
// @Failure      404          {object}  utils.Response
// @Security     BearerAuth
// @Router       /bookings/passengers/{passengerID}/meal [post]
func (h *BookingHandler) AddMeal(w http.ResponseWriter, r *http.Request) {
	passengerID, ok := utils.UUIDParam(w, r, "passengerID")
	if !ok {
		return
	}
	var req dto.AddMealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	meal, err := h.svc.AddMeal(r.Context(), passengerID, req)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", dto.ToPassengerMealResponse(meal))

}

// RemoveMeal godoc
// @Summary      Hapus pilihan meal
// @Description  Menghapus pilihan meal penumpang.
// @Tags         Booking
// @Produce      json
// @Param        mealID  path      string  true  "Meal UUID"
// @Success      200     {object}  utils.Response
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Security     BearerAuth
// @Router       /bookings/meal/{mealID} [delete]
func (h *BookingHandler) RemoveMeal(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.UUIDParam(w, r, "mealID")
	if !ok {
		return
	}
	err := h.svc.RemoveMeal(r.Context(), id)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "success", nil)
}
