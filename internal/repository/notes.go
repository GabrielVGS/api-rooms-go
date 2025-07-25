package repository

import (
	"api-go/internal/models"

	"gorm.io/gorm"
)

type NotesRepository struct {
	DB *gorm.DB
}

func NewNotesRepository(db *gorm.DB) *NotesRepository {
	return &NotesRepository{
		DB: db,
	}
}

func (r *NotesRepository) Create(userID, roomID uint, title, content string) (*models.Note, error) {
	note := models.Note{
		UserID:  userID,
		RoomID:  roomID,
		Title:   title,
		Content: content,
	}

	if err := r.DB.Create(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NotesRepository) GetByID(id uint) (*models.Note, error) {
	var note models.Note
	if err := r.DB.Preload("User").Preload("Room").First(&note, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &note, nil
}

func (r *NotesRepository) GetByRoomID(roomID uint) ([]models.Note, error) {
	var notes []models.Note
	if err := r.DB.Preload("User").Where("room_id = ?", roomID).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *NotesRepository) GetByUserID(userID uint) ([]models.Note, error) {
	var notes []models.Note
	if err := r.DB.Preload("Room").Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *NotesRepository) Update(id uint, title, content string) error {
	updates := map[string]interface{}{
		"title":   title,
		"content": content,
	}

	if err := r.DB.Model(&models.Note{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (r *NotesRepository) Delete(id uint) error {
	if err := r.DB.Delete(&models.Note{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (r *NotesRepository) GetByUserAndRoom(userID, roomID uint) ([]models.Note, error) {
	var notes []models.Note
	if err := r.DB.Where("user_id = ? AND room_id = ?", userID, roomID).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *NotesRepository) GetAll() ([]models.Note, error) {
	var notes []models.Note
	if err := r.DB.Preload("User").Preload("Room").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}
