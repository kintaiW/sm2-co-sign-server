package model

import "time"

type Key struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	D2        string    `json:"-" db:"d2"`
	D2Inv     string    `json:"-" db:"d2_inv"`
	PublicKey string    `json:"publicKey" db:"public_key"`
	HMACKey   string    `json:"-" db:"hmac_key"`
	Status    int       `json:"status" db:"status"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// KeyStatus 密钥状态常量
const (
	KeyStatusDisabled = 0
	KeyStatusEnabled  = 1
)

// IsEnabled 检查密钥是否启用
func (k *Key) IsEnabled() bool {
	return k.Status == KeyStatusEnabled
}
