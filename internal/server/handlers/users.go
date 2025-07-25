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
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepository *repository.UserRepository
}

func (uh *UserHandler) RegisterUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", uh.CreateUserHandler)
		r.Get("/", uh.GetUserByIDHandler)
		r.Get("/", uh.GetAllUsersHandler)
		r.Put("/{user_id}", uh.UpdateUserHandler)
		r.Delete("/{user_id}", uh.DeleteUserHandler)
		r.Get("/by-email", uh.GetUserByEmailHandler)
	})
}

// CreateUserHandler creates a new user
//
//	@Summary		Create user
//	@Description	Create a new user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.CreateUserRequest	true	"User creation details"
//	@Success		201		{object}	dtos.UserResponse
//	@Failure		400		{object}	dtos.ErrorResponse
//	@Failure		409		{object}	dtos.ErrorResponse
//	@Failure		500		{object}	dtos.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users [post]
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

	user, err := uh.UserRepository.Create(req.Email, req.Name, req.Password)
	if err != nil {
		if err.Error() == fmt.Sprintf("usuário com email %s já existe", req.Email) {
			utils.RespondWithError(w, http.StatusConflict, "Usuário já existe")
			return
		}
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

// GetUserByEmailHandler gets user by email
//
//	@Summary		Get user by email
//	@Description	Retrieve user information by email address
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			email	query		string	true	"User email address"
//	@Success		200		{object}	dtos.UserResponse
//	@Failure		400		{object}	dtos.ErrorResponse
//	@Failure		404		{object}	dtos.ErrorResponse
//	@Failure		500		{object}	dtos.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/by-email [get]
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

// GetUserByIDHandler gets user by ID
//
//	@Summary		Get user by ID
//	@Description	Retrieve user information by user ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	query		int	true	"User ID"
//	@Success		200		{object}	dtos.UserResponse
//	@Failure		400		{object}	dtos.ErrorResponse
//	@Failure		404		{object}	dtos.ErrorResponse
//	@Failure		500		{object}	dtos.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users [get]
func (uh *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("user_id")

	if ID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Parâmetro 'user_id' é obrigatório")
		return
	}

	userID, err := strconv.Atoi(ID)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Parâmetro 'user_id' deve ser um número")
		return
	}

	user, err := uh.UserRepository.GetByID(uint(userID))

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

func (uh *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	users, err := uh.UserRepository.GetAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "Nenhum usuario")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	userList := make([]dtos.UserResponse, len(users))
	for i, u := range users {
		userList[i] = dtos.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}
	response := dtos.UserListResponse{
		Users: userList,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Falha ao codificar resposta: %v", err), http.StatusInternalServerError)
	}

}

// UpdateUserHandler updates user information
//
//	@Summary		Update user
//	@Description	Update user information by user ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		int						true	"User ID"
//	@Param			request	body		dtos.UpdateUserRequest	true	"User update details"
//	@Success		200		{object}	dtos.UserResponse
//	@Failure		400		{object}	dtos.ErrorResponse
//	@Failure		404		{object}	dtos.ErrorResponse
//	@Failure		409		{object}	dtos.ErrorResponse
//	@Failure		500		{object}	dtos.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/{user_id} [put]
func (uh *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dtos.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Falha ao decodificar o corpo da requisição")
		return
	}

	if req.Email == "" && req.Name == "" && req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Pelo menos um campo ('email', 'name', 'password') deve ser fornecido")
		return
	}

	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Parâmetro 'user_id' é obrigatório")
		return
	}

	id, err := strconv.Atoi(userID)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Parâmetro 'user_id' deve ser um número")
		return
	}

	user, err := uh.UserRepository.GetByID(uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "Usuário não encontrado")
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Falha ao buscar usuário: %v", err))
		return
	}

	if req.Email != "" {
		if !utils.IsValidEmail(req.Email) {
			utils.RespondWithError(w, http.StatusBadRequest, "Email inválido")
			return
		}

		existingUser, err := uh.UserRepository.GetByEmail(req.Email)
		if err == nil && existingUser.ID != user.ID {
			utils.RespondWithError(w, http.StatusConflict, "Email já está em uso por outro usuário")
			return
		}

		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Falha ao hashear a senha")
			return
		}
		user.Password = hashedPassword
	}

	if err := uh.UserRepository.Update(user); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Falha ao atualizar usuário: %v", err))
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

// DeleteUserHandler deletes a user
//
//	@Summary		Delete user
//	@Description	Delete user by user ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path	int	true	"User ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	dtos.ErrorResponse
//	@Failure		404		{object}	dtos.ErrorResponse
//	@Failure		500		{object}	dtos.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/{user_id} [delete]
func (uh *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Parâmetro 'user_id' é obrigatório")
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Parâmetro 'user_id' deve ser um número")
		return
	}

	if err := uh.UserRepository.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "Usuário não encontrado")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Falha ao deletar usuário: %v", err))
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
