package service

import (
	"github.com/sm2-cosign/backend/internal/crypto"
	"github.com/sm2-cosign/backend/internal/model"
	"github.com/sm2-cosign/backend/internal/repository"
	"github.com/sm2-cosign/backend/pkg/response"
	"github.com/sm2-cosign/backend/pkg/utils"
)

// CosignService 协同签名服务
type CosignService struct {
	keyRepo   *repository.KeyRepository
	auditRepo *repository.AuditLogRepository
}

// NewCosignService 创建协同签名服务实例
func NewCosignService() *CosignService {
	return &CosignService{
		keyRepo:   repository.NewKeyRepository(),
		auditRepo: repository.NewAuditLogRepository(),
	}
}

// KeyInitRequest 密钥初始化请求
type KeyInitRequest struct {
	UserID string `json:"userId" validate:"required"`
	P1     string `json:"p1" validate:"required"`
}

type KeyInitResponse struct {
	P2        string `json:"p2"`
	PublicKey string `json:"publicKey"`
}

// KeyInit 密钥初始化（重新生成密钥）
func (s *CosignService) KeyInit(req *KeyInitRequest, ipAddress string) (*KeyInitResponse, response.Code) {
	// 解码 P1
	p1, err := crypto.DecodeFromBase64(req.P1)
	if err != nil || len(p1) != 64 {
		return nil, response.CodeInvalidParam
	}

	// 生成协同密钥对
	keyResult, err := crypto.CoopKeyGenInit(p1)
	if err != nil {
		return nil, response.CodeCryptoError
	}

	// 更新或创建密钥记录
	existingKey, _ := s.keyRepo.FindByUserID(req.UserID)
	if existingKey != nil {
		existingKey.D2 = crypto.EncodeToBase64(keyResult.D2)
		existingKey.D2Inv = crypto.EncodeToBase64(keyResult.D2Inv)
		existingKey.PublicKey = crypto.EncodeToBase64(keyResult.Pa)
		if err := s.keyRepo.Update(existingKey); err != nil {
			return nil, response.CodeDBError
		}
	} else {
		key := &model.Key{
			ID:        utils.GenerateUUID(),
			UserID:    req.UserID,
			D2:        crypto.EncodeToBase64(keyResult.D2),
			D2Inv:     crypto.EncodeToBase64(keyResult.D2Inv),
			PublicKey: crypto.EncodeToBase64(keyResult.Pa),
			Status:    model.KeyStatusEnabled,
		}
		if err := s.keyRepo.Create(key); err != nil {
			return nil, response.CodeDBError
		}
	}

	// 记录审计日志
	auditLog := &model.AuditLog{
		ID:        utils.GenerateUUID(),
		UserID:    req.UserID,
		Action:    model.ActionKeyGen,
		IPAddress: ipAddress,
	}
	s.auditRepo.Create(auditLog)

	return &KeyInitResponse{
		P2:        crypto.EncodeToBase64(keyResult.P2),
		PublicKey: crypto.EncodeToBase64(keyResult.Pa),
	}, response.CodeSuccess
}

// SignRequest 签名请求
type SignRequest struct {
	UserID string `json:"userId" validate:"required"`
	Q1     string `json:"q1" validate:"required"`
	E      string `json:"e" validate:"required"`
}

type SignResponse struct {
	R  string `json:"r"`
	S2 string `json:"s2"`
	S3 string `json:"s3"`
}

// Sign 协同签名
func (s *CosignService) Sign(req *SignRequest, ipAddress string) (*SignResponse, response.Code) {
	// 获取密钥
	key, err := s.keyRepo.FindByUserID(req.UserID)
	if err != nil {
		return nil, response.CodeKeyNotFound
	}

	// 解码参数
	q1, err := crypto.DecodeFromBase64(req.Q1)
	if err != nil || len(q1) != 64 {
		return nil, response.CodeInvalidParam
	}

	e, err := crypto.DecodeFromBase64(req.E)
	if err != nil || len(e) != 32 {
		return nil, response.CodeInvalidParam
	}

	// 解码 D2Inv
	d2Inv, err := crypto.DecodeFromBase64(key.D2Inv)
	if err != nil {
		return nil, response.CodeCryptoError
	}

	// 执行协同签名
	result, err := crypto.CoopSign(d2Inv, q1, e)
	if err != nil {
		return nil, response.CodeCryptoError
	}

	// 记录审计日志
	auditLog := &model.AuditLog{
		ID:        utils.GenerateUUID(),
		UserID:    req.UserID,
		Action:    model.ActionSign,
		IPAddress: ipAddress,
	}
	s.auditRepo.Create(auditLog)

	return &SignResponse{
		R:  crypto.EncodeToBase64(result.R),
		S2: crypto.EncodeToBase64(result.S2),
		S3: crypto.EncodeToBase64(result.S3),
	}, response.CodeSuccess
}

// DecryptRequest 解密请求
type DecryptRequest struct {
	UserID string `json:"userId" validate:"required"`
	T1     string `json:"t1" validate:"required"`
}

type DecryptResponse struct {
	T2 string `json:"t2"`
}

// Decrypt 协同解密
func (s *CosignService) Decrypt(req *DecryptRequest, ipAddress string) (*DecryptResponse, response.Code) {
	// 获取密钥
	key, err := s.keyRepo.FindByUserID(req.UserID)
	if err != nil {
		return nil, response.CodeKeyNotFound
	}

	// 解码参数
	t1, err := crypto.DecodeFromBase64(req.T1)
	if err != nil || len(t1) != 64 {
		return nil, response.CodeInvalidParam
	}

	// 解码 D2Inv
	d2Inv, err := crypto.DecodeFromBase64(key.D2Inv)
	if err != nil {
		return nil, response.CodeCryptoError
	}

	// 执行协同解密
	t2, err := crypto.CoopDecrypt(d2Inv, t1)
	if err != nil {
		return nil, response.CodeCryptoError
	}

	// 记录审计日志
	auditLog := &model.AuditLog{
		ID:        utils.GenerateUUID(),
		UserID:    req.UserID,
		Action:    model.ActionDecrypt,
		IPAddress: ipAddress,
	}
	s.auditRepo.Create(auditLog)

	return &DecryptResponse{
		T2: crypto.EncodeToBase64(t2),
	}, response.CodeSuccess
}

// GetKeyByUserID 根据用户ID获取密钥信息
func (s *CosignService) GetKeyByUserID(userID string) (*model.Key, response.Code) {
	key, err := s.keyRepo.FindByUserID(userID)
	if err != nil {
		return nil, response.CodeKeyNotFound
	}
	return key, response.CodeSuccess
}

// DeleteKey 删除密钥
func (s *CosignService) DeleteKey(keyID string) response.Code {
	if err := s.keyRepo.Delete(keyID); err != nil {
		return response.CodeDBError
	}
	return response.CodeSuccess
}

// ListKeys 获取密钥列表
func (s *CosignService) ListKeys(page, pageSize int) ([]model.Key, int64, response.Code) {
	keys, total, err := s.keyRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, response.CodeDBError
	}
	return keys, total, response.CodeSuccess
}
