package utils

import "github.com/gin-gonic/gin"

func GetPostForm(c *gin.Context, key, def string) string {
	val, ok := c.GetPostForm(key)
	if !ok {
		return def
	}
	return val
}
