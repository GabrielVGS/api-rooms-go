package repository

import (
	"api-go/internal/models"

	"gorm.io/gorm"
)

type RoomsRepository struct {
	DB *gorm.DB
}

func NewRoomsRepository(db *gorm.DB) *RoomsRepository {
	return &RoomsRepository{
		DB: db,
	}
}

func (r *RoomsRepository) Create(name string, description string, capacity int) error {
	room := models.Room{
		Name:        name,
		Description: description,
		Capacity:    capacity,
	}

	if err := r.DB.Create(&room).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoomsRepository) GetByID(id uint) (*models.Room, error) {
	var room models.Room
	if err := r.DB.First(&room, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *RoomsRepository) GetAll() ([]models.Room, error) {
	var rooms []models.Room
	if err := r.DB.Find(&rooms).Error; err != nil {
		return nil, err // Retorna o erro
	}
	return rooms, nil // Retorna a lista de salas
}

func (r *RoomsRepository) Update(id uint, name string, description string, capacity int) error {
	room := models.Room{
		Model:       gorm.Model{ID: id},
		Name:        name,
		Description: description,
		Capacity:    capacity,
	}

	if err := r.DB.Save(&room).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoomsRepository) Delete(id uint) error {
	if err := r.DB.Delete(&models.Room{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // Retorna nil se não encontrar o registro para excluir
		}
		return err
	}
	return nil
}

func (r *RoomsRepository) GetByName(name string) (*models.Room, error) {
	var room models.Room
	if err := r.DB.Where("name = ?", name).First(&room).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err // Retorna o erro se houver outro problema
	}
	return &room, nil
}
