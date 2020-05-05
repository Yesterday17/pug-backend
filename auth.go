package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authorize(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := c.Request.Cookie("pub_session")
		if err != nil || !auth.HttpOnly {
			c.Abort()
			c.Status(401)
			return
		}

		token, err := jwt.Parse(auth.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if token == nil || err != nil {
			c.Abort()
			c.Status(401)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.Abort()
			c.Status(401)
			return
		}

		c.Set("uuid", claims["uuid"].(string))
		c.Next()
	}
}
