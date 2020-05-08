package auth

import (
	"fmt"

	"github.com/Yesterday17/pug-backend/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authorize(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)

	auth, err := c.Request.Cookie("pug_session")
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	token, err := jwt.Parse(auth.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return cfg.KeyPublic, nil
	})
	if token == nil || err != nil {
		c.AbortWithStatus(401)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatus(401)
		return
	}

	c.Set("uuid", claims["sub"].(string))
	c.Set("level", claims["level"].(int))
	c.Next()
}
