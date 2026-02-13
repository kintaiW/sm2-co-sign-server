package repository

import (
	"github.com/sm2-cosign/backend/internal/model"
)

// KeyRepository 密钥数据访问
type KeyRepository struct{}

// NewKeyRepository 创建密钥数据访问实例
func NewKeyRepository() *KeyRepository {
	return &KeyRepository{}
}

// Create 创建密钥记录
func (r *KeyRepository) Create(key *model.Key) error {
	query := `INSERT INTO keys (id, user_id, d2, d2_inv, public_key, hmac_key, status, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'))`
	_, err := db.Exec(query, key.ID, key.UserID, key.D2, key.D2Inv, key.PublicKey, key.HMACKey, key.Status)
	return err
}

// FindByID 根据ID查询密钥
func (r *KeyRepository) FindByID(id string) (*model.Key, error) {
	query := `SELECT id, user_id, d2, d2_inv, public_key, hmac_key, status, created_at FROM keys WHERE id = ?`
	key := &model.Key{}
	err := db.QueryRow(query, id).Scan(
		&key.ID, &key.UserID, &key.D2, &key.D2Inv,
		&key.PublicKey, &key.HMACKey, &key.Status, &key.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// FindByUserID 根据用户ID查询密钥
func (r *KeyRepository) FindByUserID(userID string) (*model.Key, error) {
	query := `SELECT id, user_id, d2, d2_inv, public_key, hmac_key, status, created_at FROM keys WHERE user_id = ?`
	key := &model.Key{}
	err := db.QueryRow(query, userID).Scan(
		&key.ID, &key.UserID, &key.D2, &key.D2Inv,
		&key.PublicKey, &key.HMACKey, &key.Status, &key.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// List 获取密钥列表
func (r *KeyRepository) List(page, pageSize int) ([]model.Key, int64, error) {
	offset := (page - 1) * pageSize

	var total int64
	err := db.QueryRow(`SELECT COUNT(*) FROM keys`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, user_id, d2, d2_inv, public_key, hmac_key, status, created_at 
	          FROM keys ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var keys []model.Key
	for rows.Next() {
		var key model.Key
		if err := rows.Scan(
			&key.ID, &key.UserID, &key.D2, &key.D2Inv,
			&key.PublicKey, &key.HMACKey, &key.Status, &key.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		keys = append(keys, key)
	}
	return keys, total, nil
}

// Update 更新密钥
func (r *KeyRepository) Update(key *model.Key) error {
	query := `UPDATE keys SET d2 = ?, d2_inv = ?, public_key = ?, hmac_key = ?, status = ? WHERE id = ?`
	_, err := db.Exec(query, key.D2, key.D2Inv, key.PublicKey, key.HMACKey, key.Status, key.ID)
	return err
}

// UpdateHMACKey 更新HMAC密钥
func (r *KeyRepository) UpdateHMACKey(id, hmacKey string) error {
	query := `UPDATE keys SET hmac_key = ? WHERE id = ?`
	_, err := db.Exec(query, hmacKey, id)
	return err
}

// Delete 删除密钥
func (r *KeyRepository) Delete(id string) error {
	query := `DELETE FROM keys WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// DeleteByUserID 根据用户ID删除密钥
func (r *KeyRepository) DeleteByUserID(userID string) error {
	query := `DELETE FROM keys WHERE user_id = ?`
	_, err := db.Exec(query, userID)
	return err
}
