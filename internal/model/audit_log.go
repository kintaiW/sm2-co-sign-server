package model

import "time"

type AuditLog struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Action    string    `json:"action" db:"action"`
	Detail    string    `json:"detail" db:"detail"`
	IPAddress string    `json:"ipAddress" db:"ip_address"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// AuditAction 审计操作类型常量
const (
	ActionRegister = "register"
	ActionLogin    = "login"
	ActionLogout   = "logout"
	ActionSign     = "sign"
	ActionDecrypt  = "decrypt"
	ActionKeyGen   = "key_gen"
	ActionUserDel  = "user_delete"
	ActionKeyDel   = "key_delete"
)
