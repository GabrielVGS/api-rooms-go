package handlers

import (
	"api-go/internal/repository"
	"api-go/internal/server/dtos"
	"api-go/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RoomsHandler struct {
	RoomsRepository *repository.RoomsRepository
}

func (rh *RoomsHandler) RegisterRoomsRoutes(r chi.Router) {
	r.Route("/rooms", func(r chi.Router) {
		r.Post("/", rh.CreateRoomsHandler)
		r.Get("/", rh.GetAllRoomsHandler)
		r.Get("/{room_id}", rh.GetRoomByIDHandler)
		r.Put("/{room_id}", rh.UpdateRoomsHandler)
		r.Delete("/{room_id}", rh.DeleteRoomsHandler)
		r.Post("/{room_id}/join", rh.JoinRoomHandler)
		r.Delete("/{room_id}/leave", rh.LeaveRoomHandler)
		r.Get("/my-rooms", rh.GetUserRoomsHandler)
	})
}

func (rh *RoomsHandler) CreateRoomsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	var req dtos.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Subject == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Name and subject are required")
		return
	}

	room, err := rh.RoomsRepository.Create(req.Name, req.Description, req.Subject, req.Capacity, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create room")
		return
	}

	if err := rh.RoomsRepository.JoinRoom(userID, room.ID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to join created room")
		return
	}

	response := dtos.RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Subject:     room.Subject,
		Capacity:    room.Capacity,
		CreatedBy:   room.CreatedBy,
		CreatedAt:   room.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   room.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (rh *RoomsHandler) GetAllRoomsHandler(w http.ResponseWriter, r *http.Request) {
	rooms, err := rh.RoomsRepository.GetAll()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get rooms")
		return
	}

	var response []dtos.RoomResponse
	for _, room := range rooms {
		response = append(response, dtos.RoomResponse{
			ID:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			Subject:     room.Subject,
			Capacity:    room.Capacity,
			CreatedBy:   room.CreatedBy,
			CreatedAt:   room.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   room.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (rh *RoomsHandler) GetRoomByIDHandler(w http.ResponseWriter, r *http.Request) {
	roomIDStr := chi.URLParam(r, "room_id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	room, err := rh.RoomsRepository.GetByID(uint(roomID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get room")
		return
	}
	if room == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Room not found")
		return
	}

	response := dtos.RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Subject:     room.Subject,
		Capacity:    room.Capacity,
		CreatedBy:   room.CreatedBy,
		CreatedAt:   room.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   room.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	for _, member := range room.Members {
		response.Members = append(response.Members, dtos.RoomMemberResponse{
			UserID:    member.UserID,
			UserName:  member.User.Name,
			UserEmail: member.User.Email,
			Role:      member.Role,
			JoinedAt:  member.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	for _, note := range room.Notes {
		response.Notes = append(response.Notes, dtos.NoteResponse{
			ID:        note.ID,
			UserID:    note.UserID,
			RoomID:    note.RoomID,
			Title:     note.Title,
			Content:   note.Content,
			UserName:  note.User.Name,
			UserEmail: note.User.Email,
			CreatedAt: note.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (rh *RoomsHandler) UpdateRoomsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	roomIDStr := chi.URLParam(r, "room_id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	room, err := rh.RoomsRepository.GetByID(uint(roomID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get room")
		return
	}
	if room == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Room not found")
		return
	}

	if room.CreatedBy != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Only room creator can update the room")
		return
	}

	var req dtos.UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Subject == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Name and subject are required")
		return
	}

	if err := rh.RoomsRepository.Update(uint(roomID), req.Name, req.Description, req.Subject, req.Capacity); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update room")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Room updated successfully"}`))
}

func (rh *RoomsHandler) DeleteRoomsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	roomIDStr := chi.URLParam(r, "room_id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	room, err := rh.RoomsRepository.GetByID(uint(roomID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get room")
		return
	}
	if room == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Room not found")
		return
	}

	if room.CreatedBy != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Only room creator can delete the room")
		return
	}

	if err := rh.RoomsRepository.Delete(uint(roomID)); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete room")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Room deleted successfully"}`))
}

func (rh *RoomsHandler) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	roomIDStr := chi.URLParam(r, "room_id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	room, err := rh.RoomsRepository.GetByID(uint(roomID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get room")
		return
	}
	if room == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Room not found")
		return
	}

	if rh.RoomsRepository.IsUserInRoom(userID, uint(roomID)) {
		utils.RespondWithError(w, http.StatusConflict, "User already in room")
		return
	}

	if err := rh.RoomsRepository.JoinRoom(userID, uint(roomID)); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to join room")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Joined room successfully"}`))
}

func (rh *RoomsHandler) LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	roomIDStr := chi.URLParam(r, "room_id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	if !rh.RoomsRepository.IsUserInRoom(userID, uint(roomID)) {
		utils.RespondWithError(w, http.StatusNotFound, "User not in room")
		return
	}

	if err := rh.RoomsRepository.LeaveRoom(userID, uint(roomID)); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to leave room")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Left room successfully"}`))
}

func (rh *RoomsHandler) GetUserRoomsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	rooms, err := rh.RoomsRepository.GetUserRooms(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get user rooms")
		return
	}

	var response []dtos.RoomResponse
	for _, room := range rooms {
		response = append(response, dtos.RoomResponse{
			ID:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			Subject:     room.Subject,
			Capacity:    room.Capacity,
			CreatedBy:   room.CreatedBy,
			CreatedAt:   room.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   room.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
