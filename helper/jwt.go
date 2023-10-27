package helper

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const SECRET_JWT = "mytodo"

type JWT_TOKEN struct {
	Token string
}

func CreateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_JWT))
}

// func ExtractToken(token *jwt.Token) any{
// 	if token.Valid {
// 		var claims = token.Claims
// 		expTime, _ := claims.GetExpirationTime()
// 		fmt.Println(expTime.Time.Compare(time.Now()))
// 		if expTime.Time.Compare(time.Now()) > 0 {

// 			return token.Claims
// 		}

// 		logrus.Error("Token expired")
// 		return nil

// 	}
// 	return nil
// }
