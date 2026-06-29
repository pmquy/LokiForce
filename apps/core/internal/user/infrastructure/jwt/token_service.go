package jwt

import (
	"time"

	"lokiforce.com/apps/core/internal/user/application"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey string
}

func NewJWTService(secret string) application.TokenService {
	return &jwtService{secretKey: secret}
}

func (j *jwtService) GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	}

	return "", jwt.ErrInvalidKey
}
