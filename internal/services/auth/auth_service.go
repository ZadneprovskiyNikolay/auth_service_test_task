package auth

import (
	"auth/internal/config"
	jwtutils "auth/internal/utils/jwt"
	logutils "auth/internal/utils/log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	refreshTokenStorage  RefreshTokenStorage
	emailService         EmailService
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	jwtPrivateKey        []byte
	Emails               config.Emails
}

//go:generate mockery --name RefreshTokenStorage --filename refresh_token_storage.go
type RefreshTokenStorage interface {
	Create(token *RefreshToken) (uuid.UUID, error)
	Get(id uuid.UUID) (*RefreshToken, error)
	Delete(id uuid.UUID) error
}

//go:generate mockery --name EmailService --filename email_service.go
type EmailService interface {
	SendEmailToUser(from string, userID uuid.UUID, msg []byte) error
}

func NewAuthService(
	refershTokenStorage RefreshTokenStorage,
	emailService EmailService,
	jwtPrivateKey []byte,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
	emails config.Emails) *AuthService {
	return &AuthService{
		refreshTokenStorage:  refershTokenStorage,
		emailService:         emailService,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		jwtPrivateKey:        jwtPrivateKey,
		Emails:               emails,
	}
}

func (s *AuthService) CreateAccessAndRefreshTokens(userID uuid.UUID, requestIP string) (string, string, error) {
	refreshBytes, err := generateRefreshTokenBytes()
	if err != nil {
		return "", "", errors.Wrap(err, "generate refresh token bytes")
	}
	refreshHash, err := bcrypt.
		GenerateFromPassword(refreshBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", "", errors.Wrap(err, "hash refresh token")
	}

	refreshExpTime := time.Now().Add(s.refreshTokenDuration)
	refresh := &RefreshToken{
		Hash:      refreshHash,
		ExpiresAt: refreshExpTime,
	}
	refreshTokenID, err := s.refreshTokenStorage.Create(refresh)
	if err != nil {
		return "", "", errors.Wrap(err, "create refresh token")
	}

	accessExpTime := time.Now().Add(s.accessTokenDuration).Unix()
	accessJWT := jwt.NewWithClaims(jwtSigningMethod,
		jwt.MapClaims{
			UserIDClaim:         userID.String(),
			UserIPClaim:         requestIP,
			ExpTimeClaim:        accessExpTime,
			RefreshTokenIDClaim: refreshTokenID,
		})
	access, err := accessJWT.SignedString(s.jwtPrivateKey)
	if err != nil {
		return "", "", errors.Wrap(err, "sign access token")
	}

	return access, string(refreshBytes), nil
}

func (s *AuthService) RefreshAccessToken(accessToken, refreshTokenStr, requestIP string) (string, string, error) {
	jwtToken, err := jwtutils.ParseAndValidateJWTToken(accessToken, s.jwtPrivateKey, jwtSigningMethod.Name)
	if err != nil {
		return "", "", UnauthorizedError{"parse error"}
	}

	jwtClaims, err := parseJWTClaims(jwtToken)
	if err != nil {
		return "", "", UnauthorizedError{}
	}

	refreshToken, err := s.refreshTokenStorage.Get(jwtClaims.refreshTokenID)
	if err != nil {
		return "", "", UnauthorizedError{}
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return "", "", UnauthorizedError{"refresh token expired"}
	}

	if bcrypt.CompareHashAndPassword(refreshToken.Hash, []byte(refreshTokenStr)) != nil {
		return "", "", UnauthorizedError{}
	}

	if requestIP != jwtClaims.userIP {
		err := s.emailService.SendEmailToUser(
			s.Emails.SupportEmail,
			jwtClaims.userID,
			RefreshRequestNewIPEmail(requestIP))
		if err != nil {
			logutils.Error("send email error", err)
		}
	}

	go func() {
		err := s.refreshTokenStorage.Delete(refreshToken.ID)
		if err != nil {
			logutils.Error("delete old refresh token error", err)
		}
	}()

	return s.CreateAccessAndRefreshTokens(
		jwtClaims.userID,
		requestIP,
	)
}
