package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"time"

	"github.com/Yesterday17/pug-backend/auth"
	"github.com/Yesterday17/pug-backend/config"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UserLogin(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	cfg := c.MustGet("config").(*config.Config)

	form, err := c.MultipartForm()
	if err != nil {
		c.Abort()
		c.Status(400)
		return
	}

	var name, pass string
	if form.Value["username"] != nil && len(form.Value["username"]) == 1 {
		name = form.Value["username"][0]
	}
	if form.Value["password"] != nil && len(form.Value["password"]) == 1 {
		pass = form.Value["password"][0]
	}

	if name == "" || pass == "" {
		c.Abort()
		c.Status(400)
		return
	}

	password, err := base64.StdEncoding.DecodeString(pass)
	if err != nil {
		c.Abort()
		c.Status(400)
		return
	}

	realPassword, err := rsa.DecryptPKCS1v15(rand.Reader, cfg.KeyPrivate, []byte(password))
	if err != nil {
		c.Abort()
		c.Status(400)
		return
	}

	var user models.User
	db.First(&user, "username = ?", name)
	if user.UUID == "" || user.Password != string(realPassword) {
		c.Abort()
		c.Status(400)
		return
	}

	token, err := auth.GenerateToken(cfg.KeyPrivate, user.UUID, time.Hour*12)
	if err != nil {
		c.Abort()
		c.Status(500)
	}

	c.SetCookie("pug_session", token, int(time.Hour*24), "", "", false, true)
	c.Status(200)
}

func UserLogout(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)

	token, err := auth.CancelToken(cfg.KeyPrivate)
	if err != nil {
		c.Abort()
		c.Status(500)
	}

	c.SetCookie("pug_session", token, int(time.Minute*1), "", "", false, true)
	c.Status(200)
}

func UserRegister(c *gin.Context) {

}
