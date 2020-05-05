package main

import (
	"log"

	"github.com/Yesterday17/pug-backend/auth"
	"github.com/Yesterday17/pug-backend/config"
	"github.com/Yesterday17/pug-backend/controllers"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	r := gin.Default()
	db := models.InitModels(&cfg.ModelSettings)

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(func(ctx *gin.Context) {
		ctx.Set("db", db)
		ctx.Set("config", cfg)
		ctx.Next()
	})

	// put frontend here
	r.StaticFile("/", "./public")

	// methods does not need authorization
	r.POST("/user/login", controllers.UserLogin)
	r.PUT("/user/register", controllers.UserRegister)

	// methods need authorization
	r.POST("/user/logout", auth.Authorize, controllers.UserLogout)
	r.GET("/user", auth.Authorize, nil)

	if err := r.Run(cfg.Listen); err != nil {
		log.Fatal(err)
	}
}
