package helper

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GenerateJWT(signKey string, userID uint) map[string]any {
	res := map[string]any{}
	accessToken := generateToken(signKey, userID)
	if accessToken == "" {
		return nil
	}
	res["access_token"] = accessToken
	return res
}

func generateToken(signKey string, id uint) string {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := token.SignedString([]byte(signKey))
	if err != nil {
		return ""
	}
	return validToken
}

func ExtractToken(name string, c echo.Context) jwt.MapClaims {
	user := c.Get(name).(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}
