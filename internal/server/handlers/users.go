package handlers

import (
	"api-go/internal/repository"
	"api-go/internal/server/dtos"
	"api-go/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepository *repository.UserRepository
}

func (uh *UserHandler) RegisterUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", uh.CreateUserHandler)
		r.Get("/by-email", uh.GetUserByEmailHandler)
	})
}

func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Falha ao decodificar o corpo da requisição")
		return
	}

	if req.Email == "" || req.Name == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Campos 'email', 'name', e 'password' são obrigatórios")
		return
	}

	// 1. Chame o método Create, que agora retorna o usuário criado
	user, err := uh.UserRepository.Create(req.Email, req.Name, req.Password)
	if err != nil {
		// Verifica se o erro é de usuário duplicado
		if err.Error() == fmt.Sprintf("usuário com email %s já existe", req.Email) {
			utils.RespondWithError(w, http.StatusConflict, "Usuário já existe")
			return
		}
		// Para outros erros, retorna um erro genérico de servidor
		utils.RespondWithError(w, http.StatusInternalServerError, "Falha ao criar usuário")
		return
	}

	response := dtos.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Falha ao codificar resposta JSON: %v", err)
	}
}

func (uh *UserHandler) GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Parâmetro 'email' é obrigatório")
		return
	}

	user, err := uh.UserRepository.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "Usuário não encontrado")
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Falha ao buscar usuário: %v", err))
		return
	}

	response := dtos.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Falha ao codificar resposta: %v", err), http.StatusInternalServerError)
	}
}
