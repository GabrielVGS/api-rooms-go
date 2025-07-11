package repository

import (
	"api-go/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(email string, name string, password string) (*models.User, error) {
	var userCount int64
	// Verifica se um usuário com o email já existe
	if err := r.DB.Model(&models.User{}).Where("email = ?", email).
		Count(&userCount).Error; err != nil {
		return nil, fmt.Errorf("falha ao verificar se o usuário existe: %w", err)
	}
	if userCount > 0 {
		return nil, fmt.Errorf("usuário com email %s já existe", email)
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("falha ao criar usuário: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("usuário com ID %d não encontrado", id)
		}
		return nil, fmt.Errorf("falha ao buscar usuário por ID: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("falha ao buscar usuários: %w", err)
	}
	return users, nil
}

func (r *UserRepository) Update(user *models.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return fmt.Errorf("falha ao atualizar usuário: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(id uint) error {
	if err := r.DB.Delete(&models.User{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("usuário com ID %d não encontrado para exclusão", id)
		}
		return fmt.Errorf("falha ao excluir usuário: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("usuário com email %s não encontrado", email)
		}
		return nil, fmt.Errorf("falha ao buscar usuário por email: %w", err)
	}
	return &user, nil
}
