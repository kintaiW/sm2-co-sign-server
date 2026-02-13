package service

import (
	"encoding/hex"
	"log"

	"github.com/sm2-cosign/backend/internal/config"
	"github.com/sm2-cosign/backend/internal/crypto"
	"github.com/sm2-cosign/backend/internal/model"
	"github.com/sm2-cosign/backend/internal/repository"
	"github.com/sm2-cosign/backend/pkg/utils"
)

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "please-change-password"
)

func getAdminCredentials() (username, password string) {
	username = DefaultAdminUsername
	password = DefaultAdminPassword

	if config.AppConfig != nil && config.AppConfig.Admin.Username != "" {
		username = config.AppConfig.Admin.Username
	}
	if config.AppConfig != nil && config.AppConfig.Admin.Password != "" {
		password = config.AppConfig.Admin.Password
	}

	return username, password
}

func InitAdminUser() error {
	username, password := getAdminCredentials()

	userRepo := repository.NewUserRepository()

	exists, err := userRepo.ExistsByUsername(username)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Admin user already exists")
		return nil
	}

	salt, err := utils.GenerateSalt()
	if err != nil {
		return err
	}
	passwordHash := crypto.SM3HashWithPassword([]byte(password), salt)

	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		return err
	}

	publicKeyBytes := crypto.EncodeToBase64(append(keyPair.PublicKey.X.Bytes(), keyPair.PublicKey.Y.Bytes()...))

	userID := utils.GenerateUUID()
	user := &model.User{
		ID:           userID,
		Username:     username,
		PasswordHash: hex.EncodeToString(salt) + hex.EncodeToString(passwordHash),
		PublicKey:    publicKeyBytes,
		Status:       model.UserStatusEnabled,
	}
	if err := userRepo.Create(user); err != nil {
		return err
	}

	log.Printf("Admin user created successfully: %s", username)
	return nil
}
