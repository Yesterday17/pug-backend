package controllers

import (
	e "github.com/Yesterday17/pug-backend/error"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/Yesterday17/pug/api"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ModuleInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Author      []string `json:"author"`
	Usage       string   `json:"usage"`
}

var ModulePipeRestriction []models.ModuleRestrictRule = nil

func InitModulePipeRestriction(c *gin.Context) {
	if ModulePipeRestriction == nil {
		db := c.MustGet("db").(*gorm.DB)

		db.Find(&ModulePipeRestriction)
		if db.Error != nil {
			ModulePipeRestriction = nil
			c.Abort()
			c.JSON(500, e.ErrDBRead)
		}
	}
}

func GetAllModuleInfo(c *gin.Context) {
	mgr := c.MustGet("manager").(api.ModuleManager)
	var ret []ModuleInfo

	modules := mgr.Modules()
	for _, name := range modules {
		m := mgr.Module(name)
		ret = append(ret, ModuleInfo{
			Name:        m.Name(),
			Description: m.Description(),
			Author:      m.Author(),
			Usage:       m.Usage(),
		})
	}

	c.JSON(200, ret)
}

func GetModuleInfo(c *gin.Context) {
	mgr := c.MustGet("manager").(api.ModuleManager)

	m := mgr.Module(c.Param("id"))
	if m == nil {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}

	c.JSON(200, ModuleInfo{
		Name:        m.Name(),
		Description: m.Description(),
		Author:      m.Author(),
		Usage:       m.Usage(),
	})
}
