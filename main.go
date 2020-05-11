package main

import (
	"log"

	"github.com/Yesterday17/pug-backend/auth"
	"github.com/Yesterday17/pug-backend/config"
	"github.com/Yesterday17/pug-backend/controllers"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/Yesterday17/pug/modules"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

func main() {
	// Global config
	cfg := config.LoadConfig()

	// Global PUG Module Manager
	mgr := modules.NewManager()

	// Global DB
	r := gin.Default()
	db := models.InitModels(&cfg.ModelSettings)

	// gin Debug mode
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// global variables
	r.Use(func(ctx *gin.Context) {
		ctx.Set("config", cfg)
		ctx.Set("manager", mgr)
		ctx.Set("db", db)
		ctx.Next()
	})

	// frontend
	r.StaticFile("/", "./public")

	if cfg.CrossOrigin {
		r.Use(cors.Middleware(cors.Config{
			Origins:     "*",
			Methods:     "GET, PUT, POST, PATCH, DELETE",
			Credentials: true,
		}))
	}

	// Session
	r.POST("/session", controllers.SessionCreate, controllers.SessionUpdate)
	r.PUT("/session", controllers.UserRegister)
	r.GET("/session/key", controllers.SessionGetKey)
	r.DELETE("/session", auth.Authorize, controllers.SessionRevoke)

	// User
	r.GET("/user", auth.Authorize, controllers.SessionUpdate, controllers.UserInfoGet)
	r.GET("/user/setting", auth.Authorize, controllers.SessionUpdate, controllers.UserSettingGet)
	r.PATCH("/user/setting", auth.Authorize, controllers.SessionUpdate, controllers.UserSettingPatch)

	// Module
	r.GET("/module", auth.Authorize, controllers.SessionUpdate, controllers.InitModulePipeRestriction, controllers.GetAllModuleInfo)
	r.GET("/module/:module", auth.Authorize, controllers.SessionUpdate, controllers.InitModulePipeRestriction, controllers.GetModuleInfo)
	r.PATCH("/module/:module", auth.Authorize, controllers.SessionUpdate, controllers.InitModulePipeRestriction, controllers.EditModulePipeRestriction)

	// Pipe
	r.GET("/pipe", auth.Authorize, controllers.SessionUpdate, controllers.InitModulePipeRestriction, controllers.GetAllPipeInfo)
	r.GET("/pipe/:module", auth.Authorize, controllers.SessionUpdate, controllers.InitModulePipeRestriction, controllers.GetPipeInfo)
	r.PATCH("/pipe/:module/:pipe", auth.Authorize, controllers.SessionUpdate, controllers.InitModulePipeRestriction, controllers.EditModulePipeRestriction)

	if err := r.Run(cfg.Listen); err != nil {
		log.Fatal(err)
	}
}
