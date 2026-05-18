package jwt

import (
	"encoding/json"
	"time"

	pnerrors "github.com/PuppyNote/puppynote-server-golang/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenDuration  = 2 * time.Hour
	refreshTokenDuration = 30 * 24 * time.Hour
)

type LoginUserInfo struct {
	UserID int64  `json:"userId"`
	Email  string `json:"email"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	jwt.RegisteredClaims
}

var secretKey []byte

func Init(secret string) {
	secretKey = []byte(secret)
}

func GeneratePair(info LoginUserInfo) (*TokenPair, error) {
	accessToken, err := generate(info, accessTokenDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generate(info, refreshTokenDuration)
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func generate(info LoginUserInfo, duration time.Duration) (string, error) {
	subject, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   string(subject),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(secretKey)
}

func Parse(tokenString string) (*LoginUserInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pnerrors.ErrTokenInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		if err.Error() == jwt.ErrTokenExpired.Error() {
			return nil, pnerrors.ErrTokenExpired
		}
		return nil, pnerrors.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, pnerrors.ErrTokenInvalid
	}

	var info LoginUserInfo
	if err = json.Unmarshal([]byte(claims.Subject), &info); err != nil {
		return nil, pnerrors.ErrTokenInvalid
	}

	return &info, nil
}
