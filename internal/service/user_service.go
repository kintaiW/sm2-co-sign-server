package service

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/sm2-cosign/backend/internal/config"
	"github.com/sm2-cosign/backend/internal/crypto"
	"github.com/sm2-cosign/backend/internal/model"
	"github.com/sm2-cosign/backend/internal/repository"
	"github.com/sm2-cosign/backend/pkg/response"
	"github.com/sm2-cosign/backend/pkg/utils"
)

// UserService 用户服务
type UserService struct {
	userRepo    *repository.UserRepository
	keyRepo     *repository.KeyRepository
	sessionRepo *repository.SessionRepository
	auditRepo   *repository.AuditLogRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo:    repository.NewUserRepository(),
		keyRepo:     repository.NewKeyRepository(),
		sessionRepo: repository.NewSessionRepository(),
		auditRepo:   repository.NewAuditLogRepository(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=6,max=64"`
	P1       string `json:"p1" validate:"required"`
}

type RegisterResponse struct {
	UserID    string `json:"userId"`
	PublicKey string `json:"publicKey"`
	P2        string `json:"p2"`
}

// Register 用户注册
func (s *UserService) Register(req *RegisterRequest, ipAddress string) (*RegisterResponse, response.Code) {
	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, response.CodeDBError
	}
	if exists {
		return nil, response.CodeUserExists
	}

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

	// 生成用户ID
	userID := utils.GenerateUUID()

	// 生成密码哈希
	salt, err := utils.GenerateSalt()
	if err != nil {
		return nil, response.CodeInternalError
	}
	passwordHash := crypto.SM3HashWithPassword([]byte(req.Password), salt)

	// 创建用户
	user := &model.User{
		ID:           userID,
		Username:     req.Username,
		PasswordHash: hex.EncodeToString(salt) + hex.EncodeToString(passwordHash),
		PublicKey:    crypto.EncodeToBase64(keyResult.Pa),
		Status:       model.UserStatusEnabled,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, response.CodeDBError
	}

	// 创建密钥记录
	key := &model.Key{
		ID:        utils.GenerateUUID(),
		UserID:    userID,
		D2:        crypto.EncodeToBase64(keyResult.D2),
		D2Inv:     crypto.EncodeToBase64(keyResult.D2Inv),
		PublicKey: crypto.EncodeToBase64(keyResult.Pa),
		Status:    model.KeyStatusEnabled,
	}
	if err := s.keyRepo.Create(key); err != nil {
		return nil, response.CodeDBError
	}

	// 记录审计日志
	auditLog := &model.AuditLog{
		ID:        utils.GenerateUUID(),
		UserID:    userID,
		Action:    model.ActionRegister,
		Detail:    `{"username":"` + req.Username + `"}`,
		IPAddress: ipAddress,
	}
	s.auditRepo.Create(auditLog)

	return &RegisterResponse{
		UserID:    userID,
		PublicKey: crypto.EncodeToBase64(keyResult.Pa),
		P2:        crypto.EncodeToBase64(keyResult.P2),
	}, response.CodeSuccess
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
	UserID    string `json:"userId"`
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest, ipAddress string) (*LoginResponse, response.Code) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, response.CodeUserNotFound
	}

	// 验证密码
	storedHash := user.PasswordHash
	if len(storedHash) < 64 {
		return nil, response.CodePasswordError
	}
	salt, _ := hex.DecodeString(storedHash[:32])
	expectedHash := storedHash[32:]
	actualHash := hex.EncodeToString(crypto.SM3HashWithPassword([]byte(req.Password), salt))
	if actualHash != expectedHash {
		return nil, response.CodePasswordError
	}

	// 检查用户状态
	if !user.IsEnabled() {
		return nil, response.CodeUserDisabled
	}

	// 生成 Token
	token, err := utils.GenerateToken()
	if err != nil {
		return nil, response.CodeInternalError
	}

	// 计算过期时间
	expiresAt := utils.CalculateTokenExpiry(config.AppConfig.Auth.TokenExpire)

	// 创建会话
	session := &model.Session{
		ID:        token,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, response.CodeDBError
	}

	// 记录审计日志
	auditLog := &model.AuditLog{
		ID:        utils.GenerateUUID(),
		UserID:    user.ID,
		Action:    model.ActionLogin,
		IPAddress: ipAddress,
	}
	s.auditRepo.Create(auditLog)

	return &LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
		UserID:    user.ID,
	}, response.CodeSuccess
}

// Logout 用户登出
func (s *UserService) Logout(token string) response.Code {
	// 删除会话
	if err := s.sessionRepo.Delete(token); err != nil {
		return response.CodeDBError
	}
	return response.CodeSuccess
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID string) (*model.User, response.Code) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, response.CodeUserNotFound
	}
	return user, response.CodeSuccess
}

// ValidateSession 验证会话
func (s *UserService) ValidateSession(token string) (*model.Session, response.Code) {
	session, err := s.sessionRepo.FindByID(token)
	if err != nil {
		return nil, response.CodeTokenInvalid
	}
	if session.IsExpired() {
		s.sessionRepo.Delete(token)
		return nil, response.CodeTokenExpired
	}
	return session, response.CodeSuccess
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID string) response.Code {
	// 删除用户（级联删除密钥和会话）
	if err := s.userRepo.Delete(userID); err != nil {
		return response.CodeDBError
	}
	return response.CodeSuccess
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(userID string, status int) response.Code {
	if err := s.userRepo.UpdateStatus(userID, status); err != nil {
		return response.CodeDBError
	}
	return response.CodeSuccess
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(page, pageSize int) ([]model.User, int64, response.Code) {
	users, total, err := s.userRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, response.CodeDBError
	}
	return users, total, response.CodeSuccess
}

var (
	ErrInvalidP1 = errors.New("invalid P1 format")
	ErrInvalidQ1 = errors.New("invalid Q1 format")
	ErrInvalidE  = errors.New("invalid E format")
	ErrInvalidT1 = errors.New("invalid T1 format")
)
