package controllers

import (
	"strconv"
	"strings"

	e "github.com/Yesterday17/pug-backend/error"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/Yesterday17/pug-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetAllPipelines(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

	var ppl []models.Pipeline
	db.Set("gorm:auto_preload", true).Find(&ppl, "owner = ? OR public = ?", uuid, true)
	if db.Error != nil {
		c.AbortWithStatusJSON(500, e.ErrDBRead)
		return
	}

	c.JSON(200, &ppl)
}

func CreatePipeline(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uuid := c.MustGet("uuid").(string)

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

	ps, ok := c.GetPostForm("pipes")
	if !ok {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}
	var pipes []models.PipeConstructed
	for _, p := range strings.Split(ps, ",") {
		pn, err := strconv.ParseUint(p, 10, 0)
		if err != nil {
			c.AbortWithStatusJSON(400, e.ErrInputInvalid)
			return
		}

		var pp models.PipeConstructed
		db.First(&pp, "id = ?", pn)
		if db.Error != nil {
			c.AbortWithStatusJSON(500, e.ErrDBRead)
			return
		} else if pp.Owner == "" {
			c.AbortWithStatusJSON(404, e.ErrConstructedPipeNotFound)
			return
		}

		pipes = append(pipes, pp)
	}

	pl := models.Pipeline{
		Owner:       uuid,
		Name:        name,
		Description: desc,
		Public:      public,
		Pipes:       pipes,
	}
	db.Create(&pl)
	if db.Error != nil {
		c.AbortWithStatusJSON(500, e.ErrDBWrite)
		return
	}

	db.Set("gorm:auto_preload", true).First(&pl, "id = ?", pl.ID)
	c.JSON(200, &pl)
}
