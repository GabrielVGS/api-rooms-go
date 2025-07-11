package server

import (
	"net/http"

	// Imports corretos para seus pacotes internos
	"api-go/internal/repository"
	"api-go/internal/server/handlers"

	"api-go/internal/server/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger) // Middleware para logar as requisições

	// Configuração do CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rota de healthcheck

	r.Get("/health", s.healthHandler)

	// Criação do repositórios

	userRepo := repository.NewUserRepository(s.db.GetDB())
	// TODO: Criar os repositórios de Rooms e Reservations

	// Criação dos Handlers
	userHandler := handlers.UserHandler{
		UserRepository: userRepo,
	}

	authHandler := handlers.AuthHandler{
		UserRepository: userRepo,
	}
	// TODO: Criar os Handlers de Rooms e Reservations

	// Registro das rotas
	// Rotas terão o prefixo /api | EX /api/users, /api/rooms, /api/reservations
	r.Route("/api", func(r chi.Router) {
		authHandler.RegisterAuthRoutes(r)
		r.Group(func(r chi.Router) {
			// Aqui irao ficar as rotas que precisam de autenticação
			r.Use(middlewares.AuthMiddleware)
			userHandler.RegisterUserRoutes(r)

			// TODO: Registrar as rotas de Rooms e Reservations
		})
	})

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "up"}`))
}
