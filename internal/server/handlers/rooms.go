package handlers

import (
	"api-go/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RoomsHandler struct {
	RoomsRepository *repository.RoomsRepository
}

func (rh *RoomsHandler) RegisterRoomsRoutes(r chi.Router) {
	r.Route("/rooms", func(r chi.Router) {
		r.Post("/", rh.CreateRoomsHandler)
		r.Get("/", rh.GetRoomsByIDHandler)
		r.Put("/{rooms_id}", rh.UpdateRoomsHandler)
		r.Delete("/{rooms_id}", rh.DeleteRoomsHandler)
	})
}

func (rh *RoomsHandler) CreateRoomsHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func (rh *RoomsHandler) GetRoomsByIDHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func (rh *RoomsHandler) UpdateRoomsHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func (rh *RoomsHandler) DeleteRoomsHandler(w http.ResponseWriter, r *http.Request) {
	//
}
