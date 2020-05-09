package controllers

import (
	"strconv"
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
	// No Restriction
	return 0
}

func canUseModule(module string, level int) bool {
	return canUsePipe(module, "", level)
}

func canUsePipe(module, pipe string, level int) bool {
	return restrictionLevel(module, pipe) <= level
}

func setModulePipeRestriction(db *gorm.DB, module, pipe string, level int) error {
	m, ok := ModulePipeRestriction.Load(module)
	if !ok {
		m = &sync.Map{}
		ModulePipeRestriction.Store(module, m)
	}

	var rule models.ModuleRestrictRule
	if pipe == "" {
		// Module level restriction
		if db != nil {
			// Write to db first
			db.First(&rule, "module_name = ?", module)
			if db.Error != nil {
				return &e.ErrDBRead
			}
			rule.MinAvailableUserLevel = level
			db.Save(&rule)
			if db.Error != nil {
				return &e.ErrDBWrite
			}
		}

		m.(*sync.Map).Store(ModuleLevelKeyword, level)
	} else {
		// Pipe level restriction
		if db != nil {
			// Write to db first
			db.First(&rule, "module_name = ? AND pipe_name = ?", module, pipe)
			if db.Error != nil {
				return &e.ErrDBRead
			}
			rule.MinAvailableUserLevel = level
			db.Save(&rule)
			if db.Error != nil {
				return &e.ErrDBWrite
			}
		}
		m.(*sync.Map).Store(pipe, level)
	}

	return nil
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
			_ = setModulePipeRestriction(nil, r.ModuleName, r.PipeName, r.MinAvailableUserLevel)
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
	db := c.MustGet("db").(*gorm.DB)
	mgr := c.MustGet("manager").(api.ModuleManager)
	level := c.MustGet("level").(int)

	ul, ok := c.GetPostForm("level")
	if !ok {
		ul = "0"
	}

	userLevel, err := strconv.Atoi(ul)
	if err != nil {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}

	if userLevel < 0 {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	} else if userLevel > level {
		c.AbortWithStatusJSON(400, e.ErrCannotRestrictSelf)
		return
	}

	module := c.Param("module")
	pipe := c.Param("pipe")

	// module does not exist
	if m := mgr.Module(module); m == nil {
		c.AbortWithStatusJSON(404, e.ErrModuleNotFound)
		return
	} else if pipe != "" {
		// pipe param found but pipe not exist
		if p := mgr.Pipe(module, pipe); p == nil {
			c.AbortWithStatusJSON(404, e.ErrPipeNotFound)
			return
		}
	}

	// cannot edit level
	if !canUsePipe(module, pipe, level) {
		c.AbortWithStatusJSON(403, e.ErrPermissionDeny)
		return
	}

	err = setModulePipeRestriction(db, module, pipe, userLevel)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	var ret []string
	if pipe == "" {
		// Return available Modules
		for _, m := range mgr.Modules() {
			if canUseModule(m, level) {
				ret = append(ret, m)
			}
		}
	} else {
		// Return available Pipes
		for _, p := range mgr.Module(module).Pipes() {
			if canUsePipe(module, p, level) {
				ret = append(ret, p)
			}
		}
	}
	c.JSON(200, ret)
}
