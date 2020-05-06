package controllers

import (
	e "github.com/Yesterday17/pug-backend/error"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/Yesterday17/pug-backend/utils"
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

func UserSettingGet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

	var settings models.UserSettings
	db.First(&settings, "uuid = ?", uuid)
	if db.Error != nil || settings.UUID == "" {
		c.Abort()
		c.JSON(500, e.ErrDBRead)
		return
	}

	var ret interface{}
	ca := c.Param("category")
	ke := c.Param("key")

	if ca != "" {
		// Get category
		ret = utils.GetFieldByTag(settings, "json", ca)
		if ret == nil {
			c.Abort()
			c.JSON(404, e.ErrInputNotFound)
			return
		}

		if ke != "" {
			// Get key
			ret = utils.GetFieldByTag(ret, "json", ke)
			if ret == nil {
				c.Abort()
				c.JSON(404, e.ErrInputNotFound)
				return
			}
			ret = map[string]interface{}{
				ke: ret,
			}
		}
	}

	c.JSON(200, ret)
}
