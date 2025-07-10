package handlers

import (
	"api-go/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ReservationsHandler struct {
	ReservationsRepository *repository.ReservationsRepository
}

func (rh *ReservationsHandler) RegisterReservationsRoutes(r chi.Router) {
	r.Route("/reservations", func(r chi.Router) {
		r.Post("/", rh.CreateReservationHandler)
		r.Get("/{reservation_id}", rh.GetReservationByIDHandler)
		r.Put("/{reservation_id}", rh.UpdateReservationHandler)
		r.Delete("/{reservation_id}", rh.DeleteReservationHandler)
		r.Get("/by-user/{user_id}", rh.GetReservationsByUserIDHandler)
		r.Get("/by-room/{room_id}", rh.GetReservationsByRoomIDHandler)
	})
}

func (rh *ReservationsHandler) CreateReservationHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func (rh *ReservationsHandler) GetReservationByIDHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func (rh *ReservationsHandler) UpdateReservationHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func (rh *ReservationsHandler) DeleteReservationHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func (rh *ReservationsHandler) GetReservationsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func (rh *ReservationsHandler) GetReservationsByRoomIDHandler(w http.ResponseWriter, r *http.Request) {
	//
}
