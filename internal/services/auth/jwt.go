package auth

import (
	"fmt"
	"time"

	jwtutils "auth/internal/utils/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS512
)

const (
	UserIDClaim         = "sub"
	UserIPClaim         = "sub_ip"
	RefreshTokenIDClaim = "refresh_token_id"
	ExpTimeClaim        = "exp"
)

type jwtClaims struct {
	userID         uuid.UUID
	refreshTokenID uuid.UUID
	userIP         string
	expTime        time.Time
}

func parseJWTClaims(token *jwt.Token) (*jwtClaims, error) {
	var claims jwtClaims
	if claimsMap, ok := token.Claims.(jwt.MapClaims); ok {
		userIDStr, err := jwtutils.GetStringJWTClaim(claimsMap, UserIDClaim)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("parse %s claim", UserIDClaim))
		}
		claims.userID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("parse %s claim", UserIDClaim))
		}
		claims.userIP, err = jwtutils.GetStringJWTClaim(claimsMap, UserIPClaim)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("parse %s claim", UserIPClaim))
		}
		refreshTokenID, err := jwtutils.GetStringJWTClaim(claimsMap, RefreshTokenIDClaim)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("parse %s claim", RefreshTokenIDClaim))
		}
		claims.refreshTokenID, err = uuid.Parse(refreshTokenID)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("parse %s claim", RefreshTokenIDClaim))
		}
		claims.expTime, err = jwtutils.GetTimeJWTClaim(claimsMap, ExpTimeClaim)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("parse %s claim", ExpTimeClaim))
		}

		return &claims, nil
	}

	return nil, fmt.Errorf("claim missing")
}
