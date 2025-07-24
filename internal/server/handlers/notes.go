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

type NotesHandler struct {
	NotesRepository *repository.NotesRepository
	RoomsRepository *repository.RoomsRepository
}

func (nh *NotesHandler) RegisterNotesRoutes(r chi.Router) {
	r.Route("/notes", func(r chi.Router) {
		r.Post("/", nh.CreateNoteHandler)
		r.Get("/{note_id}", nh.GetNoteByIDHandler)
		r.Put("/{note_id}", nh.UpdateNoteHandler)
		r.Delete("/{note_id}", nh.DeleteNoteHandler)
		r.Get("/room/{room_id}", nh.GetNotesByRoomHandler)
		r.Get("/my-notes", nh.GetUserNotesHandler)
	})
}

func (nh *NotesHandler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	var req dtos.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Content == "" || req.RoomID == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Title, content, and room_id are required")
		return
	}

	if !nh.RoomsRepository.IsUserInRoom(userID, req.RoomID) {
		utils.RespondWithError(w, http.StatusForbidden, "User is not a member of this room")
		return
	}

	note, err := nh.NotesRepository.Create(userID, req.RoomID, req.Title, req.Content)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create note")
		return
	}

	response := dtos.NoteResponse{
		ID:        note.ID,
		UserID:    note.UserID,
		RoomID:    note.RoomID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (nh *NotesHandler) GetNoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	noteIDStr := chi.URLParam(r, "note_id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	note, err := nh.NotesRepository.GetByID(uint(noteID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get note")
		return
	}
	if note == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Note not found")
		return
	}

	if !nh.RoomsRepository.IsUserInRoom(userID, note.RoomID) {
		utils.RespondWithError(w, http.StatusForbidden, "User is not a member of this room")
		return
	}

	response := dtos.NoteResponse{
		ID:        note.ID,
		UserID:    note.UserID,
		RoomID:    note.RoomID,
		Title:     note.Title,
		Content:   note.Content,
		UserName:  note.User.Name,
		UserEmail: note.User.Email,
		RoomName:  note.Room.Name,
		CreatedAt: note.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (nh *NotesHandler) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	noteIDStr := chi.URLParam(r, "note_id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	note, err := nh.NotesRepository.GetByID(uint(noteID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get note")
		return
	}
	if note == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Note not found")
		return
	}

	if note.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Only note creator can update the note")
		return
	}

	var req dtos.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Content == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	if err := nh.NotesRepository.Update(uint(noteID), req.Title, req.Content); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update note")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Note updated successfully"}`))
}

func (nh *NotesHandler) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	noteIDStr := chi.URLParam(r, "note_id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	note, err := nh.NotesRepository.GetByID(uint(noteID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get note")
		return
	}
	if note == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Note not found")
		return
	}

	if note.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Only note creator can delete the note")
		return
	}

	if err := nh.NotesRepository.Delete(uint(noteID)); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete note")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Note deleted successfully"}`))
}

func (nh *NotesHandler) GetNotesByRoomHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	roomIDStr := chi.URLParam(r, "room_id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	if !nh.RoomsRepository.IsUserInRoom(userID, uint(roomID)) {
		utils.RespondWithError(w, http.StatusForbidden, "User is not a member of this room")
		return
	}

	notes, err := nh.NotesRepository.GetByRoomID(uint(roomID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get notes")
		return
	}

	var response []dtos.NoteResponse
	for _, note := range notes {
		response = append(response, dtos.NoteResponse{
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

func (nh *NotesHandler) GetUserNotesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	notes, err := nh.NotesRepository.GetByUserID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get user notes")
		return
	}

	var response []dtos.NoteResponse
	for _, note := range notes {
		response = append(response, dtos.NoteResponse{
			ID:        note.ID,
			UserID:    note.UserID,
			RoomID:    note.RoomID,
			Title:     note.Title,
			Content:   note.Content,
			RoomName:  note.Room.Name,
			CreatedAt: note.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}