package repository

import (
	"database/sql"

	"github.com/sm2-cosign/backend/internal/model"
)

// AuditLogRepository 审计日志数据访问
type AuditLogRepository struct{}

// NewAuditLogRepository 创建审计日志数据访问实例
func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{}
}

// Create 创建审计日志
func (r *AuditLogRepository) Create(log *model.AuditLog) error {
	query := `INSERT INTO audit_logs (id, user_id, action, detail, ip_address, created_at) 
	          VALUES (?, ?, ?, ?, ?, datetime('now'))`
	_, err := db.Exec(query, log.ID, log.UserID, log.Action, log.Detail, log.IPAddress)
	return err
}

// List 获取审计日志列表
func (r *AuditLogRepository) List(page, pageSize int, action, userID string) ([]model.AuditLog, int64, error) {
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if action != "" {
		whereClause += " AND action = ?"
		args = append(args, action)
	}
	if userID != "" {
		whereClause += " AND user_id = ?"
		args = append(args, userID)
	}

	// 获取总数
	countQuery := `SELECT COUNT(*) FROM audit_logs ` + whereClause
	var total int64
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	query := `SELECT id, user_id, action, detail, ip_address, created_at 
	          FROM audit_logs ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	for rows.Next() {
		var log model.AuditLog
		var userID, detail, ipAddress sql.NullString
		if err := rows.Scan(
			&log.ID, &userID, &log.Action, &detail, &ipAddress, &log.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		log.UserID = userID.String
		log.Detail = detail.String
		log.IPAddress = ipAddress.String
		logs = append(logs, log)
	}
	return logs, total, nil
}
