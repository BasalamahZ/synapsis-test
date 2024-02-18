package service

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/synapsis-test/internal/user"
)

// jwtClaimss is the claims encapsulated in JWT-generated token.
type jwtClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// parseTokenData parse token data from jwt claims.
func (jwtc jwtClaims) parseTokenData() user.TokenData {
	return user.TokenData{
		UserID: jwtc.UserID,
		Email:  jwtc.Email,
	}
}

// formatTokenData format token data into jwt claims.
func formatTokenData(data user.TokenData) jwtClaims {
	return jwtClaims{
		UserID: data.UserID,
		Email:  data.Email,
	}
}

func (s *service) LoginBasic(ctx context.Context, email string, password string) (string, user.TokenData, error) {
	// validate the given values
	if email == "" {
		return "", user.TokenData{}, user.ErrInvalidEmail
	}
	if password == "" {
		return "", user.TokenData{}, user.ErrInvalidPassword
	}

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return "", user.TokenData{}, err
	}

	// get user current data
	current, err := pgStoreClient.GetUserByEmail(ctx, email)
	if err != nil {
		return "", user.TokenData{}, err
	}

	// check password
	err = s.checkPassword(ctx, current, password)
	if err != nil {
		return "", user.TokenData{}, err
	}

	// generate token
	tokenData := user.TokenData{
		UserID: current.ID,
		Email:  current.Email,
	}
	token, err := s.generateToken(ctx, tokenData)
	if err != nil {
		return "", user.TokenData{}, err
	}

	return token, tokenData, nil
}

// TODO: check for expired token error from internal JWT library
func (s *service) ValidateToken(ctx context.Context, token string) (user.TokenData, error) {
	if token == "" {
		return user.TokenData{}, user.ErrInvalidToken
	}

	// get jwt token object
	jwtToken, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.config.TokenSecretKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return user.TokenData{}, user.ErrInvalidToken
		}
		return user.TokenData{}, err
	}

	// check whether token is valid or not (from expirations time)
	if !jwtToken.Valid {
		return user.TokenData{}, user.ErrExpiredToken
	}

	// parse jwt claims
	claims, ok := jwtToken.Claims.(*jwtClaims)
	if !ok {
		return user.TokenData{}, user.ErrInvalidToken
	}

	return claims.parseTokenData(), nil
}

func (s *service) RefreshToken(ctx context.Context, token string) (string, error) {
	// validate token
	data, err := s.ValidateToken(ctx, token)
	if err != nil {
		return "", err
	}

	// generate a new token
	return s.generateToken(ctx, data)
}

// generateToken returns a new token that encapsulates the
// given token data with some additional information:
//   - token expiration time
//
// Token is generated using JWT HS256.
func (s *service) generateToken(ctx context.Context, data user.TokenData) (string, error) {
	claims := formatTokenData(data)

	// add expirations time
	expiresAt := s.timeNow().Add(s.config.TokenExpiration)
	claims.ExpiresAt = jwt.NewNumericDate(expiresAt)

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token with secret key
	signedToken, err := token.SignedString([]byte(s.config.TokenSecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
