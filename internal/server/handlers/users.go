package handlers

import (
	"api-go/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm" // Importe o GORM para checar o erro específico
)

// MUDANÇA 1: Usar um ponteiro para o repositório.
type UserHandler struct {
	UserRepository *repository.UserRepository
}

// DTO (Data Transfer Object) para a criação de usuário.
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// DTO para a resposta, omitindo dados sensíveis.
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (uh *UserHandler) RegisterUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", uh.CreateUserHandler)
		r.Get("/by-email", uh.GetUserByEmailHandler)
	})
}

func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	// Decodifica o corpo da requisição
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	// Valida os dados de entrada
	if req.Email == "" || req.Name == "" || req.Password == "" {
		http.Error(w, "Campos 'email', 'name', e 'password' são obrigatórios", http.StatusBadRequest)
		return
	}

	// MUDANÇA 2: O handler não faz mais o hash. Ele passa a senha em texto plano
	// para o repositório, que é o responsável por essa lógica.
	// MUDANÇA 3: Chamando o método corrigido 'Create'.
	err := uh.UserRepository.Create(req.Email, req.Name, req.Password)
	if err != nil {
		// Se o erro for de usuário já existente, retorne um status mais apropriado.
		if err.Error() == fmt.Sprintf("usuário com email %s já existe", req.Email) {
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
			return
		}
		http.Error(w, fmt.Sprintf("Falha ao criar usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Responde com uma mensagem de sucesso.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuário criado com sucesso"})
}

func (uh *UserHandler) GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Parâmetro 'email' é obrigatório", http.StatusBadRequest)
		return
	}

	// MUDANÇA 4: Chamando o método corrigido 'GetByEmail'.
	user, err := uh.UserRepository.GetByEmail(email)
	if err != nil {
		// MUDANÇA 5: Tratamento de erro específico para "não encontrado".
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound) // 404 Not Found
			return
		}
		// Para todos os outros erros, usamos 500.
		http.Error(w, fmt.Sprintf("Falha ao buscar usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// MUDANÇA 6 (CORREÇÃO DE SEGURANÇA): Criamos uma resposta segura sem a senha.
	response := UserResponse{
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
