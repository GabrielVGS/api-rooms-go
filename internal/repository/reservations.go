package repository

import (
	"api-go/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReservationsRepository struct {
	DB *gorm.DB
}

func NewReservationsRepository(db *gorm.DB) *ReservationsRepository {
	return &ReservationsRepository{
		DB: db,
	}
}

func (r *ReservationsRepository) Create(userID uint, roomID uint, startTime time.Time, endTime time.Time) error {
	reservation := models.Reservation{
		UserID:    userID,
		RoomID:    roomID,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := r.DB.Create(&reservation).Error; err != nil {
		return err
	}
	return nil
}

func (r *ReservationsRepository) GetByID(id uint) (*models.Reservation, error) {
	//
	return nil, nil
}

func (r *ReservationsRepository) Update(id uint, userID uint, roomID uint, startTime time.Time, endTime time.Time) error {
	//
	return nil
}
func (r *ReservationsRepository) Delete(id uint) error {
	return nil
}

func (r *ReservationsRepository) GetByUserID(userID uint) ([]models.Reservation, error) {
	return nil, nil

}

func (r *ReservationsRepository) GetByRoomID(roomID uint) ([]models.Reservation, error) {
	return nil, nil
}
