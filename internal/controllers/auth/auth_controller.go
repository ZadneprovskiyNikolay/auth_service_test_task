package authcontroller

import (
	"auth/internal/controllers"
	"auth/internal/controllers/httputils"
	"auth/internal/services/auth"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var _ controllers.Controller = &AuthController{}

type AuthController struct {
	authService AuthService
}

type AuthService interface {
	CreateAccessAndRefreshTokens(userID uuid.UUID, requestIP string) (string, string, error)
	RefreshAccessToken(accessToken, refreshToken, requestIP string) (string, string, error)
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) createSession(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.URL.Query().Get("userID"))
	if err != nil {
		httputils.BadRequest(w, r, errors.Wrap(err, "parse userID"))
		return
	}

	accessToken, refreshToken, err := c.authService.CreateAccessAndRefreshTokens(
		userID,
		httputils.RequestIP(r))
	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	if err != nil {
		slog.Error(errors.Wrap(err, "create access and refresh tokens").Error())
		httputils.InternalError(w, r)
		return
	}

	render.JSON(w, r, Session{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenBase64,
	})
}

func (c *AuthController) refreshSession(w http.ResponseWriter, r *http.Request) {
	var session Session
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		httputils.BadRequest(w, r, err)
		return
	}
	refreshTokenDecoded, err := base64.StdEncoding.DecodeString(session.RefreshToken)
	if err != nil {
		httputils.BadRequest(w, r, errors.Wrap(err, "decode base64 refresh token"))
	}

	accessToken, refreshToken, err := c.authService.RefreshAccessToken(
		session.AccessToken,
		string(refreshTokenDecoded),
		httputils.RequestIP(r))
	switch err.(type) {
	case nil:
	case auth.UnauthorizedError:
		httputils.UnauthorizedError(w, r, err)
		return
	default:
		httputils.InternalError(w, r)
		return
	}

	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	render.JSON(w, r, Session{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenBase64,
	})
}

func (c *AuthController) RegisterRoutes(router chi.Router) {
	router.Route("/session", func(r chi.Router) {
		r.Get("/", c.createSession)
		r.Post("/refresh", c.refreshSession)
	})
}
