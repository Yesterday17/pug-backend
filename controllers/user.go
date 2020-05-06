package controllers

import (
	e "github.com/Yesterday17/pug-backend/error"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UserInfoGet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

	var user models.User
	db.First(&user, "uuid = ?", uuid)
	if db.Error != nil {
		c.Abort()
		c.JSON(500, e.ErrDBRead)
		return
	}

	c.JSON(200, user)
}
