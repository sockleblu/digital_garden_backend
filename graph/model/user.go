package model

import (
	"time"
)

type WrongUsernameOrPasswordError struct{}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}

type User struct {
	ID         int       `json:"id" gorm:"primary_key"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Token      string    `json:"token"`
	Articles   []*Article `json:"articles" gorm:"foreignKey:ID"`
	CreatedAt time.Time  `json:"createdAt"`
}