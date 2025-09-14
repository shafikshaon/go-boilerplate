package interfaces

import "go-boilerplate/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll(offset, limit int) ([]*models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
}
