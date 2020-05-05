package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UserLogin(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	username := c.Request.Form.Get("username")
	password := c.Request.Form.Get("password")

	if username == "" || password == "" {
		c.Abort()
		c.Status(400)
		return
	}

	_ = db
}

func UserLogout(c *gin.Context) {

}

func UserRegister(c *gin.Context) {

}
