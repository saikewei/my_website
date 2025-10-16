package auth

import (
	"errors"

	"github.com/saikewei/my_website/back/internal/database"
	"github.com/saikewei/my_website/back/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrPasswordNotFound     = errors.New("password not found")
	ErrIncoreectOldPassword = errors.New("incorrect old password")
)

func ChangePasswordStore(req *ChangePasswordRequest) error {
	var currentPassword model.SystemPassword

	if err := database.DB.First(&currentPassword).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPasswordNotFound
		}
		return err
	}

	err := bcrypt.CompareHashAndPassword([]byte(currentPassword.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		return ErrIncoreectOldPassword
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	result := database.DB.Model(&model.SystemPassword{}).Where("id = ?", currentPassword.ID).Update("password_hash", string(newHashedPassword))
	return result.Error
}

func createPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newPassword := model.SystemPassword{
		PasswordHash: string(hashedPassword),
	}

	result := database.DB.Create(&newPassword)
	return result.Error
}

func loginStore(req *LoginRequest) error {
	var currentPassword model.SystemPassword
	if err := database.DB.First(&currentPassword).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPasswordNotFound
		}
		return err
	}

	err := bcrypt.CompareHashAndPassword([]byte(currentPassword.PasswordHash), []byte(req.Password))
	if err != nil {
		return ErrIncoreectOldPassword
	}

	return nil
}
