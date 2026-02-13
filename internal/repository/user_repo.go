package repository

import (
	"github.com/sm2-cosign/backend/internal/model"
)

// UserRepository 用户数据访问
type UserRepository struct{}

// NewUserRepository 创建用户数据访问实例
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (id, username, password_hash, public_key, status, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))`
	_, err := db.Exec(query, user.ID, user.Username, user.PasswordHash, user.PublicKey, user.Status)
	return err
}

// FindByID 根据ID查询用户
func (r *UserRepository) FindByID(id string) (*model.User, error) {
	query := `SELECT id, username, password_hash, public_key, status, created_at, updated_at FROM users WHERE id = ?`
	user := &model.User{}
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.PublicKey,
		&user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByUsername 根据用户名查询用户
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, password_hash, public_key, status, created_at, updated_at FROM users WHERE username = ?`
	user := &model.User{}
	err := db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.PublicKey,
		&user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// List 获取用户列表
func (r *UserRepository) List(page, pageSize int) ([]model.User, int64, error) {
	offset := (page - 1) * pageSize
	
	// 获取总数
	var total int64
	err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	query := `SELECT id, username, password_hash, public_key, status, created_at, updated_at 
	          FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID, &user.Username, &user.PasswordHash, &user.PublicKey,
			&user.Status, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	return users, total, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	query := `UPDATE users SET password_hash = ?, public_key = ?, status = ?, updated_at = datetime('now') WHERE id = ?`
	_, err := db.Exec(query, user.PasswordHash, user.PublicKey, user.Status, user.ID)
	return err
}

// UpdateStatus 更新用户状态
func (r *UserRepository) UpdateStatus(id string, status int) error {
	query := `UPDATE users SET status = ?, updated_at = datetime('now') WHERE id = ?`
	_, err := db.Exec(query, status, id)
	return err
}

// Delete 删除用户
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&count)
	return count > 0, err
}
