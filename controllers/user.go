package controllers

import (
	e "github.com/Yesterday17/pug-backend/error"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/Yesterday17/pug-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func UserInfoGet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

	var user models.User
	err := db.Set("gorm:auto_preload", true).First(&user, "uuid = ?", uuid).Error
	if err != nil {
		c.JSON(500, e.ErrDBRead)
		return
	}

	c.JSON(200, user)
}

func UserSettingGet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

	var settings models.UserSettings
	db.Set("gorm:auto_preload", true).First(&settings, "uuid = ?", uuid)
	if db.Error != nil || settings.UUID == "" {
		c.JSON(500, e.ErrDBRead)
		return
	}

	var ret interface{}
	ca := c.Query("category")
	ke := c.Query("key")

	if ca != "" {
		// Get category
		ret = utils.GetFieldByTag(settings, "json", ca)
		if ret == nil {
			c.JSON(404, e.ErrInputNotFound)
			return
		}

		if ke != "" {
			// Get key
			ret = utils.GetFieldByTag(ret, "json", ke)
			if ret == nil {
				c.JSON(404, e.ErrInputNotFound)
				return
			}
			ret = map[string]interface{}{
				ke: ret,
			}
		}
	} else {
		ret = settings
	}

	c.JSON(200, ret)
}

func UserSettingPatch(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

	var settings models.UserSettings
	db.Set("gorm:auto_preload", true).First(&settings, "uuid = ?", uuid)
	if db.Error != nil || settings.UUID == "" {
		c.JSON(500, e.ErrDBRead)
		return
	}

	ca, ok := c.GetPostForm("category")
	if !ok || ca == "" {
		c.JSON(400, e.ErrInputInvalid)
		return
	}

	data, ok := c.GetPostForm("data")
	if !ok || ca == "" {
		c.JSON(400, e.ErrInputInvalid)
		return
	}

	var input map[string]interface{}
	err := json.Unmarshal([]byte(data), &input)
	if err != nil {
		c.JSON(400, e.ErrInputInvalid)
		return
	}

	category := utils.GetFieldByTag(&settings, "json", ca)
	if category == nil {
		c.JSON(404, e.ErrInputNotFound)
		return
	}

	var success = 0
	for k, v := range input {
		err = utils.SetFieldByTag(category, "json", k, v)
		if err == nil {
			success++
		}
	}

	// zero success, no one would think it's successful
	if success == 0 {
		c.JSON(400, e.ErrInputInvalid)
		return
	}

	db.Save(&settings)
	if db.Error != nil {
		c.JSON(500, e.ErrDBWrite)
		return
	}

	c.JSON(200, e.NoError)
}
