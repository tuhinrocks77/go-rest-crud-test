package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

type User struct {
	Id   int
	Name string
}

const DummyUserId = 1
const DummyUserName = "dummy user"

const jwtSecretKey = "dummy-secret-key" // TODO: load from .env

func MakeDummyUserToken() string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Id":   DummyUserId,
			"Name": DummyUserName,
		})
	s, _ := t.SignedString([]byte(jwtSecretKey))
	fmt.Printf("new token in case old one expires: %v \n", s)
	bearerPrefix := "Bearer "
	return bearerPrefix + s
}

func ValidateToken(tokenString string) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// TODO: find why const jwtSecretKey can't be used when running test cases
			return []byte("dummy-secret-key"), nil
		})
	if err != nil {
		return false
	}
	claims = token.Claims.(jwt.MapClaims)
	return claims["Id"].(float64) == DummyUserId && claims["Name"].(string) == DummyUserName
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := authHeader{}
		errText := "Please provide valid Authorization header"
		if err := ctx.ShouldBindHeader(&h); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": errText,
			})
			ctx.Abort()
		}

		tokenLength := 126
		bearerPrefix := "Bearer "
		token := h.Token
		if len(token) == 0 {
			// TODO: implement common error function
			errText = "Please provide Authorization token"
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": errText,
			})
			ctx.Abort()
			return
		} else if len(token) != tokenLength {
			errText = "Invalid Authorization token"
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": errText,
			})
			ctx.Abort()
			return
		} else if !strings.HasPrefix(token, bearerPrefix) {
			errText = "Invalid Authorization token"
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": errText,
			})
			ctx.Abort()
			return
		}
		token = strings.TrimPrefix(token, bearerPrefix)
		if !ValidateToken(token) {
			errText = "Invalid Authorization token"
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": errText,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
