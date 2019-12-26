package utils

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const DEFAULT_COST = 10

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type UserInfo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Salary   int    `json:"salary"`
}

type UserInfoResult struct {
	UserInfo UserInfo `json:"data"`
	Status   bool     `json:"status"`
	Message  string   `json:"message"`
}

type AuthInfo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func HashPassword(password string) (string, error) {
	var hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), DEFAULT_COST)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ClaimsJWT(username string, jwtKey []byte) (string, error) {
	expirationTime := time.Now().Add(100 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CompareHashAndPassword(HashPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(HashPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
