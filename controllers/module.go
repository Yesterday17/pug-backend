package controllers

import (
	"sync"

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

var ModulePipeRestriction *sync.Map

const ModuleLevelKeyword = "__PUG_BACKEND__"

func restrictionLevel(module, pipe string) int {
	m, ok := ModulePipeRestriction.Load(module)
	if ok {
		if pipe != "" {
			// check whether a pipe has been restricted
			if p, ok := m.(*sync.Map).Load(pipe); ok {
				return p.(int)
			}
		} else if g, ok := m.(*sync.Map).Load(ModuleLevelKeyword); ok {
			// check whether a module has been restricted
			return g.(int)
		}
	}
	return 0
}

func canUseModule(module string, level int) bool {
	return canUsePipe(module, "", level)
}

func canUsePipe(module, pipe string, level int) bool {
	return restrictionLevel(module, pipe) <= level
}

func InitModulePipeRestriction(c *gin.Context) {
	if ModulePipeRestriction == nil {
		db := c.MustGet("db").(*gorm.DB)

		var restrictions []models.ModuleRestrictRule
		db.Find(&restrictions)
		if db.Error != nil {
			c.AbortWithStatusJSON(500, e.ErrDBRead)
		}

		ModulePipeRestriction = &sync.Map{}

		for _, r := range restrictions {
			m, ok := ModulePipeRestriction.Load(r.ModuleName)
			if !ok {
				m = &sync.Map{}
				ModulePipeRestriction.Store(r.ModuleName, m)
			}

			if r.PipeName == "" {
				// Module level restriction
				m.(*sync.Map).Store(ModuleLevelKeyword, r.MinAvailableUserLevel)
			} else {
				m.(*sync.Map).Store(r.PipeName, r.MinAvailableUserLevel)
			}
		}
	}
}

func GetAllModuleInfo(c *gin.Context) {
	mgr := c.MustGet("manager").(api.ModuleManager)
	level := c.MustGet("level").(int)

	var ret []ModuleInfo

	modules := mgr.Modules()
	for _, name := range modules {
		if canUseModule(name, level) {
			m := mgr.Module(name)
			ret = append(ret, ModuleInfo{
				Name:        m.Name(),
				Description: m.Description(),
				Author:      m.Author(),
				Usage:       m.Usage(),
			})
		}
	}

	c.JSON(200, ret)
}

func GetModuleInfo(c *gin.Context) {
	mgr := c.MustGet("manager").(api.ModuleManager)
	level := c.MustGet("level").(int)
	name := c.Param("module")

	if !canUseModule(name, level) {
		c.AbortWithStatusJSON(403, e.ErrPermissionDeny)
		return
	}

	m := mgr.Module(name)
	if m == nil {
		c.AbortWithStatusJSON(404, e.ErrModuleNotFound)
		return
	}

	c.JSON(200, ModuleInfo{
		Name:        m.Name(),
		Description: m.Description(),
		Author:      m.Author(),
		Usage:       m.Usage(),
	})
}

func EditModulePipeRestriction(c *gin.Context) {

}
