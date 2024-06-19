package auth

import (
	"encoding/base64"
	"testing"
	"time"

	"auth/internal/config"
	"auth/internal/services/auth"
	"auth/internal/services/auth/mocks"
	jwtutils "auth/internal/utils/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	jwtPrivateKey            = []byte("private-key")
	signingMethod            = jwt.SigningMethodHS512
	accessTokenDuration      = time.Hour * 24
	refreshTokenDuration     = time.Hour * 24 * 7
	userID, _                = uuid.Parse("8798e65e-dc84-4a7d-879e-2a52e67d86da")
	accessTokenStr           = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg4NzA2NDEsInJlZnJlc2hfdG9rZW5faWQiOiIzZTAyZWViOS1kZTlhLTRlMGEtODU3Yi0xMjkzYzI1YmQ3NzYiLCJzdWIiOiI4Nzk4ZTY1ZS1kYzg0LTRhN2QtODc5ZS0yYTUyZTY3ZDg2ZGEiLCJzdWJfaXAiOiIxMjcuMC4wLjEifQ.jdtW8iIC3dOqxRb3WItCmhiAK0qnvy91YcBc7e43AHbtEWJolBiisb3cbPvu1QiRqVi6gkSvKKI-Cx8SYkxE4g"
	refreshTokenID, _        = uuid.Parse("3e02eeb9-de9a-4e0a-857b-1293c25bd776")
	refreshTokenBase64Str    = "9W0/xxXxSSSprySP/JRTRQ=="
	refreshTokenStr          = mustNewRefreshTokenFromBase64(refreshTokenBase64Str)
	refreshTokenExpiresAt, _ = time.Parse("2006-01-02 15:04:05.999999", "2024-06-26 20:04:01.115831")
	refreshToken             = auth.RefreshToken{ID: refreshTokenID, Hash: []byte("$2a$10$snAg.KGa.uk0OrGkcp4.Au4pfVfl3pKn5DV8rSz5g5RKkBjHTFI7i"), ExpiresAt: refreshTokenExpiresAt}
	ip                       = "127.0.0.1"
	emails                   = config.Emails{SupportEmail: "support@company.com"}
)

func TestCreateAccessAndRefreshTokens_Simple(t *testing.T) {
	service, refreshTokenStorage, _ := newServiceAndMocks(t)
	refreshTokenStorage.
		On("Create", mock.AnythingOfType("*auth.RefreshToken")).
		Return(uuid.New(), nil)
	startTime := time.Now().Truncate(time.Second) // truncate time since jwt claim "exp" truncates it to seconds

	accessStr, _, err := service.CreateAccessAndRefreshTokens(userID, ip)
	assert.NoError(t, err)
	accessToken, err := jwtutils.ParseAndValidateJWTToken(accessStr, jwtPrivateKey, signingMethod.Name)
	assert.NoError(t, err)
	claimsMap, ok := accessToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	userIDClaim, err := jwtutils.GetStringJWTClaim(claimsMap, auth.UserIDClaim)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), userIDClaim)
	userIPClaim, err := jwtutils.GetStringJWTClaim(claimsMap, auth.UserIPClaim)
	assert.NoError(t, err)
	assert.Equal(t, ip, userIPClaim)
	accessExpTime, err := jwtutils.GetTimeJWTClaim(claimsMap, auth.ExpTimeClaim)
	assert.NoError(t, err)
	assert.True(t, !accessExpTime.Before(startTime.Add(accessTokenDuration))) // user !Before instead of After because time is truncated to seconds and two values can be equal
}

func TestRefreshAccessToken_Simple(t *testing.T) {
	service, refreshTokenStorage, _ := newServiceAndMocks(t)

	refreshTokenStorage.
		On("Get", refreshToken.ID).
		Return(&refreshToken, nil)
	refreshTokenStorage.
		On("Create", mock.AnythingOfType("*auth.RefreshToken")).
		Return(uuid.New(), nil)
	refreshTokenStorage.
		On("Delete", refreshToken.ID).
		Return(nil)

	_, _, err := service.RefreshAccessToken(accessTokenStr, refreshTokenStr, ip)
	assert.NoError(t, err)
}

func TestRefreshAccessToken_TokensDontMatch(t *testing.T) {
	service, refreshTokenStorage, _ := newServiceAndMocks(t)

	refreshTokenStorage.
		On("Get", refreshToken.ID).
		Return(&refreshToken, nil)

	_, _, err := service.RefreshAccessToken(accessTokenStr, refreshTokenStr+"a", ip)
	assert.ErrorAs(t, err, &auth.UnauthorizedError{})
}

func TestRefreshAccessToken_RefreshTokenExpired(t *testing.T) {
	service, refreshTokenStorage, _ := newServiceAndMocks(t)

	refreshToken := refreshToken
	refreshToken.ExpiresAt = time.Now().Add(-time.Second)
	refreshTokenStorage.
		On("Get", refreshToken.ID).
		Return(&refreshToken, nil)

	_, _, err := service.RefreshAccessToken(accessTokenStr, refreshTokenStr, ip)
	assert.ErrorAs(t, err, &auth.UnauthorizedError{})
}

func TestRefreshAccessToken_NewIP(t *testing.T) {
	service, refreshTokenStorage, emailService := newServiceAndMocks(t)

	refreshTokenStorage.
		On("Get", refreshToken.ID).
		Return(&refreshToken, nil)
	refreshTokenStorage.
		On("Create", mock.AnythingOfType("*auth.RefreshToken")).
		Return(uuid.New(), nil)
	refreshTokenStorage.
		On("Delete", refreshToken.ID).
		Return(nil)
	emailService.
		On("SendEmailToUser", emails.SupportEmail, userID, mock.Anything).
		Return(nil)

	_, _, err := service.RefreshAccessToken(accessTokenStr, refreshTokenStr, ip+"1")
	assert.NoError(t, err)
}

func newServiceAndMocks(t *testing.T) (*auth.AuthService, *mocks.RefreshTokenStorage, *mocks.EmailService) {
	trefreshTokenStorage := mocks.NewRefreshTokenStorage(t)
	emailService := mocks.NewEmailService(t)
	service := auth.NewAuthService(
		trefreshTokenStorage,
		emailService,
		jwtPrivateKey,
		accessTokenDuration,
		refreshTokenDuration,
		emails,
	)

	return service, trefreshTokenStorage, emailService
}

func mustNewRefreshTokenFromBase64(base64Str string) string {
	refreshTokenBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		panic(err)
	}
	return string(refreshTokenBytes)
}
