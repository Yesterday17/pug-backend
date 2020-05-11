package controllers

import (
	"strconv"

	e "github.com/Yesterday17/pug-backend/error"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/Yesterday17/pug-backend/utils"
	"github.com/Yesterday17/pug/api"
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

	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}

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
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)
	mgr := c.MustGet("manager").(api.ModuleManager)

	module, ok := c.GetPostForm("module")
	if !ok {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	} else if mgr.Module(module) == nil {
		c.AbortWithStatusJSON(400, e.ErrModuleNotFound)
		return
	}

	pipe, ok := c.GetPostForm("pipe")
	if !ok {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	} else if mgr.Pipe(module, pipe) == nil {
		c.AbortWithStatusJSON(400, e.ErrPipeNotFound)
		return
	}

	name, ok := c.GetPostForm("name")
	if !ok {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}
	desc := utils.GetPostForm(c, "description", "")

	pub := utils.GetPostForm(c, "public", "true")
	public, err := strconv.ParseBool(pub)
	if err != nil {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}

	args, ok := c.GetPostForm("arguments")
	if !ok {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}
	arguments := map[string]interface{}{}
	err = json.Unmarshal([]byte(args), &arguments)
	if err != nil {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}

	pc := mgr.Pipe(module, pipe)
	_, err = pc.Build(arguments)
	if err != api.PipeNoError {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}

	db.Create(&models.PipeConstructed{
		ModelIONDP: models.ModelIONDP{
			Owner:       uuid,
			Name:        name,
			Description: desc,
			Public:      public,
		},
		Module:    module,
		Pipe:      pipe,
		Arguments: args,
	})
	if db.Error != nil {
		c.AbortWithStatusJSON(500, e.ErrDBWrite)
		return
	}

	c.JSON(200, e.NoError)
}

func DeleteExtrudedPipe(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}

	var pipe models.PipeConstructed
	db.First(&pipe, "id = ?", id)
	if db.Error != nil || pipe.Owner == "" {
		c.AbortWithStatusJSON(500, e.ErrDBRead)
		return
	}

	if pipe.Owner != uuid {
		if pipe.Public {
			c.AbortWithStatusJSON(401, e.ErrCannotDeleteNotOwnedPipe)
		} else {
			c.AbortWithStatusJSON(401, e.ErrCannotVisitPrivate)
		}
		return
	}

	db.Delete(&models.PipeConstructed{
		ModelIONDP: models.ModelIONDP{ModelWithID: models.ModelWithID{ID: uint(id)}},
	})
	if db.Error != nil {
		c.AbortWithStatusJSON(500, e.ErrDBDelete)
		return
	}

	c.JSON(200, e.NoError)
}
