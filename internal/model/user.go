package model

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	PublicKey    string    `json:"publicKey" db:"public_key"`
	Status       int       `json:"status" db:"status"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// UserStatus 用户状态常量
const (
	UserStatusDisabled = 0
	UserStatusEnabled  = 1
)

// IsEnabled 检查用户是否启用
func (u *User) IsEnabled() bool {
	return u.Status == UserStatusEnabled
}
