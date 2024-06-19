package email

import (
	"auth/internal/config"
	"net/smtp"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type EmailService struct {
	smtpConfig  config.SMTP
	userStorage UserStorage
	auth        smtp.Auth
}

type UserStorage interface {
	GetUserEmail(userID uuid.UUID) (string, error)
}

func NewEmailService(
	smtpConfig config.SMTP,
	userStorage UserStorage,
) *EmailService {
	auth := smtp.PlainAuth("", smtpConfig.UserName, smtpConfig.Password, smtpConfig.Host)
	return &EmailService{
		smtpConfig:  smtpConfig,
		userStorage: userStorage,
		auth:        auth,
	}
}

func (s *EmailService) SendEmail(from string, to []string, msg []byte) error {
	return smtp.SendMail(s.smtpConfig.Host+s.smtpConfig.Port, s.auth, from, to, msg)
}

func (s *EmailService) SendEmailToUser(from string, userID uuid.UUID, msg []byte) error {
	to, err := s.userStorage.GetUserEmail(userID)
	if err != nil {
		return errors.Wrap(err, "get user email")
	}

	return s.SendEmail(from, []string{to}, msg)
}
