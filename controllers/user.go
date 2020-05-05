package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UserLogin(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	_ = db
}

func UserLogout(c *gin.Context) {

}

func UserRegister(c *gin.Context) {

}
