package models

type User struct {
	BaseModel
	Name     string `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Email    string `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password string `json:"-" gorm:"not null" validate:"required,min=6"`
}

func (User) TableName() string {
	return "users"
}
