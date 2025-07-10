package server

import (
	"net/http"

	// Imports corretos para seus pacotes internos
	"api-go/internal/repository"
	"api-go/internal/server/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Configuração do CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// --- Rotas Públicas ---

	r.Get("/health", s.healthHandler)

	// --- Conectando o UserHandler ---

	// 1. Inicialize o Repositório
	// Pegamos a conexão GORM (*gorm.DB) do nosso serviço de banco de dados.
	userRepo := repository.NewUserRepository(s.db.GetDB())

	// 2. Inicialize o Handler
	// Injetamos o repositório que acabamos de criar no handler.
	// Note que usamos um ponteiro para o repositório, como corrigimos antes.
	userHandler := handlers.UserHandler{
		UserRepository: userRepo,
	}

	// 3. Registre as Rotas do Usuário
	// É uma boa prática agrupar as rotas da API sob um prefixo como "/api".
	r.Route("/api", func(r chi.Router) {
		userHandler.RegisterUserRoutes(r) // Isso criará as rotas /api/users/...
	})

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "up"}`))
}
