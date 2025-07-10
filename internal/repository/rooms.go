package repository

import (
	"api-go/internal/models"

	"gorm.io/gorm"
)

type RoomsRepository struct {
	DB *gorm.DB
}

// A função New agora espera *gorm.DB, tornando a inicialização mais direta.
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
		// É uma boa prática verificar se o erro é "record not found"
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Retorna nil se não encontrar o registro
		}
		return nil, err // Retorna o erro se houver outro problema
	}
	return &room, nil
}

func (r *RoomsRepository) GetAll() ([]models.Room, error) {
	var rooms []models.Room
	if err := r.DB.Find(&rooms).Error; err != nil {
		return nil, err // Retorna o erro se houver um problema ao buscar os registros
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
		return err // Retorna o erro se houver um problema ao atualizar o registro
	}
	return nil // Retorna nil se a atualização for bem-sucedida
}

func (r *RoomsRepository) Delete(id uint) error {
	if err := r.DB.Delete(&models.Room{}, id).Error; err != nil {
		// É uma boa prática verificar se o erro é "record not found"
		if err == gorm.ErrRecordNotFound {
			return nil // Retorna nil se não encontrar o registro para excluir
		}
		return err // Retorna o erro se houver outro problema
	}
	return nil // Retorna nil se a exclusão for bem-sucedida
}

func (r *RoomsRepository) GetByName(name string) (*models.Room, error) {
	var room models.Room
	if err := r.DB.Where("name = ?", name).First(&room).Error; err != nil {
		// É uma boa prática verificar se o erro é "record not found"
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Retorna nil se não encontrar o registro
		}
		return nil, err // Retorna o erro se houver outro problema
	}
	return &room, nil // Retorna a sala encontrada
}
