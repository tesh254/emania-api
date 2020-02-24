package services

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

// Claims defines jwt claims
type Claims struct {
	UserID string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken handles generation of a jwt code
func GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(1440 * time.Minute)

	claims := &Claims{
			UserID: userID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

// DecodeToken handles decoding of a jwt token
func DecodeToken(tknStr string) (string, string, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "Token is invalid", err
		}
		return "", "Error decoding token", err
	}

	if !tkn.Valid {
		return "", "Token is Invalid", err
	}

	return claims.UserID, "Lock and Key success", nil
}

// RefreshToken handles refreshing a jwt token
func RefreshToken(tknStr string) (string, string, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error){
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "Token is invalid", err
		}
		return "", "Error decoding token", err
	}

	if !tkn.Valid {
		return "", "Token is invalid", err
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return "", "Token has passed allocated renewal time", nil
	}

	expirationTime := time.Now().Add(5*time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "Error sigining string", err
	}

	return tokenString, "Token sign successful", nil
}
