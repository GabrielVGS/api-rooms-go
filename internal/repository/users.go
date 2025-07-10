package repository

import (
	"api-go/internal/models"
	"api-go/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// A função New agora espera *gorm.DB, tornando a inicialização mais direta.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// MUDANÇA 2: O nome do método é mais simples e direto.
// Já que estamos no UserRepository, "Create" é suficiente.
func (r *UserRepository) Create(email string, name string, password string) error {
	// MUDANÇA 3: Não precisamos mais chamar GetDB(). Usamos r.DB diretamente.
	// E removemos a verificação de nil, pois agora é responsabilidade de quem
	// cria o UserRepository (o NewUserRepository) garantir que o DB não seja nulo.

	var userCount int64
	// Verifica se um usuário com o email já existe
	if err := r.DB.Model(&models.User{}).Where("email = ?", email).
		Count(&userCount).Error; err != nil {
		return fmt.Errorf("falha ao verificar se o usuário existe: %w", err)
	}
	if userCount > 0 {
		return fmt.Errorf("usuário com email %s já existe", email)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("falha ao gerar hash da senha: %w", err)
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	if err := r.DB.Create(&user).Error; err != nil {
		return fmt.Errorf("falha ao criar usuário: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	// Usamos r.DB diretamente aqui também.
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		// É uma boa prática verificar se o erro é "record not found"
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("usuário com email %s não encontrado", email)
		}
		return nil, fmt.Errorf("falha ao buscar usuário por email: %w", err)
	}
	return &user, nil
}
