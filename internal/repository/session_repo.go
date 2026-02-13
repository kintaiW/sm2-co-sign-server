package repository

import (
	"time"

	"github.com/sm2-cosign/backend/internal/model"
)

// SessionRepository 会话数据访问
type SessionRepository struct{}

// NewSessionRepository 创建会话数据访问实例
func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

// Create 创建会话
func (r *SessionRepository) Create(session *model.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at, created_at) 
	          VALUES (?, ?, ?, datetime('now'))`
	_, err := db.Exec(query, session.ID, session.UserID, session.ExpiresAt)
	return err
}

// FindByID 根据ID查询会话
func (r *SessionRepository) FindByID(id string) (*model.Session, error) {
	query := `SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ?`
	session := &model.Session{}
	err := db.QueryRow(query, id).Scan(
		&session.ID, &session.UserID, &session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// FindByUserID 根据用户ID查询会话列表
func (r *SessionRepository) FindByUserID(userID string) ([]model.Session, error) {
	query := `SELECT id, user_id, expires_at, created_at FROM sessions WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.Session
	for rows.Next() {
		var session model.Session
		if err := rows.Scan(
			&session.ID, &session.UserID, &session.ExpiresAt, &session.CreatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

// Delete 删除会话
func (r *SessionRepository) Delete(id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// DeleteByUserID 根据用户ID删除所有会话
func (r *SessionRepository) DeleteByUserID(userID string) error {
	query := `DELETE FROM sessions WHERE user_id = ?`
	_, err := db.Exec(query, userID)
	return err
}

// DeleteExpired 删除过期会话
func (r *SessionRepository) DeleteExpired() error {
	query := `DELETE FROM sessions WHERE expires_at < ?`
	_, err := db.Exec(query, time.Now())
	return err
}

// UpdateExpiresAt 更新会话过期时间
func (r *SessionRepository) UpdateExpiresAt(id string, expiresAt time.Time) error {
	query := `UPDATE sessions SET expires_at = ? WHERE id = ?`
	_, err := db.Exec(query, expiresAt, id)
	return err
}

// CleanupExpired 清理过期会话
func (r *SessionRepository) CleanupExpired() (int64, error) {
	result, err := db.Exec(`DELETE FROM sessions WHERE expires_at < datetime('now')`)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
