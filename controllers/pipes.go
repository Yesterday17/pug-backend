package controllers

import (
	e "github.com/Yesterday17/pug-backend/error"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetAllExtrudedPipeInfo(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

	var p []models.PipeConstructed
	db.Set("gorm:auto_preload", true).Find(&p, "owner = ? OR public = ?", uuid, true)
	if db.Error != nil {
		c.AbortWithStatusJSON(500, e.ErrDBRead)
		return
	}

	c.JSON(200, p)
}

func GetExtrudedPipeInfo(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

	id := c.Param("id")

	var p models.PipeConstructed
	db.Set("gorm:auto_preload", true).First(&p, "id = ?", id)
	if db.Error != nil || p.Owner == "" {
		c.AbortWithStatusJSON(500, e.ErrDBRead)
		return
	}

	if p.Owner != uuid && !p.Public {
		c.AbortWithStatusJSON(401, e.ErrCannotVisitPrivate)
		return
	}

	// Empty user settings
	// p.OwnerUser.Setting = models.UserSettings{}
	c.JSON(200, p)
}

func ExtrudePipe(c *gin.Context) {

}
