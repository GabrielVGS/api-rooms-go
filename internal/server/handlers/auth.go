package handlers

import (
	"api-go/internal/auth"
	"api-go/internal/models"
	"api-go/internal/repository"
	"api-go/internal/server/dtos"
	"api-go/internal/server/middlewares"
	"api-go/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	UserRepository *repository.UserRepository
}

func (ah *AuthHandler) RegisterAuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", ah.LoginHandler)
		r.Post("/register", ah.RegisterHandler)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware) // Middleware de autenticação
			r.Get("/profile", ah.GetProfileHandler)
		})
	})
}

// LoginHandler authenticates user and returns JWT token
//
//	@Summary		User login
//	@Description	Authenticate user with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.AuthLoginRequest	true	"Login credentials"
//	@Success		200		{object}	dtos.AuthLoginResponse
//	@Failure		400		{object}	dtos.ErrorResponse
//	@Failure		401		{object}	dtos.ErrorResponse
//	@Failure		500		{object}	dtos.ErrorResponse
//	@Router			/auth/login [post]
func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Handle user login
	var loginRequest dtos.AuthLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate loginRequest and authenticate user
	// If successful, generate JWT token and respond with it

	// First, we need to get the user from the database
	user, _ := ah.UserRepository.GetByEmail(loginRequest.Email)

	fmt.Println("User found:", user)
	fmt.Println("Password provided:", loginRequest.Password)

	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	response := dtos.AuthLoginResponse{
		Token: token,
		User: dtos.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not encode response")
		return
	}

}

// RegisterHandler creates new user account
//
//	@Summary		User registration
//	@Description	Register new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.AuthRegisterRequest	true	"Registration details"
//	@Success		201		{object}	dtos.AuthLoginResponse
//	@Failure		400		{object}	dtos.ErrorResponse
//	@Failure		409		{object}	dtos.ErrorResponse
//	@Failure		500		{object}	dtos.ErrorResponse
//	@Router			/auth/register [post]
func (ah *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Handle user registration
	var registerRequest dtos.AuthRegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	//print registerRequest

	// Check if the user already exists
	existingUser, _ := ah.UserRepository.GetByEmail(registerRequest.Email)

	if existingUser != nil {
		utils.RespondWithError(w, http.StatusConflict, "User already exists")
		return
	}
	// Validate registerRequest and create user

	hashedPassword, err := utils.HashPassword(registerRequest.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not hash password")
		return
	}

	user := &models.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: hashedPassword,
	}

	createdUser, err := ah.UserRepository.Create(
		user.Email,
		user.Name,
		user.Password,
	)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	token, err := auth.GenerateToken(createdUser)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	response := dtos.AuthLoginResponse{
		Token: token,
		User: dtos.UserResponse{
			ID:    createdUser.ID,
			Name:  createdUser.Name,
			Email: createdUser.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not encode response")
		return
	}
}

// GetProfileHandler gets current user profile
//
//	@Summary		Get user profile
//	@Description	Get current authenticated user profile
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.AuthProfileResponse
//	@Failure		401	{object}	dtos.ErrorResponse
//	@Failure		500	{object}	dtos.ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/profile [get]
func (ah *AuthHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Get full user data from repository
	user, err := ah.UserRepository.GetByID(claims.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not fetch user profile")
		return
	}

	if user == nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Create response
	response := dtos.AuthProfileResponse{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not encode response")
		return
	}
}
