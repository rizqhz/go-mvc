package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rizghz/api/configs"
)

func CreateToken(userId int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["expired"] = time.Now().Add(time.Hour * 2).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	env, _ := configs.NewJwtEnv()
	return token.SignedString([]byte(env["SECRET_KEY"].(string)))
}

func ExtractToken(ctx echo.Context) int {
	user := ctx.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["userId"].(int)
		return userId
	}
	return 0
}
