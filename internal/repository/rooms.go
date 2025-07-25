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

func (r *RoomsRepository) Create(name string, description string, subject string, capacity int, createdBy uint) (*models.Room, error) {
	room := models.Room{
		Name:        name,
		Description: description,
		Subject:     subject,
		Capacity:    capacity,
		CreatedBy:   createdBy,
	}

	if err := r.DB.Create(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomsRepository) GetByID(id uint) (*models.Room, error) {
	var room models.Room
	if err := r.DB.Preload("Members.User").Preload("Notes.User").First(&room, id).Error; err != nil {
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

func (r *RoomsRepository) Update(id uint, name string, description string, subject string, capacity int) error {
	updates := map[string]interface{}{
		"name":        name,
		"description": description,
		"subject":     subject,
		"capacity":    capacity,
	}

	if err := r.DB.Model(&models.Room{}).Where("id = ?", id).Updates(updates).Error; err != nil {
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
		return nil, err
	}
	return &room, nil
}

func (r *RoomsRepository) JoinRoom(userID, roomID uint) error {
	member := models.RoomMember{
		UserID: userID,
		RoomID: roomID,
		Role:   "member",
	}
	return r.DB.Create(&member).Error
}

func (r *RoomsRepository) LeaveRoom(userID, roomID uint) error {
	return r.DB.Where("user_id = ? AND room_id = ?", userID, roomID).Delete(&models.RoomMember{}).Error
}

func (r *RoomsRepository) IsUserInRoom(userID, roomID uint) bool {
	var count int64
	r.DB.Model(&models.RoomMember{}).Where("user_id = ? AND room_id = ?", userID, roomID).Count(&count)
	return count > 0
}

func (r *RoomsRepository) GetUserRooms(userID uint) ([]models.Room, error) {
	var rooms []models.Room
	err := r.DB.Joins("JOIN room_members ON rooms.id = room_members.room_id").
		Where("room_members.user_id = ?", userID).
		Find(&rooms).Error
	return rooms, err
}

func (r *RoomsRepository) GetRoomMemberCount(roomID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.RoomMember{}).Where("room_id = ?", roomID).Count(&count).Error
	return count, err
}
