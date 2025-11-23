package domain

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre    string    `json:"nombre"`
	Apellido  string    `json:"apellido"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Reparto   string    `json:"reparto"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
