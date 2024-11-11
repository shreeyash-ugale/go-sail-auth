package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Username string    `json:"username" gorm:"unique"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"password"`
	APIKey   []*APIKey `gorm:"foreignKey:UserID;"`
	PlanID   uuid.UUID `json:"plan_id"`
	Plan     Plan      `gorm:"foreignKey:PlanID"`
}

type APIKey struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Key    string    `gorm:"unique;not null"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User   User      `json:"user" `
}

type Action struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"type:text"`
	Plan        []*Plan   `gorm:"many2many:plan_actions;"`
}

type Plan struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"type:text"`
	Actions     []Action  `gorm:"many2many:plan_actions;"`
	Users       []User    `gorm:"many2many:user_plans;"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
